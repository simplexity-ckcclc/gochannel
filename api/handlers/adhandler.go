package handlers

import (
	"github.com/gin-gonic/gin"
	api "github.com/simplexity-ckcclc/gochannel/api/common"
	"github.com/simplexity-ckcclc/gochannel/api/entity"
	"github.com/simplexity-ckcclc/gochannel/common"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

// The request responds to a url matching:  /ad/click?appKey=foo&deviceId=bar&sig=
func ClickHandler(c *gin.Context) {
	var click entity.ClickInfo
	if err := c.ShouldBind(&click); err != nil {
		api.ResponseJSON(c, http.StatusBadRequest, api.REQUIRED_PARAMETER_MISSING)
		return
	}

	appkeySig, found := entity.SearchChannelSig(click.ChannelId)
	if !found {
		api.ResponseJSON(c, http.StatusOK, api.CHANNEL_NOT_FOUND)
		return
	}

	if appkeySig.AppKey != click.AppKey {
		api.ResponseJSON(c, http.StatusOK, api.APP_KEY_UNMATCHED)
		return
	}

	//verify app key sig
	sig := c.Query("sig")
	if sig == "" {
		api.ResponseJSON(c, http.StatusBadRequest, api.REQUIRED_PARAMETER_MISSING)
		return
	}
	sourceURL := buildSourceURL(click)
	valid, err := api.VerifyBase64WithRSAPubKey(sourceURL, appkeySig.PublicKey, sig)
	if err != nil {
		common.ApiLogger.WithFields(logrus.Fields{
			"pubKey": appkeySig.PublicKey,
		}).Error("Verify click signature - VerifyBase64WithRSAPubKey error : ", err)
		api.ResponseJSONWithExtraMsg(c, http.StatusInternalServerError, api.INTERNAL_SERVER_ERROR, err.Error())
		return
	} else if !valid {
		common.ApiLogger.WithFields(logrus.Fields{
			"clickInfo": click,
			"pubKey":    appkeySig.PublicKey,
			"sig":       sig,
		}).Info("Invalid signature")
		api.ResponseJSON(c, http.StatusOK, api.SIG_INVALID)
		return
	}

	if err := click.InsertDB(common.DB); err != nil {
		common.ApiLogger.WithFields(logrus.Fields{
			"clickInfo": click,
			"pubKey":    appkeySig.PublicKey,
			"sig":       sig,
		}).Error("Insert DB error : ", err)
		api.ResponseJSONWithExtraMsg(c, http.StatusInternalServerError, api.INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	common.ApiLogger.WithFields(logrus.Fields{
		"clickInfo": click,
	}).Info("Insert click info")
	api.ResponseJSON(c, http.StatusOK, api.SUCCESS)
}

func buildSourceURL(click entity.ClickInfo) string {
	var sb strings.Builder
	sb.WriteString("appKey=")
	sb.WriteString(click.AppKey)
	sb.WriteString("&channelId=")
	sb.WriteString(click.ChannelId)
	sb.WriteString("&deviceId=")
	sb.WriteString(click.DeviceId)
	sb.WriteString("&clickTime=")
	sb.WriteString(strconv.FormatInt(click.ClickTime, 10))
	return sb.String()
}
