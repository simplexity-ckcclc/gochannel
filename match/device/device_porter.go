package device

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/simplexity-ckcclc/gochannel/common"
	"github.com/simplexity-ckcclc/gochannel/common/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/olivere/elastic.v6"
	"strconv"
	"time"
)

type DevicePorter struct {
	db       *sql.DB
	esClient *elastic.Client
}

func NewDevicePorter(database *sql.DB) (*DevicePorter, error) {
	client, err := elastic.NewClient(elastic.SetURL(config.GetString(config.EsServer)))
	if err != nil {
		return nil, err
	}

	return &DevicePorter{
		db:       database,
		esClient: client,
	}, nil
}

func (porter DevicePorter) TransferDevices() {
	esDeviceIndex := config.GetString(config.EsDeviceIndex)
	for {
		devices, err := porter.getSdkDevices(config.GetInt(config.EsDeviceBatchSize))
		if err != nil {
			common.MatchLogger.WithFields(logrus.Fields{}).Error("Get devices from db error : ", err)
			continue
		}

		if len(devices) > 0 {
			porter.putDevicesIntoEs(devices, esDeviceIndex)
		}

		time.Sleep(10 * time.Second)

	}
}

func (porter DevicePorter) getSdkDevices(limit int) ([]Device, error) {
	rows, err := porter.db.Query(`SELECT id, idfa, imei, app_key, os_type, os_version, source_ip, 
        language, resolution, receive_time FROM sdk_device_report limit` + strconv.Itoa(limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []Device
	for rows.Next() {
		device := new(Device)
		if err = rows.Scan(&device.Id, &device.Idfa, &device.Imei, &device.AppKey, &device.OsType, &device.OsVersion,
			&device.SourceIp, &device.Language, &device.Resolution, &device.ReceiveTime); err != nil {
			return nil, err
		}
		devices = append(devices, *device)
	}
	err = rows.Err()
	return devices, err
}

func (porter DevicePorter) putDevicesIntoEs(devices []Device, index string) {
	bulkRequest := porter.esClient.Bulk()
	for _, device := range devices {
		deviceJson, err := json.Marshal(device)
		if err != nil {
			common.MatchLogger.WithFields(logrus.Fields{
				"device": device,
			}).Error("JSON marshal device error : ", err)
			continue
		}

		req := elastic.NewBulkIndexRequest().
			Index(index).
			Type(device.AppKey).
			Doc(string(deviceJson))
		bulkRequest.Add(req)
	}

	bulkResponse, err := bulkRequest.Do(context.Background())
	if err != nil {
		common.MatchLogger.WithFields(logrus.Fields{
			"devices": devices,
		}).Error("Bulk put device doc error : ", err)
		return
	}

	failed := bulkResponse.Failed()
	for _, failedResp := range failed {
		common.MatchLogger.WithFields(logrus.Fields{
			"id":       failedResp.Id,
			"errCause": failedResp.Error,
		}).Error("Bulk put device doc error : ", err)
	}
}
