package initialize

/**
*@Author lyer
*@Date 2/20/21 15:22
*@Describe
**/
import (
	"github.com/biningo/boil-gin/middleware"
	"github.com/biningo/boil-gin/router"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.Cors())

	PublicGroup := r.Group("")
	{
		router.InitUserRouter(PublicGroup)
		router.InitTagRouter(PublicGroup)
		router.InitBoilRouter(PublicGroup)
	}

	PrivateGroup := r.Group("")
	PrivateGroup.Use(middleware.JWTAuth())
	{
		router.InitPingRouter(PrivateGroup)
	}
	return r
}
