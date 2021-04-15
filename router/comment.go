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
		CommentRouter.POST("/publish", v1.CommentPublish)
		CommentRouter.GET("/list/boil/:bid", v1.CommentBoilList)
		CommentRouter.GET("/comment/:cid", v1.CommentCommentList)
	}
}
