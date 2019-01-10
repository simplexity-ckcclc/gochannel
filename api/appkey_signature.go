package api

import (
	"database/sql"
	"fmt"
)

var (
	appKeys = make(map[string]appKeySig)
)

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

	for rows.Next() {
		err = rows.Scan(&appkeySig.appKey, &appkeySig.publicKey, &appkeySig.privateKey)
		if err != nil {
			fmt.Print(err.Error())
		}
		appKeys[appkeySig.appKey] = appkeySig
	}

	err = rows.Err()
	return err
}

func searchAppKeySig(appkey string) (appKeySig, bool) {
	appkeySig, ok := appKeys[appkey]
	return appkeySig, ok
}

func verify(encryptedUrl string, privateKey string, signature string) (bool, error) {
	return true, nil
}
