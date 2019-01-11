package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/simplexity-ckcclc/gochannel/api/entity"
	"github.com/simplexity-ckcclc/gochannel/api/errors"
	"github.com/simplexity-ckcclc/gochannel/common"
	"net/http"
)

// The request responds to a url matching:  /ad/click?appKey=foo&deviceId=bar
func ClickHandler(c *gin.Context) {
	var click entity.ClickInfo
	if err := c.ShouldBind(&click); err != nil {
		c.JSON(http.StatusBadRequest, errors.REQUIRED_PARAMETER_MISSING)
		return
	}

	appkeySig, ok := entity.SearchAppKeySig(click.AppKey)
	if !ok {
		c.JSON(http.StatusOK, errors.APP_KEY_NOT_FOUND)
		return
	}
	fmt.Println(appkeySig)	// todo verify appkey signature

	if err := entity.InsertClickInfo(common.DB, click); err != nil {
		c.JSON(http.StatusInternalServerError, errors.INTERNAL_SERVER_ERROR)
		return
	}

	c.JSON(http.StatusOK, errors.SUCCESS)
}
