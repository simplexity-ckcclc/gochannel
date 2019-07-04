package api

import (
	"github.com/gin-gonic/gin"
	api "github.com/simplexity-ckcclc/gochannel/api/common"
	"github.com/simplexity-ckcclc/gochannel/api/entity"
	"github.com/simplexity-ckcclc/gochannel/api/errorcode"
	"github.com/simplexity-ckcclc/gochannel/api/handlers"
	"github.com/simplexity-ckcclc/gochannel/common"
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

	if err := entity.LoadAppKeySigs(common.DB); err != nil {
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
		internalGroup.POST("/appkey/:appkey/evict", handlers.EvictAppKeyHandler)
		internalGroup.POST("/appkey/:appkey/register", handlers.RegisterAppKeyHandler)
	}
	return router
}

func authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		nonce := c.Query("nonce")
		sig := c.Query("sig")
		if nonce == "" || sig == "" {
			c.JSON(http.StatusBadRequest, errorcode.REQUIRED_PARAMETER_MISSING)
			c.Abort()
			return
		}

		if valid := validateNonce(nonce); !valid {
			c.JSON(http.StatusOK, errorcode.DUPLICATE_NONCE)
			c.Abort()
			return
		}

		//valid, err := api.VerifyBase64WithRSAPubKey(nonce, common.Conf.Api.Internal.PublicKey, sig)
		//if err != nil {
		//	api.ApiLog.WithFields(logrus.Fields{
		//		"pubKey": common.Conf.Api.Internal.PublicKey,
		//	}).Error("Verify internal signature - VerifyBase64WithRSAPubKey error : ", err)
		//	c.JSON(http.StatusInternalServerError, errorcode.INTERNAL_SERVER_ERROR)
		//	c.Abort()
		//	return
		//} else if !valid {
		//	c.JSON(http.StatusOK, errorcode.SIG_INVALID)
		//	c.Abort()
		//	return
		//}

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
