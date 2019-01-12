package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/simplexity-ckcclc/gochannel/api/errors"
	"github.com/simplexity-ckcclc/gochannel/api/handlers"
	"github.com/simplexity-ckcclc/gochannel/api/util"
	"net/http"
	"strings"
	"time"
)

var (
	token = ""
	publicKey = ""
	nonces = make(map[string]time.Time)
)

func Router() *gin.Engine {
	//router := gin.Default()
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	adGroup := router.Group("/ad")
	{
		adGroup.POST("/click", handlers.ClickHandler)
	}

	internalGroup := router.Group("/internal")
	internalGroup.Use(authInternal())
	internalGroup.Use()
	{
		internalGroup.POST("/appkey/:appkey/evict", handlers.EvictAppKeyHandler)
		internalGroup.POST("/appkey/:appkey/register", handlers.RegisterAppKeyHandler)
	}
	return router
}

func authInternal() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("auth internal")
		nonce := c.Query("nonce")
		sig := c.Query("sig")
		if nonce == "" || sig == "" {
			c.JSON(http.StatusOK, errors.REQUIRED_PARAMETER_MISSING)
			return
		}

		sourceURL := buildSourceText(nonce)
		if valid := validateNonce(nonce); !valid {
			c.JSON(http.StatusOK, errors.DUPLICATE_NONCE)
			return
		}

		valid, err := util.VerifySig(sourceURL, publicKey, sig)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.INTERNAL_SERVER_ERROR)
			return
		} else if !valid {
			c.JSON(http.StatusOK, errors.SIG_INVALID)
			return
		}

		c.Next()
	}
}

func validateNonce(nonce string) bool {
	// evict all expired nonce, inefficient, but works since there is little internal request
	now := time.Now()
	for nonce, timestamp := range nonces {
		if now.Sub(timestamp) > time.Duration(5 * time.Minute) {
			delete(nonces, nonce)
		}
	}

	if _, found := nonces[nonce]; found {
		return false
	}
	nonces[nonce] = time.Now()
	return true
}

func buildSourceText(nonce string) string {
	var sb strings.Builder
	sb.WriteString(token)
	sb.WriteString(":")
	sb.WriteString(nonce)
	return sb.String()
}


