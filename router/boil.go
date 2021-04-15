package router

import (
	v1 "github.com/biningo/boil-gin/api/v1"
	"github.com/gin-gonic/gin"
)

/**
*@Author lyer
*@Date 4/15/21 13:58
*@Describe
**/
func InitBoilRouter(Router *gin.RouterGroup) {
	BoilRouter := Router.Group("/boil")
	{
		BoilRouter.POST("/publish", v1.BoilPublish)
		BoilRouter.DELETE("/:id", v1.BoilDelete)
		BoilRouter.GET("/like", v1.BoilLike)
		BoilRouter.GET("/unlike", v1.BoilUnLike)
		BoilRouter.GET("/all", v1.BoilAll)
		BoilRouter.GET("/list/user/:uid", v1.BoilUserList)
		BoilRouter.GET("/list/user/like/:uid", v1.BoilUserLikeList)
	}

}
