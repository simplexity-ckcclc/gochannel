package api

import (
	"github.com/gin-gonic/gin"
	api "github.com/simplexity-ckcclc/gochannel/api/common"
	"github.com/simplexity-ckcclc/gochannel/api/entity"
	"github.com/simplexity-ckcclc/gochannel/api/handlers"
	"github.com/simplexity-ckcclc/gochannel/common"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var (
	nonces = make(map[string]time.Time)
)

func Serve() {

	conf := common.Conf
	if err := api.InitLogger(conf); err != nil {
		panic(err)
	}

	db, err := api.OpenDB(common.Conf.Core.DSN)
	defer db.Close()
	if err != nil {
		panic(err)
	}

	if err := entity.LoadChannelSigs(api.DB); err != nil {
		panic(err)
	}

	server := &http.Server{
		Addr:    conf.Api.Address,
		Handler: router(),
	}
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
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

	internalGroup := router.Group("/internal")
	internalGroup.Use(authRequired())
	{
		internalGroup.POST("/channel/:channel/evict", handlers.EvictChannelHandler)
		internalGroup.POST("/channel/:channel/register", handlers.RegisterChannelHandler)
	}
	return router
}

func authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		nonce := c.Query("nonce")
		sig := c.Query("sig")
		if nonce == "" || sig == "" {
			api.ResponseJSON(c, http.StatusBadRequest, api.REQUIRED_PARAMETER_MISSING)
			c.Abort()
			return
		}

		if valid := validateNonce(nonce); !valid {
			api.ResponseJSON(c, http.StatusOK, api.DUPLICATE_NONCE)
			c.Abort()
			return
		}

		valid, err := api.VerifyBase64WithRSAPubKey(nonce, common.Conf.Api.Internal.PublicKey, sig)
		if err != nil {
			api.ApiLog.WithFields(logrus.Fields{
				"pubKey": common.Conf.Api.Internal.PublicKey,
			}).Error("Verify internal signature - VerifyBase64WithRSAPubKey error : ", err)
			api.ResponseJSONWithExtraMsg(c, http.StatusInternalServerError, api.INTERNAL_SERVER_ERROR, err.Error())
			c.Abort()
			return
		} else if !valid {
			api.ResponseJSON(c, http.StatusOK, api.SIG_INVALID)
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
		if now.Sub(timestamp) > time.Duration(5*time.Minute) {
			delete(nonces, nonce)
		}
	}

	if _, found := nonces[nonce]; found {
		return false
	}
	nonces[nonce] = time.Now()
	return true
}
