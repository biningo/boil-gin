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
	BoilRouter := Router.Group("/boils")
	{
		BoilRouter.POST("/publish", v1.BoilPublish)
		BoilRouter.GET("/list", v1.BoilAll)
		BoilRouter.GET("/list/tag/:tid", v1.BoilListByTag)
		BoilRouter.GET("/list/user/:uid", v1.BoilListByUser)
		BoilRouter.GET("/list/user/:uid/like", v1.BoilListUserLike)
		BoilRouter.GET("/list/user/:uid/comment", v1.BoilListUserComment)
		BoilRouter.GET("/list/following", v1.BoilListByFollowing)
		BoilRouter.DELETE("/:bid", v1.BoilDelete)
		BoilRouter.GET("/boil/:bid", v1.GetBoilById)

		BoilRouter.GET("/user/:uid/like/:bid", v1.BoilUserLike)
		BoilRouter.GET("/user/:uid/unlike/:bid", v1.BoilUserUnLike)
		BoilRouter.GET("/user/:uid/islike/:bid", v1.BoilUserIsLike)
	}

}
