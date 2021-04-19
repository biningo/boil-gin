package initialize

/**
*@Author lyer
*@Date 2/20/21 15:22
*@Describe
**/
import (
	"github.com/biningo/boil-gin/global"
	"github.com/biningo/boil-gin/middleware"
	"github.com/biningo/boil-gin/router"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	gin.SetMode(global.G_CONFIG.Server.Mode)
	r.Use(middleware.Cors())
	r.Use(middleware.UserInfoMiddleware())
	PublicGroup := r.Group("")
	{
		router.InitUserRouter(PublicGroup)
		router.InitTagRouter(PublicGroup)
		router.InitBoilRouter(PublicGroup)
		router.InitCommentRouter(PublicGroup)
	}

	PrivateGroup := r.Group("")
	PrivateGroup.Use(middleware.JWTAuth())
	{
		router.InitPingRouter(PrivateGroup)
	}
	return r
}
