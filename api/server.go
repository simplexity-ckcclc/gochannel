package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/simplexity-ckcclc/gochannel/api/handlers"
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
	}
	return router
}

func authInternal() gin.HandlerFunc {
	return func(context *gin.Context) {
		fmt.Println("auth internal")
		context.Next()
	}
}
