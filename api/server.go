package api

import (
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()
	adGroup := router.Group("/ad")
	{
		adGroup.POST("/click", clickHandler)
	}
	return router
}
