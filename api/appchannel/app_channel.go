package appchannel

import (
	"database/sql"
	api "github.com/simplexity-ckcclc/gochannel/api/common"
	"github.com/simplexity-ckcclc/gochannel/common"
	"github.com/simplexity-ckcclc/gochannel/common/logger"
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
		logger.ApiLogger.Error("Load channel signature error : ", err)
		return err
	}

	var ac *common.AppChannel
	var appKey, channelId, channelType, pubKey, priKey string
	appChannelHolder.Lock()
	defer appChannelHolder.Unlock()
	for rows.Next() {
		err = rows.Scan(&appKey, &channelId, &channelType, &pubKey, &priKey)
		if err != nil {
			logger.ApiLogger.Error("Load channel signature error : ", err)
			return err
		}

		ac = &common.AppChannel{
			AppKey:      appKey,
			ChannelId:   channelId,
			ChannelType: common.ParseChannelType(channelType),
			PublicKey:   pubKey,
			PrivateKey:  priKey,
		}
		appChannelHolder.channels[appKeyChannel(ac.AppKey, ac.ChannelId)] = ac
	}

	err = rows.Err()
	return err
}

func SearchAppChannel(appKey string, channel string) (*common.AppChannel, bool) {
	appChannelHolder.RLock()
	defer appChannelHolder.RUnlock()
	channelSig, ok := appChannelHolder.channels[appKeyChannel(appKey, channel)]
	return channelSig, ok
}

func EvictAppChannel(db *sql.DB, appKey string, channel string) error {
	appChannelHolder.Lock()
	defer appChannelHolder.Unlock()
	delete(appChannelHolder.channels, appKeyChannel(appKey, channel))

	stmt, err := db.Prepare("DELETE FROM app_channel WHERE app_key = ? and channel_id = ?")
	defer stmt.Close()
	if err != nil {
		return err
	}

	if _, err = stmt.Exec(appKey, channel); err != nil {
		return err
	}
	return nil
}

func RegisterAppChannel(db *sql.DB, ac *common.AppChannel) error {
	pubKey, priKey, err := api.GenerateRSAKeyPair()
	if err != nil {
		return err
	}

	ac.PublicKey = pubKey
	ac.PrivateKey = priKey

	stmt, err := db.Prepare("INSERT INTO app_channel (app_key, channel_id, channel_type, public_key, private_key, callback_url) " +
		"VALUES (?, ?, ?, ?, ?, ?)")
	defer stmt.Close()
	if err != nil {
		return err
	}

	if _, err = stmt.Exec(ac.AppKey, ac.ChannelId, ac.ChannelType.String(), pubKey, priKey, ac.CallbackUrl); err != nil {
		return err
	}

	appChannelHolder.Lock()
	defer appChannelHolder.Unlock()
	appChannelHolder.channels[appKeyChannel(ac.AppKey, ac.ChannelId)] = ac
	return nil
}

func appKeyChannel(appKey string, channel string) string {
	return appKey + "-" + channel
}
