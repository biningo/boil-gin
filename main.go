package main

import (
	"github.com/biningo/boil-gin/global"
	"github.com/biningo/boil-gin/initialize"
	"github.com/biningo/boil-gin/server"
)

/**
*@Author lyer
*@Date 2/20/21 15:18
*@Describe
**/

func main() {
	global.G_VP = initialize.InitViper()
	global.G_DB = initialize.InitDB()
	global.Routers = initialize.InitRouter()
	server.RunServer(":8080",global.Routers)
}
