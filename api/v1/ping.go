package v1

import (
	"github.com/biningo/boil-gin/global"
	"github.com/gin-gonic/gin"
)

/**
*@Author lyer
*@Date 4/1/21 14:03
*@Describe
**/

func PingApi() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, global.G_CONFIG)
	}
}
