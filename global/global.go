package global

/**
*@Author lyer
*@Date 2/20/21 15:13
*@Describe
**/

import (
	"github.com/biningo/boil-gin/config"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

var (
	G_VP     *viper.Viper
	G_CONFIG config.Config
	G_DB     *sqlx.DB
	Routers  *gin.Engine
	RedisClient		*redis.Client
)
