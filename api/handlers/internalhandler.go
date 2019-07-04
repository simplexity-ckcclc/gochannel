package handlers

import (
	"github.com/gin-gonic/gin"
	api "github.com/simplexity-ckcclc/gochannel/api/common"
	"github.com/simplexity-ckcclc/gochannel/api/entity"
	"github.com/simplexity-ckcclc/gochannel/common"
	"github.com/sirupsen/logrus"
	"net/http"
)

// The request responds to a url matching:  /internal/appkey/:appkey/evict?nonce=xx&sig=
func EvictAppKeyHandler(c *gin.Context) {
	appkey := c.Param("appkey")
	if err := entity.EvictAppKeySig(common.DB, appkey); err != nil {
		api.ResponseJSONWithExtraMsg(c, http.StatusInternalServerError, api.INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	api.ApiLog.WithFields(logrus.Fields{
		"appkey": appkey,
	}).Info("Evict app key")
	api.ResponseJSON(c, http.StatusOK, api.SUCCESS)
}

// The request responds to a url matching:  /internal/appkey/:appkey/evict?nonce=xx&sig=
func RegisterAppKeyHandler(c *gin.Context) {
	appkey := c.Param("appkey")
	if len(appkey) == 0 {
		api.ResponseJSON(c, http.StatusOK, api.REQUIRED_PARAMETER_MISSING)
		return
	}

	if _, found := entity.SearchAppKeySig(appkey); found {
		api.ResponseJSON(c, http.StatusOK, api.DUPLICATE_APP_KEY)
		return
	}

	sig, err := entity.RegisterAppKeySig(common.DB, appkey)
	if err != nil {
		api.ApiLog.Warning("RegisterAppKeySig error : ", err)
		api.ResponseJSONWithExtraMsg(c, http.StatusInternalServerError, api.INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	api.ApiLog.WithFields(logrus.Fields{
		"appkeySig": sig,
	}).Info("Register app key")
	api.ResponseJSONWithData(c, http.StatusOK, api.SUCCESS, sig)
}
