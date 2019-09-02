package device

import (
	"database/sql"
	"encoding/json"
	"github.com/simplexity-ckcclc/gochannel/common"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ResponseCode int32

const (
	Success        ResponseCode = 0
	ParameterError ResponseCode = 10000
)

type Callbacker struct {
	db       *sql.DB
	stopChan chan bool
}

type callbackInfo struct {
	Id           int64  `json:"-"`
	AppKey       string `json:"app_key"`
	Channel      string `json:"channel"`
	DeviceId     string `json:"device_id"`
	ClickTime    int64  `json:"click_time"`
	ActivateTime int64  `json:"activate_time"`
}

type callbackResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

func NewCallbacker(db *sql.DB) (cb *Callbacker) {
	cb = &Callbacker{
		db: db,
	}
	go cb.processCallbackInfo()
	return
}

func (cb Callbacker) stop() {
	cb.stopChan <- true
}

func (cb Callbacker) preHandle(devices []*Device) (err error) {
	sqlStr := "INSERT INTO callback_info(app_key, channel_id, device_id, os_type, click_time, activate_time) VALUES "
	vals := []interface{}{}

	for _, device := range devices {
		var deviceId string
		switch device.OsType {
		case common.IOS:
			deviceId = device.Idfa
		case common.Android:
			deviceId = device.Imei
		}

		sqlStr += "(?, ?, ?, ?, ?, ?),"
		vals = append(vals, device.AppKey, device.MatchInfo.Channel, deviceId, device.OsType.String(),
			device.MatchInfo.ClickTime, device.ActivateTime)
	}
	sqlStr = strings.TrimSuffix(sqlStr, ",")

	var stmt *sql.Stmt
	stmt, err = cb.db.Prepare(sqlStr)
	if err != nil {
		return
	}

	_, err = stmt.Exec(vals...)
	return
}

func (cb Callbacker) processCallbackInfo() {
runningLoop:
	for {
		select {
		case <-cb.stopChan:
			common.MatchLogger.Info("Callbacker stop")
			break runningLoop
		default:
			callbackInfos, err := cb.getCallbackInfos(5)
			if err != nil {
				common.MatchLogger.Error("Get callback infos error.", err)
			}

			if len(callbackInfos) > 0 {
				for _, callbackInfo := range callbackInfos {
					//go cb.callback(callbackInfo)
					cb.callback(callbackInfo)
				}
			} else {
				time.Sleep(10 * time.Second)
			}
		}
	}
}

func (cb Callbacker) callback(info callbackInfo) {
	callbackUrl := cb.getCallbackUrl(info.AppKey, info.Channel)
	reqBody, err := json.Marshal(info)
	if err != nil {
		common.MatchLogger.Error("Marshal callback info error.", err)
		return
	}

	resp, err := http.Post(callbackUrl, "application/json", strings.NewReader(string(reqBody[:])))
	if err != nil {
		common.MatchLogger.WithFields(logrus.Fields{
			"callbackUrl": callbackUrl,
			"reqBody":     string(reqBody[:]),
		}).Error("Post callback info error.", err)
		return
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	var cbResp callbackResponse
	if err = json.Unmarshal(respBody, &cbResp); err != nil {
		common.MatchLogger.WithFields(logrus.Fields{
			"respBody": respBody,
		}).Error("UnMarshal callback response error.", err)
		return
	}

	if cbResp.Code == int32(Success) {
		common.MatchLogger.WithFields(logrus.Fields{
			"callbackInfo": info,
			"callbackUrl":  callbackUrl,
		}).Info("Callback success.")
		cb.deleteCallbackInfo(info)
	} else {
		common.MatchLogger.WithFields(logrus.Fields{
			"callbackInfo": info,
			"callbackUrl":  callbackUrl,
			"response":     cbResp,
		}).Error("Callback error.")
		return
	}

}

func (cb *Callbacker) getCallbackInfos(limit int) ([]callbackInfo, error) {
	rows, err := cb.db.Query(`SELECT id, app_key, channel_id, device_id, click_time, activate_time from callback_info limit ?`, strconv.Itoa(limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var callbackInfos []callbackInfo
	for rows.Next() {
		callback := new(callbackInfo)
		if err = rows.Scan(&callback.Id, &callback.AppKey, &callback.Channel, &callback.DeviceId, &callback.ClickTime, &callback.ActivateTime); err != nil {
			return nil, err
		}
		callbackInfos = append(callbackInfos, *callback)
	}
	err = rows.Err()
	return callbackInfos, err
}

func (cb *Callbacker) deleteCallbackInfo(info callbackInfo) error {
	_, err := cb.db.Exec(`DELETE FROM callback_info WHERE id = ?`, info.Id)
	return err
}

func (cb Callbacker) getCallbackUrl(appKey string, channel string) (callbackUrl string) {
	if err := cb.db.QueryRow("SELECT callback_url FROM app_channel WHERE app_key = ? AND channel_id = ?", appKey, channel).
		Scan(&callbackUrl); err != nil {
		common.MatchLogger.WithFields(logrus.Fields{
			"appKey":  appKey,
			"channel": channel,
		}).Error("Get callback url error.", err)
	}

	return
}
