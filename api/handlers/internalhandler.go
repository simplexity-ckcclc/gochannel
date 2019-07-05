package handlers

import (
	"github.com/gin-gonic/gin"
	api "github.com/simplexity-ckcclc/gochannel/api/common"
	"github.com/simplexity-ckcclc/gochannel/api/entity"
	"github.com/sirupsen/logrus"
	"net/http"
)

// The request responds to a url matching:  /internal/channel/:channel/evict?nonce=xx&sig=
func EvictChannelHandler(c *gin.Context) {
	channelId := c.Param("channel")
	if err := entity.EvictChannelSig(api.DB, channelId); err != nil {
		api.ResponseJSONWithExtraMsg(c, http.StatusInternalServerError, api.INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	api.ApiLog.WithFields(logrus.Fields{
		"channelId": channelId,
	}).Info("Evict channel")
	api.ResponseJSON(c, http.StatusOK, api.SUCCESS)
}

// The request responds to a url matching:  /internal/channel/:channel/register?nonce=xx&sig=
func RegisterChannelHandler(c *gin.Context) {
	appkey := c.Query("appKey")
	channelId := c.Param("channel")
	if len(appkey) == 0 || len(channelId) == 0 {
		api.ResponseJSON(c, http.StatusOK, api.REQUIRED_PARAMETER_MISSING)
		return
	}

	if _, found := entity.SearchChannelSig(channelId); found {
		api.ResponseJSON(c, http.StatusOK, api.DUPLICATE_CHANNEL)
		return
	}

	sig, err := entity.RegisterChannelSig(api.DB, appkey, channelId)
	if err != nil {
		api.ApiLog.Warning("RegisterChannelSig error : ", err)
		api.ResponseJSONWithExtraMsg(c, http.StatusInternalServerError, api.INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	api.ApiLog.WithFields(logrus.Fields{
		"channelSig": sig,
	}).Info("Register channel")
	api.ResponseJSONWithData(c, http.StatusOK, api.SUCCESS, sig)
}
