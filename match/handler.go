package match

import (
	"github.com/golang/protobuf/proto"
	"github.com/simplexity-ckcclc/gochannel/common"
	pb "github.com/simplexity-ckcclc/gochannel/match/proto"
	"github.com/sirupsen/logrus"
)

type messageHandler interface {
	handle(message []byte)
}

type MatchHandler struct {
}

func (handler MatchHandler) handle(message []byte) {
	device := &pb.SdkDeviceReport{}
	if err := proto.Unmarshal(message, device); err != nil {
		common.MatchLogger.Error("Parse device error : ", err)
	} else {
		if err := insertIntoDB(device); err != nil {
			common.MatchLogger.WithFields(logrus.Fields{
				"Device ": device,
			}).Error("Insert into DB error", err)
		} else {
			common.MatchLogger.WithFields(logrus.Fields{
				"Device ": device,
			}).Debug("Insert into DB")
		}
	}

}

func insertIntoDB(device *pb.SdkDeviceReport) error {
	stmt, err := common.DB.Prepare("INSERT INTO sdk_device_report (imei, idfa, app_key, channel_id, resolution, " +
		"language, os_type, os_version, receive_time, source_ip) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	defer stmt.Close()
	if err != nil {
		return err
	}

	_, err = stmt.Exec(device.Imei, device.Idfa, device.AppKey, device.Channel, device.Resolution,
		device.Language, device.OsType, device.OsVersion, device.ReceiveTime, device.SourceIp)
	return err
}
