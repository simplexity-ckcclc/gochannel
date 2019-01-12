package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/simplexity-ckcclc/gochannel/api/entity"
	"github.com/simplexity-ckcclc/gochannel/api/errors"
	"github.com/simplexity-ckcclc/gochannel/api/util"
	"github.com/simplexity-ckcclc/gochannel/common"
	"net/http"
	"strings"
)

// The request responds to a url matching:  /ad/click?appKey=foo&deviceId=bar&sig=
func ClickHandler(c *gin.Context) {
	var click entity.ClickInfo
	if err := c.ShouldBind(&click); err != nil {
		c.JSON(http.StatusBadRequest, errors.REQUIRED_PARAMETER_MISSING)
		return
	}

	appkeySig, found := entity.SearchAppKeySig(click.AppKey)
	if !found {
		c.JSON(http.StatusOK, errors.APP_KEY_NOT_FOUND)
		return
	}

	// verify app key sig
	sig := c.Query("sig")
	if sig == "" {
		c.JSON(http.StatusBadRequest, errors.REQUIRED_PARAMETER_MISSING)
		return
	}
	sourceURL := buildSourceURL(click)
	valid, err := util.VerifySig(sourceURL, appkeySig.PublicKey, sig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.INTERNAL_SERVER_ERROR)
		return
	} else if !valid {
		c.JSON(http.StatusOK, errors.SIG_INVALID)
		return
	}

	if err := entity.InsertClickInfo(common.DB, click); err != nil {
		c.JSON(http.StatusInternalServerError, errors.INTERNAL_SERVER_ERROR)
		return
	}

	c.JSON(http.StatusOK, errors.SUCCESS)
}


func buildSourceURL(click entity.ClickInfo) string {
	var sb strings.Builder
	sb.WriteString("appkey=")
	sb.WriteString(click.AppKey)
	sb.WriteString("deviceId=")
	sb.WriteString(click.DeviceId)
	return sb.String()
}