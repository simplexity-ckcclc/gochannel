package api

import (
	"github.com/gin-gonic/gin"
)

const DSN = "ckcclc:141421@/gochannel"

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
		panic(err)
	}

	if _, err := insertClickInfo(db, click); err != nil {
		panic(err)
	}
}
