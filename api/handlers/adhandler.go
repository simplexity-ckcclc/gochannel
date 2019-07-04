package handlers

import (
	"github.com/gin-gonic/gin"
	api "github.com/simplexity-ckcclc/gochannel/api/common"
	"github.com/simplexity-ckcclc/gochannel/api/entity"
	"github.com/simplexity-ckcclc/gochannel/api/errorcode"
	"github.com/simplexity-ckcclc/gochannel/common"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

// The request responds to a url matching:  /ad/click?appKey=foo&deviceId=bar&sig=
func ClickHandler(c *gin.Context) {
	var click entity.ClickInfo
	if err := c.ShouldBind(&click); err != nil {
		c.JSON(http.StatusBadRequest, errorcode.REQUIRED_PARAMETER_MISSING)
		return
	}

	appkeySig, found := entity.SearchAppKeySig(click.AppKey)
	if !found {
		c.JSON(http.StatusOK, errorcode.APP_KEY_NOT_FOUND)
		return
	}

	//verify app key sig
	sig := c.Query("sig")
	if sig == "" {
		c.JSON(http.StatusBadRequest, errorcode.REQUIRED_PARAMETER_MISSING)
		return
	}
	sourceURL := buildSourceURL(click)
	valid, err := api.VerifyBase64WithRSAPubKey(sourceURL, appkeySig.PublicKey, sig)
	if err != nil {
		api.ApiLog.WithFields(logrus.Fields{
			"pubKey": appkeySig.PublicKey,
		}).Error("Verify click signature - VerifyBase64WithRSAPubKey error : ", err)
		c.JSON(http.StatusInternalServerError, errorcode.INTERNAL_SERVER_ERROR)
		return
	} else if !valid {
		api.ApiLog.WithFields(logrus.Fields{
			"clickInfo": click,
			"pubKey":    appkeySig.PublicKey,
			"sig":       sig,
		}).Info("Invalid signature")
		c.JSON(http.StatusOK, errorcode.SIG_INVALID)
		return
	}

	if err := click.InsertDB(common.DB); err != nil {
		c.JSON(http.StatusInternalServerError, errorcode.INTERNAL_SERVER_ERROR)
		return
	}

	api.ApiLog.WithFields(logrus.Fields{
		"clickInfo": click,
	}).Info("Insert click info")
	c.JSON(http.StatusOK, errorcode.SUCCESS)
}

func buildSourceURL(click entity.ClickInfo) string {
	var sb strings.Builder
	sb.WriteString("appKey=")
	sb.WriteString(click.AppKey)
	sb.WriteString("&deviceId=")
	sb.WriteString(click.DeviceId)
	return sb.String()
}
