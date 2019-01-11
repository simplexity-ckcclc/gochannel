package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/simplexity-ckcclc/gochannel/api/entity"
	"github.com/simplexity-ckcclc/gochannel/api/errors"
	"net/http"
)

// The request responds to a url matching:  /internal/appkey/:appkey/evict?signature=
func EvictAppKeyHandler(c *gin.Context) {
	appkey := c.Param("appkey")
	//signature := c.Query("signature")

	entity.EvictAppKeySig(appkey)
	c.JSON(http.StatusOK, errors.SUCCESS)
}
