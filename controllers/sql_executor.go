package controllers

import (
	"errors"
	"sync"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"sql_executor/life"
	"sql_executor/models"
	"sql_executor/utils"
)

type SqlExecutorController struct {
	beego.Controller
	Lmg   *life.Manager    // 用于启动beego服务，并确保程序接收到退出信号后所有任务运行完成后才退出程序
	Model *models.Executor // 存储beego orm 实例
}

// Query 查询接口 传入一个 SELECT 查询语句，执行查询任务，并将查询结果返回
func (c *SqlExecutorController) Query() {

	var msg string
	// 获取查询允许的重试次数
	retryCount, err := c.GetInt("retry")
	if err != nil {
		msg = "retryCount input is abnormal"
	}
	// 若接收到重试次数 count < 0 则默认不重试 即只执行一次
	if retryCount <= 0 {
		retryCount = 0
	}

	// 用于获取查询sql
	sql := c.GetString("sql")
	err = utils.QuerySqlValidate(sql) // sql 合法性校验
	if err != nil {
		// 若查询 SQL 不合法返回异常信息
		logs.Error("请求查询的SQL: %v 存在语法错误", err)
		c.Data["json"] = utils.ReturnQueryError(utils.FAILQUERY, sql, err)
		_ = c.ServeJSON()
		return
	}

	// 执行查询任务
	count, retryCount, items, err := c.Model.Query(sql, retryCount)
	if err != nil {
		msg = msg + " and " + err.Error()
	}

	// 生成查询接口返回对象
	c.Data["json"] = utils.ReturnQuerySuccess(sql, msg, items, count, retryCount)

	_ = c.ServeJSON()
}

// Modify 修改接口 执行 UPDATE、DELETE、INSERT 语句，一次接收一条或多条 SQL，并返回 SQL 的执行结果；
// 这些接口接收到的 SQL 语句，最终传递到后台数据库执行
func (c *SqlExecutorController) Modify() {

	req := new(utils.RequestBody)
	err := c.BindJSON(req)
	if err != nil {
		// 返回请求参数异常信息
		c.Data["json"] = utils.ReturnModifyParamError(utils.FAILQUERY, err)
		_ = c.ServeJSON()
		return
	}

	// 输入参数合法性校验
	resp, err := utils.TransactionsValidate(req)
	if err != nil {
		// 返回不合法的SQL语句和异常信息
		c.Data["json"] = resp
		_ = c.ServeJSON()
		return
	}

	// 运行 UPDATE、DELETE、INSERT 语句任务
	c.Data["json"] = modifyTaskRunners(req, c.Model)

	_ = c.ServeJSON()
}

// modifyTaskRunners 并发执行不同事务
func modifyTaskRunners(task *utils.RequestBody, model *models.Executor) *utils.ModifyJson {

	// 用于存储执行结果
	res := utils.ModifyJson{
		Code:   utils.SUCCESSMODIFY,
		Items:  make([]utils.Runner, 0),
		Count:  len(task.Transactions),
		ErrMsg: "所有都任务执行成功",
	}

	runner := func(t *utils.TransactionParamInfo, runnerInfo *utils.Runner) error {
		// 传入子任务信息调用数据库执行修改任务
		err := model.Modify(t, runnerInfo)
		if err != nil {
			return err
		}

		return nil
	}

	var m sync.Mutex       // 用于确保数据并发安全 后期可以考虑改成原子操作以提高并发性能
	wg := sync.WaitGroup{} // 用于保证所有子任务协程退出后modifyTaskRunners()才能退出 否则会panic引起程序退出
	runnerLogic := func(t *utils.TransactionParamInfo, taskInfo *utils.Runner) {

		defer func() {
			if err := recover(); err != nil {
				// 若发生panic()则捕获异常并打印日志，使程序继续执行而不退出
				logs.Error(err)
			}
		}()

		// 子任务调用逻辑
		defer wg.Done() // 子任务完成时，WaitGroup 计数器递减 1

		err := runner(t, taskInfo)
		// 若数据库执行任务返回错误，则根据设置的允许重试次数重试任务
		for err != nil && !errors.Is(err, models.ERROUTRETRYTIME) {
			err = runner(t, taskInfo)
			res.Code = utils.FAILMODIFYEXIST
			res.ErrMsg = "存在未执行成功的子任务，请排查，" + err.Error()
		}

		m.Lock() // 加锁保证并发安全，后期可以看看能不能改成原子操作
		res.Items = append(res.Items, *taskInfo)
		m.Unlock()
	}

	for _, t := range task.Transactions {
		// 根据输入的任务设置子任务执行信息
		result := new(utils.Runner)
		result.ID = t.ID
		result.Name = t.Name
		result.Count = int64(len(t.Sqls))

		wg.Add(1)                 // 启动子任务执行逻辑前 WaitGroup 计数器递减 1 确保子任务完成前 modifyTaskRunners 被阻塞
		go runnerLogic(t, result) // 启动子任务调用逻辑
	}

	wg.Wait() // 确保子任务完成前 modifyTaskRunners 不会提前退出

	return &res
}

func (c *SqlExecutorController) Prepare() {
	c.Lmg.WaitAdd() // 每传入查询请求或修改请求 life.Manager 中的WaitGrout计数器都加一
	logs.Info("查询或修改请求传入")
}

// Finish 释放资源
func (c *SqlExecutorController) Finish() {
	// 程序接收退出信号后阻塞住程序，防止查询任务或修改任务结束前退出程序
	c.Lmg.WaitDone() // 查询请求或修改请求结束 life.Manager 中的WaitGrout计数器都加一
	logs.Info("请求执行结束")
}
