package entity

import (
	"database/sql"
	"github.com/simplexity-ckcclc/gochannel/api/common"
	"sync"
)

type appKeySig struct {
	AppKey     string
	PublicKey  string
	PrivateKey string
}

var appKeySigHolder = struct {
	sync.RWMutex
	appKeys map[string]*appKeySig
}{appKeys: make(map[string]*appKeySig)}

func LoadAppKeySigs(db *sql.DB) error {
	var appkey, pubKey, priKey string
	rows, err := db.Query("select app_key, public_key, private_key from appkey_sig;")
	if err != nil {
		common.ApiLog.Error("Load appkey signature error : ", err)
		return err
	}
	defer rows.Close()

	appKeySigHolder.Lock()
	defer appKeySigHolder.Unlock()
	for rows.Next() {
		err = rows.Scan(&appkey, &pubKey, &priKey)
		if err != nil {
			common.ApiLog.Error("Load appkey signature error : ", err)
			return err
		}
		sig := &appKeySig{
			AppKey:     appkey,
			PublicKey:  pubKey,
			PrivateKey: priKey,
		}
		appKeySigHolder.appKeys[appkey] = sig
	}

	err = rows.Err()
	return err
}

func SearchAppKeySig(appkey string) (*appKeySig, bool) {
	appKeySigHolder.RLock()
	defer appKeySigHolder.RUnlock()
	appkeySig, ok := appKeySigHolder.appKeys[appkey]
	return appkeySig, ok
}

func EvictAppKeySig(db *sql.DB, appkey string) error {
	appKeySigHolder.Lock()
	defer appKeySigHolder.Unlock()
	delete(appKeySigHolder.appKeys, appkey)

	stmt, err := db.Prepare("DELETE FROM appkey_sig WHERE app_key=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(appkey); err != nil {
		return err
	}
	return nil
}

func RegisterAppKeySig(db *sql.DB, appkey string) (*appKeySig, error) {
	pubKey, priKey, err := common.GenerateRSAKeyPair()
	if err != nil {
		return nil, err
	}

	stmt, err := db.Prepare("INSERT INTO appkey_sig (app_key, public_key, private_key) VALUES (?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(appkey, pubKey, priKey); err != nil {
		return nil, err
	}

	sig := &appKeySig{
		AppKey:     appkey,
		PublicKey:  pubKey,
		PrivateKey: priKey,
	}
	appKeySigHolder.Lock()
	defer appKeySigHolder.Unlock()
	appKeySigHolder.appKeys[appkey] = sig
	return sig, nil
}
