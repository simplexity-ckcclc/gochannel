package api

import (
	"github.com/gin-gonic/gin"
	"github.com/simplexity-ckcclc/gochannel/api/appchannel"
	"github.com/simplexity-ckcclc/gochannel/api/click"
	api "github.com/simplexity-ckcclc/gochannel/api/common"
	"github.com/simplexity-ckcclc/gochannel/api/handlers"
	"github.com/simplexity-ckcclc/gochannel/common"
	"github.com/simplexity-ckcclc/gochannel/common/config"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	nonceTTL = 5
)

var (
	nonces = make(map[string]time.Time)
)

func Serve() {

	if err := appchannel.LoadAppChannels(common.DB); err != nil {
		panic(err)
	}

	clickPorter := click.NewClickPorter(common.DB, common.EsClient)
	go clickPorter.TransferClicks()

	server := &http.Server{
		Addr:    config.GetString(config.ApiServerAddress),
		Handler: router(),
	}
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
	common.ApiLogger.Info("API server started!")

}

func router() *gin.Engine {
	//router := gin.Default()
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	adGroup := router.Group("/ad")
	{
		adGroup.POST("/click", handlers.ClickHandler)
	}

	internalGroup := router.Group("/ad/internal")
	internalGroup.Use(authRequired())
	{
		internalGroup.POST("/channel/evict", handlers.EvictChannelHandler)
		internalGroup.POST("/channel/register", handlers.RegisterChannelHandler)
	}
	return router
}

func authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		nonce := c.Query("nonce")
		sig := c.Query("sig")
		if nonce == "" || sig == "" {
			api.ResponseJSON(c, http.StatusBadRequest, api.RequiredParameterError)
			c.Abort()
			return
		}

		if valid := validateNonce(nonce); !valid {
			api.ResponseJSON(c, http.StatusOK, api.DuplicateNonce)
			c.Abort()
			return
		}

		pubKey := config.GetString(config.ApiServerIntlPubKey)
		valid, err := api.VerifyBase64WithRSAPubKey(nonce, pubKey, sig)
		if err != nil {
			common.ApiLogger.WithFields(logrus.Fields{
				"pubKey": pubKey,
			}).Error("Verify internal signature - VerifyBase64WithRSAPubKey error : ", err)
			api.ResponseJSONWithExtraMsg(c, http.StatusInternalServerError, api.InternalServerError, err.Error())
			c.Abort()
			return
		} else if !valid {
			api.ResponseJSON(c, http.StatusOK, api.SigInvalid)
			c.Abort()
			return
		}

		c.Next()
	}
}

// in case for replay attack
func validateNonce(nonce string) bool {
	// evict all expired nonce, inefficient, but works well since there is little internal request
	now := time.Now()
	for nonce, timestamp := range nonces {
		if now.Sub(timestamp) > time.Duration(nonceTTL*time.Minute) {
			delete(nonces, nonce)
		}
	}

	if _, found := nonces[nonce]; found {
		return false
	}
	nonces[nonce] = time.Now()
	return true
}
