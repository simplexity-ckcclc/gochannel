package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/simplexity-ckcclc/gochannel/api/appchannel"
	"github.com/simplexity-ckcclc/gochannel/api/click"
	api "github.com/simplexity-ckcclc/gochannel/api/common"
	"github.com/simplexity-ckcclc/gochannel/common"
	"github.com/simplexity-ckcclc/gochannel/common/logger"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

// The request responds to a url matching:  /ad/click?appKey=foo&deviceId=bar&sig=
func ClickHandler(c *gin.Context) {
	var click click.ClickInfo
	if err := c.ShouldBind(&click); err != nil {
		api.ResponseJSONWithExtraMsg(c, http.StatusBadRequest, api.RequiredParameterError, "Parameter missing")
		return
	}

	click.OsType = strings.ToLower(click.OsType)
	if ot := common.ParseOsType(click.OsType); ot == common.Unknown {
		api.ResponseJSONWithExtraMsg(c, http.StatusBadRequest, api.RequiredParameterError, "OsType invalid")
		return
	}

	appChannel, found := appchannel.SearchAppChannel(click.AppKey, click.ChannelId)
	if !found {
		api.ResponseJSON(c, http.StatusOK, api.ChannelNotFound)
		return
	}

	if appChannel.AppKey != click.AppKey {
		api.ResponseJSON(c, http.StatusOK, api.AppChannelUnmatched)
		return
	}

	//verify app key sig
	sig := c.Query("sig")
	if sig == "" {
		api.ResponseJSONWithExtraMsg(c, http.StatusBadRequest, api.RequiredParameterError, "sig is blank")
		return
	}
	sourceURL := buildSourceURL(click)
	valid, err := api.VerifyBase64WithRSAPubKey(sourceURL, appChannel.PublicKey, sig)
	if err != nil {
		logger.ApiLogger.WithFields(logrus.Fields{
			"pubKey": appChannel.PublicKey,
		}).Error("Verify click signature - VerifyBase64WithRSAPubKey error : ", err)
		api.ResponseJSONWithExtraMsg(c, http.StatusInternalServerError, api.InternalServerError, err.Error())
		return
	} else if !valid {
		logger.ApiLogger.With(logger.Fields{
			"clickInfo": click,
			"pubKey":    appChannel.PublicKey,
			"sig":       sig,
		}).Info("Invalid signature")
		api.ResponseJSON(c, http.StatusOK, api.SigInvalid)
		return
	}

	if err := click.InsertDB(common.DB); err != nil {
		logger.ApiLogger.With(logger.Fields{
			"clickInfo": click,
			"pubKey":    appChannel.PublicKey,
			"sig":       sig,
		}).Error("Insert DB error : ", err)
		api.ResponseJSONWithExtraMsg(c, http.StatusInternalServerError, api.InternalServerError, err.Error())
		return
	}

	logger.ApiLogger.With(logger.Fields{
		"clickInfo": click,
	}).Info("Insert click info")
	api.ResponseJSON(c, http.StatusOK, api.Success)
}

func buildSourceURL(click click.ClickInfo) string {
	return strings.Join([]string{click.AppKey, click.ChannelId, click.DeviceId, strconv.FormatInt(click.ClickTime, 10)}, "-")
}
