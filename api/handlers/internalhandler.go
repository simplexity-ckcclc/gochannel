package handlers

import (
	"github.com/gin-gonic/gin"
	api "github.com/simplexity-ckcclc/gochannel/api/common"
	"github.com/simplexity-ckcclc/gochannel/api/entity"
	"github.com/simplexity-ckcclc/gochannel/api/errorcode"
	"github.com/simplexity-ckcclc/gochannel/common"
	"github.com/sirupsen/logrus"
	"net/http"
)

// The request responds to a url matching:  /internal/appkey/:appkey/evict?nonce=xx&sig=
func EvictAppKeyHandler(c *gin.Context) {
	appkey := c.Param("appkey")
	if err := entity.EvictAppKeySig(common.DB, appkey); err != nil {
		c.JSON(http.StatusInternalServerError, errorcode.INTERNAL_SERVER_ERROR)
		return
	}

	api.ApiLog.WithFields(logrus.Fields{
		"appkey": appkey,
	}).Info("Evict app key")
	c.JSON(http.StatusOK, errorcode.SUCCESS)
}

// The request responds to a url matching:  /internal/appkey/:appkey/evict?nonce=xx&sig=
func RegisterAppKeyHandler(c *gin.Context) {
	appkey := c.Param("appkey")
	if len(appkey) == 0 {
		c.JSON(http.StatusOK, errorcode.REQUIRED_PARAMETER_MISSING)
		return
	}

	if _, found := entity.SearchAppKeySig(appkey); found {
		c.JSON(http.StatusOK, errorcode.DUPLICATE_APP_KEY)
		return
	}

	sig, err := entity.RegisterAppKeySig(common.DB, appkey)
	if err != nil {
		api.ApiLog.Warning("RegisterAppKeySig error : ", err)
		c.JSON(http.StatusInternalServerError, errorcode.INTERNAL_SERVER_ERROR)
		return
	}

	api.ApiLog.WithFields(logrus.Fields{
		"appkeySig": sig,
	}).Info("Register app key")
	c.JSON(http.StatusOK, errorcode.SUCCESS)
}
