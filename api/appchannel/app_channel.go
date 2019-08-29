package appchannel

import (
	"database/sql"
	api "github.com/simplexity-ckcclc/gochannel/api/common"
	"github.com/simplexity-ckcclc/gochannel/common"
	"sync"
)

type appChannel struct {
	AppKey     string
	ChannelId  string
	PublicKey  string
	PrivateKey string
}

var appChannelHolder = struct {
	sync.RWMutex
	channels map[string]*appChannel
}{channels: make(map[string]*appChannel)}

func LoadChannelSigs(db *sql.DB) error {
	var appkey, channelId, pubKey, priKey string
	rows, err := db.Query("select app_key, channel_id, public_key, private_key from app_channel;")
	defer rows.Close()
	if err != nil {
		common.ApiLogger.Error("Load channel signature error : ", err)
		return err
	}

	appChannelHolder.Lock()
	defer appChannelHolder.Unlock()
	for rows.Next() {
		err = rows.Scan(&appkey, &channelId, &pubKey, &priKey)
		if err != nil {
			common.ApiLogger.Error("Load channel signature error : ", err)
			return err
		}
		sig := &appChannel{
			AppKey:     appkey,
			ChannelId:  channelId,
			PublicKey:  pubKey,
			PrivateKey: priKey,
		}
		appChannelHolder.channels[channelId] = sig
	}

	err = rows.Err()
	return err
}

func SearchAppChannel(channelId string) (*appChannel, bool) {
	appChannelHolder.RLock()
	defer appChannelHolder.RUnlock()
	channelSig, ok := appChannelHolder.channels[channelId]
	return channelSig, ok
}

func EvictAppChannel(db *sql.DB, channelId string) error {
	appChannelHolder.Lock()
	defer appChannelHolder.Unlock()
	delete(appChannelHolder.channels, channelId)

	stmt, err := db.Prepare("DELETE FROM app_channel WHERE channel_id=?")
	defer stmt.Close()
	if err != nil {
		return err
	}

	if _, err = stmt.Exec(channelId); err != nil {
		return err
	}
	return nil
}

func RegisterChannelSig(db *sql.DB, appkey string, channelId string) (*appChannel, error) {
	pubKey, priKey, err := api.GenerateRSAKeyPair()
	if err != nil {
		return nil, err
	}

	stmt, err := db.Prepare("INSERT INTO app_channel (app_key, channel_id, public_key, private_key) VALUES (?, ?, ?, ?)")
	defer stmt.Close()
	if err != nil {
		return nil, err
	}

	if _, err = stmt.Exec(appkey, channelId, pubKey, priKey); err != nil {
		return nil, err
	}

	sig := &appChannel{
		AppKey:     appkey,
		ChannelId:  channelId,
		PublicKey:  pubKey,
		PrivateKey: priKey,
	}
	appChannelHolder.Lock()
	defer appChannelHolder.Unlock()
	appChannelHolder.channels[channelId] = sig
	return sig, nil
}
