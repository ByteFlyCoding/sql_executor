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
		panic("请确认数据库用户名设置")
	}
	password, err := beego.AppConfig.String("password")
	if err != nil {
		panic("请确认数据库用户密码设置")
	}
	localHost, err := beego.AppConfig.String("localHost")
	if err != nil {
		panic("请确认数据库主机名设置")
	}
	port, err := beego.AppConfig.String("port")
	if err != nil {
		panic("请确认数据库端口设置")
	}
	dbName, err := beego.AppConfig.String("dbName")
	if err != nil {
		panic("请确认数据库名设置")
	}
	param, err := beego.AppConfig.String("param")
	if err != nil {
		panic("请确认数据库连接参数设置")
	}

	dataSourceName := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?%v", userName, password, localHost, port, dbName, param)

	// 设置默认数据库
	if err := orm.RegisterDataBase("default", "mysql", dataSourceName); err != nil {
		logs.Error("请确认数据库状态、网络连接或数据库设置")
	}

}

func main() {

	// 生成 beego server 生命周期管理器
	lifeManager := life.NewLifeManager()

	// 生成 Executor 负责执行事务
	model := models.NewExecutor()

	// 注册路由
	routers.RegisterRouter(lifeManager, model)

	// 启动 beego server 生命周期管理器
	if err := life.Run(lifeManager); err != nil {
		log.Fatalln(err)
	}
}
