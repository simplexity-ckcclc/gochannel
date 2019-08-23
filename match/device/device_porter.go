package device

import (
	"database/sql"
	"github.com/simplexity-ckcclc/gochannel/common"
)

type DevicePorter struct {
	conf common.ConfYaml
	db   sql.DB
}

func (porter DevicePorter) transferDevices() {
	devices, err := porter.getSdkDevices(50)
	if err != nil {
		// log
	} else {
		_ := porter.putIntoEs(devices)
	}
}

func (porter DevicePorter) getSdkDevices(limit int32) ([]Device, error) {
	rows, err := porter.db.Query(`SELECT id, idfa, imei, app_key, os_type, os_version, source_ip, 
        language, resolution, receive_time FROM sdk_report`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	devices := make([]Device, limit)
	for rows.Next() {
		device := Device{}
		err = rows.Scan(&device.Id, &device.Idfa, &device.Imei, &device.AppKey, &device.OsType, &device.OsVersion,
			&device.SourceIp, &device.Language, &device.Resolution, &device.ReceiveTime)
		if err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}
	err = rows.Err()
	return devices, err
}

func (porter DevicePorter) putIntoEs(devices []Device) (err error) {
	err = nil
	return
}
