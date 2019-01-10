package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/simplexity-ckcclc/gochannel/common"
	"net/http"
)

// The request responds to a url matching:  /ad/click?appKey=foo&deviceId=bar
func clickHandler(c *gin.Context) {
	var click clickInfo
	if err := c.ShouldBind(&click); err != nil {
		c.JSON(http.StatusBadRequest, REQUIRED_PARAMETER_MISSING)
		return
	}

	appkeySig, ok := searchAppKeySig(click.AppKey)
	if !ok {
		c.JSON(http.StatusOK, APP_KEY_NOT_FOUND)
		return
	}
	fmt.Println(appkeySig)

	if err := insertClickInfo(common.DB, click); err != nil {
		c.JSON(http.StatusInternalServerError, INTERNAL_SERVER_ERROR)
		return
	}

	c.JSON(http.StatusOK, SUCCESS)
}
