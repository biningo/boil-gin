package main

import (
	"github.com/biningo/boil-gin/global"
	"github.com/biningo/boil-gin/initialize"
	"github.com/biningo/boil-gin/server"
	"time"
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
	global.RedisClient = initialize.InitRedis()
	initialize.InitRedisToMySqlCron(time.Hour * 2)
	server.RunServer(":8080", global.Routers)
}
