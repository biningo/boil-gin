package initialize

import (
	"database/sql"
	"github.com/biningo/boil-gin/global"
	"github.com/jmoiron/sqlx"
	"time"
)

/**
*@Author lyer
*@Date 2/20/21 15:23
*@Describe
**/

import (
	"github.com/go-sql-driver/mysql"
)

func InitDB() *sqlx.DB {
	dbx := sqlx.NewDb(initMySql(), "mysql")
	return dbx
}

func initMySql() *sql.DB {
	connector, _ := mysql.NewConnector(&mysql.Config{
		Addr:      global.G_CONFIG.MySql.Host + ":" + global.G_CONFIG.MySql.Port,
		User:      global.G_CONFIG.MySql.User,
		Passwd:    global.G_CONFIG.MySql.Password,
		DBName:    global.G_CONFIG.MySql.DB,
		Collation: global.G_CONFIG.MySql.Collation,
		ParseTime: true,
		Loc:       time.Local,
	})
	db := sql.OpenDB(connector)
	if err := db.Ping(); err != nil {
		panic(err)
	}
	return db
}
