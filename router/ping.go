package router

import (
	v1 "github.com/biningo/boil-gin/api/v1"
	"github.com/gin-gonic/gin"
)

/**
*@Author lyer
*@Date 2/20/21 18:23
*@Describe
**/

func InitPingRouter(Router *gin.RouterGroup) {
	PingRouter := Router.Group("/ping")
	{
		PingRouter.GET("/config", v1.PingApi())
	}
}
