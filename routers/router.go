package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"sql_executor/controllers"
	"sql_executor/life"
	"sql_executor/models"
)

func RegisterRouter(manager *life.Manager, model *models.Executor) {

	executorCtl := &controllers.SqlExecutorController{
		Lmg:   manager,
		Model: model,
	}

	beego.Router("/sql_executor/query", executorCtl, "get:Query")
	beego.Router("/sql_executor/Modify", executorCtl, "post:Modify")
}
