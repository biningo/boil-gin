package router

import (
	v1 "github.com/biningo/boil-gin/api/v1"
	"github.com/gin-gonic/gin"
)

/**
*@Author lyer
*@Date 4/15/21 14:00
*@Describe
**/
func InitCommentRouter(Router *gin.RouterGroup) {
	CommentRouter := Router.Group("/comment")
	{
		CommentRouter.POST("/publish/:bid", v1.CommentPublish)
		CommentRouter.GET("/list/:bid", v1.CommentBoilList)
		CommentRouter.DELETE("/:cid",v1.CommentDelete)
	}
}
