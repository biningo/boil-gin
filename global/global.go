package global

/**
*@Author lyer
*@Date 2/20/21 15:13
*@Describe
**/

import (
	"database/sql"
	"github.com/biningo/boil-gin/config"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	G_VP     *viper.Viper
	G_CONFIG config.Config
	G_DB     *sql.DB
	Routers  *gin.Engine
)
