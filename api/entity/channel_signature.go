package entity

import (
	"database/sql"
	api "github.com/simplexity-ckcclc/gochannel/api/common"
	"github.com/simplexity-ckcclc/gochannel/common"
	"sync"
)

type channelSig struct {
	AppKey     string
	ChannelId  string
	PublicKey  string
	PrivateKey string
}

var channelSigHolder = struct {
	sync.RWMutex
	channels map[string]*channelSig
}{channels: make(map[string]*channelSig)}

func LoadChannelSigs(db *sql.DB) error {
	var appkey, channelId, pubKey, priKey string
	rows, err := db.Query("select app_key, channel_id, public_key, private_key from channel_sig;")
	defer rows.Close()
	if err != nil {
		common.ApiLogger.Error("Load channel signature error : ", err)
		return err
	}

	channelSigHolder.Lock()
	defer channelSigHolder.Unlock()
	for rows.Next() {
		err = rows.Scan(&appkey, &channelId, &pubKey, &priKey)
		if err != nil {
			common.ApiLogger.Error("Load channel signature error : ", err)
			return err
		}
		sig := &channelSig{
			AppKey:     appkey,
			ChannelId:  channelId,
			PublicKey:  pubKey,
			PrivateKey: priKey,
		}
		channelSigHolder.channels[channelId] = sig
	}

	err = rows.Err()
	return err
}

func SearchChannelSig(channelId string) (*channelSig, bool) {
	channelSigHolder.RLock()
	defer channelSigHolder.RUnlock()
	channelSig, ok := channelSigHolder.channels[channelId]
	return channelSig, ok
}

func EvictChannelSig(db *sql.DB, channelId string) error {
	channelSigHolder.Lock()
	defer channelSigHolder.Unlock()
	delete(channelSigHolder.channels, channelId)

	stmt, err := db.Prepare("DELETE FROM channel_sig WHERE channel_id=?")
	defer stmt.Close()
	if err != nil {
		return err
	}

	if _, err = stmt.Exec(channelId); err != nil {
		return err
	}
	return nil
}

func RegisterChannelSig(db *sql.DB, appkey string, channelId string) (*channelSig, error) {
	pubKey, priKey, err := api.GenerateRSAKeyPair()
	if err != nil {
		return nil, err
	}

	stmt, err := db.Prepare("INSERT INTO channel_sig (app_key, channel_id, public_key, private_key) VALUES (?, ?, ?, ?)")
	defer stmt.Close()
	if err != nil {
		return nil, err
	}

	if _, err = stmt.Exec(appkey, channelId, pubKey, priKey); err != nil {
		return nil, err
	}

	sig := &channelSig{
		AppKey:     appkey,
		ChannelId:  channelId,
		PublicKey:  pubKey,
		PrivateKey: priKey,
	}
	channelSigHolder.Lock()
	defer channelSigHolder.Unlock()
	channelSigHolder.channels[channelId] = sig
	return sig, nil
}
