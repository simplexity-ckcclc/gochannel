package entity

import (
	"database/sql"
	"fmt"
	"sync"
)

var appKeySigHolder = struct {
	sync.RWMutex
	appKeys map[string]appKeySig
}{appKeys: make(map[string]appKeySig)}

type appKeySig struct {
	appKey     string
	PublicKey  string
	privateKey string
}

func LoadAppKeySigs(db *sql.DB) error {
	var appkeySig appKeySig
	rows, err := db.Query("select app_key, public_key, private_key from appkey_sig;")
	if err != nil {
		fmt.Print(err.Error())
		return err
	}
	defer rows.Close()

	appKeySigHolder.Lock()
	defer appKeySigHolder.Unlock()
	for rows.Next() {
		err = rows.Scan(&appkeySig.appKey, &appkeySig.PublicKey, &appkeySig.privateKey)
		if err != nil {
			fmt.Print(err.Error())
		}
		appKeySigHolder.appKeys[appkeySig.appKey] = appkeySig
	}

	err = rows.Err()
	return err
}

func SearchAppKeySig(appkey string) (appKeySig, bool) {
	appKeySigHolder.RLock()
	defer appKeySigHolder.RUnlock()
	appkeySig, ok := appKeySigHolder.appKeys[appkey]
	return appkeySig, ok
}

func EvictAppKeySig(appkey string) {
	appKeySigHolder.Lock()
	defer appKeySigHolder.Unlock()
	delete(appKeySigHolder.appKeys, appkey)
}

func RegisterAppKeySig(db *sql.DB, appkey string) error {
	var appkeySig appKeySig
	if err := db.QueryRow("select app_key, public_key, private_key from appkey_sig where app_key = ?;", appkey).Scan(&appkeySig.appKey, &appkeySig.PublicKey, &appkeySig.privateKey); err != nil {
		fmt.Print(err.Error())
		return err
	}

	appKeySigHolder.Lock()
	defer appKeySigHolder.Unlock()
	appKeySigHolder.appKeys[appkeySig.appKey] = appkeySig
	return nil
}
