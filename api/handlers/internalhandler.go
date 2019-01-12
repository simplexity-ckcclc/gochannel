package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/simplexity-ckcclc/gochannel/api/entity"
	"github.com/simplexity-ckcclc/gochannel/api/errors"
	"github.com/simplexity-ckcclc/gochannel/common"
	"net/http"
)

// The request responds to a url matching:  /internal/appkey/:appkey/evict?nonce=xx&sig=
func EvictAppKeyHandler(c *gin.Context) {
	appkey := c.Param("appkey")
	//signature := c.Query("signature")

	entity.EvictAppKeySig(appkey)
	c.JSON(http.StatusOK, errors.SUCCESS)
}


// The request responds to a url matching:  /internal/appkey/:appkey/evict?nonce=xx&sig=
func RegisterAppKeyHandler(c *gin.Context) {
	appkey := c.Param("appkey")
	//signature := c.Query("signature")

	if err := entity.RegisterAppKeySig(common.DB, appkey); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, errors.APP_KEY_NOT_FOUND)
			return
		}
		c.JSON(http.StatusInternalServerError, errors.INTERNAL_SERVER_ERROR)
		return
	}
	c.JSON(http.StatusOK, errors.SUCCESS)
}