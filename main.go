package main

import (
	"fmt"
	"log"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
	"sql_executor/life"
	"sql_executor/models"
	"sql_executor/routers"
	_ "sql_executor/routers"
)

func init() {

	userName, err := beego.AppConfig.String("userName")
	if err != nil {
		panic("please check database user name config")
	}
	password, err := beego.AppConfig.String("password")
	if err != nil {
		panic("please check database user password config")
	}
	localHost, err := beego.AppConfig.String("localHost")
	if err != nil {
		panic("please check database local host config")
	}
	port, err := beego.AppConfig.String("port")
	if err != nil {
		panic("please check database port config")
	}
	dbName, err := beego.AppConfig.String("dbName")
	if err != nil {
		panic("please check database name config")
	}
	param, err := beego.AppConfig.String("param")
	if err != nil {
		panic("please check database user param config")
	}

	dataSourceName := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?%v", userName, password, localHost, port, dbName, param)

	// 设置默认数据库
	if err := orm.RegisterDataBase("default", "mysql", dataSourceName); err != nil {
		logs.Error("please check database status, network connection with database and database database config")
	}

}

func main() {

	lifeManager := life.NewLifeManager()

	model := models.NewExecutor()

	routers.RegisterRouter(lifeManager, model)

	if err := life.Run(lifeManager); err != nil {
		log.Fatalln(err)
	}
}
