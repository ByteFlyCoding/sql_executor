package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"sql_executor/life"
	"sql_executor/models"
	"sql_executor/utils"
)

type SqlExecutorController struct {
	beego.Controller
	Lmg   *life.Manager
	Model *models.Executor
}

// Query 查询接口
func (c *SqlExecutorController) Query() {

	var msg string
	retryCount, err := c.GetInt("retry")
	if err != nil {
		msg = "retryCount input is abnormal"
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
	}
	// 若接收到重试次数 count < 0 则默认不重试 即只执行一次
	if retryCount < 0 {
		retryCount = 0
	}

	sql := c.GetString("sql")
	err = utils.SqlValidate(sql) // sql 合法性校验
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		logs.Error("some sql syntax error have exist in the sql: %v", err)
		c.Data["json"] = utils.ReturnQueryError(utils.FAILQUERY, err)
		_ = c.ServeJSON()
		return
	}

	count, retryCount, items, err := c.Model.Query(sql, retryCount)
	if err != nil {
		msg = msg + " and " + err.Error()
	}

	// 执行查询任务
	c.Data["json"] = utils.ReturnQuerySuccess(sql, msg, items, count, retryCount)

	_ = c.ServeJSON()
}

// Modify 修改接口
func (c *SqlExecutorController) Modify() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
		}
	}()

	req := new(utils.RequestBody)
	err := c.BindJSON(req)
	if err != nil {
		c.Data["json"] = utils.ReturnQueryError(utils.FAILQUERY, err)
		_ = c.ServeJSON()
		return
	}

	// Sql校验器有bug 后续再改
	resp, err := utils.ModifySqlValidate(req)
	if err != nil {
		c.Data["json"] = resp
		_ = c.ServeJSON()
		return
	}

	c.Data["json"] = modifyTaskRunners(req, c.Model)

	_ = c.ServeJSON()
}

func modifyTaskRunners(task *utils.RequestBody, model *models.Executor) *utils.ModifySuccessJson {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	res := utils.ModifySuccessJson{
		Code:  utils.SUCCESSMODIFY,
		Items: make([]utils.Runner, 0),
		Count: len(task.Transactions),
	}

	runner := func(t *utils.TransactionInfo, runnerInfo *utils.Runner) error {

		err := model.Modify(t, runnerInfo)
		if err != nil {
			return err
		}

		return nil
	}

	var m sync.Mutex
	wg := sync.WaitGroup{}
	runnerLogic := func(t *utils.TransactionInfo, taskInfo *utils.Runner) {

		defer wg.Done()

		err := runner(t, taskInfo)
		for err != nil && !errors.Is(err, models.ERROUTRETRYTIME) {
			err = runner(t, taskInfo)
			res.Code = utils.FAILMODIFYEXIST
		}

		if err != nil {
			taskInfo.ErrMsg = err.Error()
		}

		m.Lock()
		res.Items = append(res.Items, *taskInfo)
		m.Unlock()
	}

	for _, t := range task.Transactions {

		result := new(utils.Runner)
		result.ID = t.ID
		result.Retry = -1
		result.Name = t.Name
		result.Count = int64(len(t.Sqls))
		result.SqlInfo = t.Sqls

		wg.Add(1)
		go runnerLogic(t, result)
	}

	wg.Wait()

	return &res
}

func (c *SqlExecutorController) Prepare() {
	c.Lmg.WaitAdd()
	logs.Info("request input")
}

// Finish 确保任务完成后退出
func (c *SqlExecutorController) Finish() {
	c.Lmg.WaitDone()
	logs.Info("finish request")
}
