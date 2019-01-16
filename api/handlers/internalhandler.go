package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/simplexity-ckcclc/gochannel/api/entity"
	"github.com/simplexity-ckcclc/gochannel/api/errorcode"
	"github.com/simplexity-ckcclc/gochannel/common"
	"github.com/sirupsen/logrus"
	"net/http"
)

// The request responds to a url matching:  /internal/appkey/:appkey/evict?nonce=xx&sig=
func EvictAppKeyHandler(c *gin.Context) {
	appkey := c.Param("appkey")
	entity.EvictAppKeySig(appkey)
	common.ApiLog.WithFields(logrus.Fields{
		"appkey": appkey,
	}).Info("Evict app key")
	c.JSON(http.StatusOK, errorcode.SUCCESS)
}

// The request responds to a url matching:  /internal/appkey/:appkey/evict?nonce=xx&sig=
func RegisterAppKeyHandler(c *gin.Context) {
	appkey := c.Param("appkey")
	if err := entity.RegisterAppKeySig(common.DB, appkey); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, errorcode.APP_KEY_NOT_FOUND)
			return
		}
		c.JSON(http.StatusInternalServerError, errorcode.INTERNAL_SERVER_ERROR)
		return
	}
	common.ApiLog.WithFields(logrus.Fields{
		"appkey": appkey,
	}).Info("Register app key")
	c.JSON(http.StatusOK, errorcode.SUCCESS)
}
