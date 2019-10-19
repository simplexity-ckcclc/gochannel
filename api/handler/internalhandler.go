package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/simplexity-ckcclc/gochannel/api/appchannel"
	api "github.com/simplexity-ckcclc/gochannel/api/common"
	"github.com/simplexity-ckcclc/gochannel/common"
	"github.com/simplexity-ckcclc/gochannel/common/logger"
	"net/http"
	"strings"
)

// The request responds to a url matching:  /internal/channel/:channel/evict?nonce=xx&sig=
func EvictChannelHandler(c *gin.Context) {
	appKey := c.Query("appKey")
	channelId := c.Query("channel")
	if err := appchannel.EvictAppChannel(common.DB, appKey, channelId); err != nil {
		api.ResponseJSONWithExtraMsg(c, http.StatusInternalServerError, api.InternalServerError, err.Error())
		return
	}

	logger.ApiLogger.With(logger.Fields{
		"channelId": channelId,
	}).Info("Evict channel")
	api.ResponseJSON(c, http.StatusOK, api.Success)
}

// The request responds to a url matching:  /internal/channel/:channel/register?nonce=xx&sig=
func RegisterChannelHandler(c *gin.Context) {
	var ac common.AppChannel
	if err := c.BindJSON(&ac); err != nil {
		api.ResponseJSON(c, http.StatusBadRequest, api.RequiredParameterError)
		return
	}

	if len(ac.AppKey) == 0 || len(ac.ChannelId) == 0 {
		api.ResponseJSON(c, http.StatusOK, api.RequiredParameterError)
		return
	}

	if _, found := appchannel.SearchAppChannel(ac.AppKey, ac.ChannelId); found {
		api.ResponseJSON(c, http.StatusOK, api.DuplicateChannel)
		return
	}

	var ct common.ChannelType
	if ct = common.ParseChannelType(strings.ToLower(ac.ChannelType.String())); ct == common.UnknownChannelType {
		api.ResponseJSON(c, http.StatusOK, api.ChannelTypeInvalid)
		return
	} else {
		ac.ChannelType = ct
	}

	if err := appchannel.RegisterAppChannel(common.DB, &ac); err != nil {
		logger.ApiLogger.Warning("RegisterAppChannel error : ", err)
		api.ResponseJSONWithExtraMsg(c, http.StatusInternalServerError, api.InternalServerError, err.Error())
		return
	}

	logger.ApiLogger.With(logger.Fields{
		"appChannel": ac,
	}).Info("Register channel")
	api.ResponseJSONWithData(c, http.StatusOK, api.Success, ac)
}
