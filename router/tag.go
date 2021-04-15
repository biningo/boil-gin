package router

import (
	v1 "github.com/biningo/boil-gin/api/v1"
	"github.com/gin-gonic/gin"
)

/**
*@Author lyer
*@Date 4/15/21 15:32
*@Describe
**/

func InitTagRouter(Router *gin.RouterGroup) {
	TagRouter := Router.Group("/tag")
	{
		TagRouter.GET("/list", v1.TagList)
		TagRouter.GET("/create/:title", v1.TagCreate)
		TagRouter.GET("/delete/:tid", v1.TagDelete)
		TagRouter.GET("/boil/count/:tid", v1.TagBoilCount)
	}
}
