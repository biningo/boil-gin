package router

import (
	v1 "github.com/biningo/boil-gin/api/v1"
	"github.com/gin-gonic/gin"
)

/**
*@Author lyer
*@Date 4/13/21 15:13
*@Describe
**/
func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("/user")
	{
		UserRouter.GET("/all", v1.ListUser)
		UserRouter.POST("/login", v1.Login)
		UserRouter.GET("/logout", v1.Logout)
		UserRouter.POST("/registry", v1.Registry)
		UserRouter.GET("/info/:uid", v1.UserInfo)
		UserRouter.POST("/update/bio/:uid", v1.UpdateUserBio)
		UserRouter.GET("/follow/:uid", v1.UserFollow)
		UserRouter.GET("/unfollow/:uid", v1.UserUnFollow)
		UserRouter.GET("/list/follower/:uid", v1.ListUserFollower)
		UserRouter.GET("/list/following/:uid", v1.ListUserFollowing)
		UserRouter.GET("/list/recommend", v1.ListRecommendUser)
	}
}
