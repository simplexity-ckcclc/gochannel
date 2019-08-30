package appchannel

import (
	"database/sql"
	api "github.com/simplexity-ckcclc/gochannel/api/common"
	"github.com/simplexity-ckcclc/gochannel/common"
	"sync"
)

var appChannelHolder = struct {
	sync.RWMutex
	channels map[string]*common.AppChannel
}{channels: make(map[string]*common.AppChannel)}

func LoadAppChannels(db *sql.DB) error {
	rows, err := db.Query("select app_key, channel_id, channel_type, public_key, private_key from app_channel;")
	defer rows.Close()
	if err != nil {
		common.ApiLogger.Error("Load channel signature error : ", err)
		return err
	}

	var ac *common.AppChannel
	var appKey, channelId, channelType, pubKey, priKey string
	appChannelHolder.Lock()
	defer appChannelHolder.Unlock()
	for rows.Next() {
		err = rows.Scan(&appKey, &channelId, &channelType, &pubKey, &priKey)
		if err != nil {
			common.ApiLogger.Error("Load channel signature error : ", err)
			return err
		}

		ac = &common.AppChannel{
			AppKey:      appKey,
			ChannelId:   channelId,
			ChannelType: common.ParseChannelType(channelType),
			PublicKey:   pubKey,
			PrivateKey:  priKey,
		}
		appChannelHolder.channels[ac.ChannelId] = ac
	}

	err = rows.Err()
	return err
}

func SearchAppChannel(channelId string) (*common.AppChannel, bool) {
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

func RegisterAppChannel(db *sql.DB, appkey string, channelId string, channelType common.ChannelType) (*common.AppChannel, error) {
	pubKey, priKey, err := api.GenerateRSAKeyPair()
	if err != nil {
		return nil, err
	}

	stmt, err := db.Prepare("INSERT INTO app_channel (app_key, channel_id, channel_type, public_key, private_key) VALUES (?, ?, ?, ?, ?)")
	defer stmt.Close()
	if err != nil {
		return nil, err
	}

	if _, err = stmt.Exec(appkey, channelId, channelType, pubKey, priKey); err != nil {
		return nil, err
	}

	ac := &common.AppChannel{
		AppKey:      appkey,
		ChannelId:   channelId,
		ChannelType: channelType,
		PublicKey:   pubKey,
		PrivateKey:  priKey,
	}
	appChannelHolder.Lock()
	defer appChannelHolder.Unlock()
	appChannelHolder.channels[channelId] = ac
	return ac, nil
}
