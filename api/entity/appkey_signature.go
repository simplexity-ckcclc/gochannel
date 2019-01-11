package entity

import (
	"database/sql"
	"fmt"
    "sync"
)

var appKeySigHoler = struct {
        sync.RWMutex
        appKeys map[string]appKeySig
    }{appKeys:  make(map[string]appKeySig)}


type appKeySig struct {
	appKey     string
	publicKey  string
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

	appKeySigHoler.Lock()
	defer appKeySigHoler.Unlock()
	for rows.Next() {
		err = rows.Scan(&appkeySig.appKey, &appkeySig.publicKey, &appkeySig.privateKey)
		if err != nil {
			fmt.Print(err.Error())
		}
		appKeySigHoler.appKeys[appkeySig.appKey] = appkeySig
	}

	err = rows.Err()
	return err
}

func SearchAppKeySig(appkey string) (appKeySig, bool) {
    appKeySigHoler.RLock()
    defer appKeySigHoler.RUnlock()
	appkeySig, ok := appKeySigHoler.appKeys[appkey]
	return appkeySig, ok
}

func EvictAppKeySig(appkey string) {
    appKeySigHoler.Lock()
    defer appKeySigHoler.Unlock()
    delete(appKeySigHoler.appKeys, appkey)
}

func Verify(encryptedUrl string, privateKey string, signature string) (bool, error) {
	return true, nil
}
