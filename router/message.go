package router

import (
	v1 "github.com/biningo/boil-gin/api/v1"
	"github.com/gin-gonic/gin"
)

/**
*@Author lyer
*@Date 4/15/21 14:55
*@Describe
**/

func InitMessageRouter(Router *gin.RouterGroup) {
	MsgRouter := Router.Group("/msssage")
	{
		MsgRouter.GET("/list/:uid", v1.MessageList)
		MsgRouter.POST("/create/:uid", v1.MessageCreate)
		MsgRouter.DELETE("/:mid", v1.MessageDelete)
	}
}
