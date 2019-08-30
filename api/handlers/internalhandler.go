package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/simplexity-ckcclc/gochannel/api/appchannel"
	api "github.com/simplexity-ckcclc/gochannel/api/common"
	"github.com/simplexity-ckcclc/gochannel/common"
	"github.com/sirupsen/logrus"
	"net/http"
)

// The request responds to a url matching:  /internal/channel/:channel/evict?nonce=xx&sig=
func EvictChannelHandler(c *gin.Context) {
	channelId := c.Param("channel")
	if err := appchannel.EvictAppChannel(common.DB, channelId); err != nil {
		api.ResponseJSONWithExtraMsg(c, http.StatusInternalServerError, api.InternalServerError, err.Error())
		return
	}

	common.ApiLogger.WithFields(logrus.Fields{
		"channelId": channelId,
	}).Info("Evict channel")
	api.ResponseJSON(c, http.StatusOK, api.Success)
}

// The request responds to a url matching:  /internal/channel/:channel/register?nonce=xx&sig=
func RegisterChannelHandler(c *gin.Context) {
	appkey := c.Query("appKey")
	channelId := c.Param("channel")
	channelType := c.Param("channelType")
	if len(appkey) == 0 || len(channelId) == 0 {
		api.ResponseJSON(c, http.StatusOK, api.RequiredParameterMissing)
		return
	}

	if _, found := appchannel.SearchAppChannel(channelId); found {
		api.ResponseJSON(c, http.StatusOK, api.DuplicateChannel)
		return
	}

	var ct common.ChannelType
	if ct = common.ParseChannelType(channelType); ct == common.UnknownChannelType {
		api.ResponseJSON(c, http.StatusOK, api.ChannelTypeInvalid)
		return
	}

	sig, err := appchannel.RegisterAppChannel(common.DB, appkey, channelId, ct)
	if err != nil {
		common.ApiLogger.Warning("RegisterAppChannel error : ", err)
		api.ResponseJSONWithExtraMsg(c, http.StatusInternalServerError, api.InternalServerError, err.Error())
		return
	}

	common.ApiLogger.WithFields(logrus.Fields{
		"channelSig": sig,
	}).Info("Register channel")
	api.ResponseJSONWithData(c, http.StatusOK, api.Success, sig)
}
