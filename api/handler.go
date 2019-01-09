package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const DSN = "ckcclc:141421@tcp(localhost:3306)/gochannel"

type clickInfo struct {
	AppKey   string
	DeviceId string
}

// The request responds to a url matching:  /ad/click?appKey=foo&deviceId=bar
func clickHandler(c *gin.Context) {
	appKey := c.Query("appKey")
	deviceId := c.Query("deviceId")

	click := clickInfo{
		AppKey:   appKey,
		DeviceId: deviceId,
	}

	db, err := open(DSN)
	if err != nil {
		c.JSON(http.StatusOK, INTERNAL_SERVER_ERROR)
	}

	if _, err := insertClickInfo(db, click); err != nil {
		c.JSON(http.StatusOK, INTERNAL_SERVER_ERROR)
	}
	c.JSON(http.StatusOK, SUCCESS)
}
