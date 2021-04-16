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
		BoilRouter.GET("/all", v1.BoilAll)
		BoilRouter.GET("/list/tag/:tid",v1.BoilListByTag)
		BoilRouter.GET("/list/user/:uid",v1.BoilListByUser)
		BoilRouter.GET("/user/likes/:uid",v1.BoilListUserLike)
		BoilRouter.DELETE("/:bid", v1.BoilDelete)
		BoilRouter.GET("/id/:bid",v1.GetBoil)
		BoilRouter.GET("/list/user/:uid/comment",v1.BoilListUserComment)
	}

}
