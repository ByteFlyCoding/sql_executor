package models

import (
	"context"
	"fmt"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/pkg/errors"
	"sql_executor/utils"
)

var ERROUTRETRYTIME = errors.New("超过最大重试次数")

type Executor struct {
	orm.Ormer
}

func NewExecutor() *Executor {
	// 构造执行器
	return &Executor{
		orm.NewOrm(),
	}
}

// Query 执行传入的SELECT语句
func (e *Executor) Query(sql string, retry int) (int64, int, []orm.Params, error) {

	var err error
	var count int64
	result := make([]orm.Params, 0)

	var i int
	for i = 0; i <= retry; i++ {
		// 若发生错误则重试 重试次数不能超过最大允许重试次数
		count, err = e.Raw(sql).Values(&result)
		if err == nil {
			return count, i, result, nil
		}
	}

	return count, i, result, fmt.Errorf("超过最大允许重试次数：%v 查询失败", retry)
}

// Modify 执行传入的事务
func (e *Executor) Modify(t *utils.TransactionParamInfo, runner *utils.Runner) error {

	// 超过最大允许重试次数，返回异常
	if t.Retry+1 < runner.Retry {
		return ERROUTRETRYTIME
	}
	runner.Retry++

	if t.Timeout <= 0 {
		runner.Timeout = 300
	} else {
		runner.Timeout = t.Timeout
	}

	ctx, cancel := context.WithTimeout(context.Background(), runner.Timeout*time.Second) // 设置超时, 默认为300s
	defer cancel()                                                                       // 确保在函数退出时取消上下文

	// 开启事务
	tx, err := e.BeginWithCtx(ctx)
	if err != nil {
		runner.ErrMsg = err.Error()
		return err
	}

	for _, sqlInfo := range t.Sqls {
		// 该事务的执行信息
		execInfo := utils.SqlExecInfo{
			ID:   sqlInfo.ID,
			Name: sqlInfo.Name,
			Sql:  sqlInfo.Sql,
		}

		// 执行同一个事务中的INSERT、DELETE或者UPDATE语句
		result, err := e.Raw(sqlInfo.Sql).Exec()
		if err != nil {
			// 若有SQL执行出错则退出
			errRollback := tx.Rollback()
			if errRollback != nil {
				runner.ErrMsg = "事务回滚失败，" + err.Error()
			}
			execInfo.ErrMsg = "事务执行失败，已经回滚，" + err.Error()
			runner.SqlExecInfo = append(runner.SqlExecInfo, execInfo)
			return err
		}
		if result != nil {
			count, err := result.RowsAffected()
			if err != nil {
				// 这里不一定是错误，可能是没有行被修改（noLows），由不同的数据库驱动实现决定
				execInfo.ErrMsg = err.Error()
				runner.SqlExecInfo = append(runner.SqlExecInfo, execInfo)
			}
			if count > 0 {
				// 记录该 SQL 语句执行后生效的记录数、执行成功消息
				execInfo.Count = count
				execInfo.ErrMsg = "该SQL执行成功，等待提交"
				runner.SqlExecInfo = append(runner.SqlExecInfo, execInfo)
			}
		}
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		// 事务提交失败 输出日志
		runner.ErrMsg = "该事务提交失败" + err.Error()
		logs.Error(err)
		return err
	}

	runner.ErrMsg = "事务提交成功"

	return nil
}
