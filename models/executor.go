package models

import (
	"fmt"

	"github.com/beego/beego/v2/client/orm"
	"github.com/pkg/errors"
	"sql_executor/utils"
)

var ERROUTRETRYTIME = errors.New("超过最大重试次数")

type Executor struct {
	orm.Ormer
}

func NewExecutor() *Executor {
	return &Executor{
		orm.NewOrm(),
	}
}

func (e *Executor) Query(sql string, retry int) (int64, int, []orm.Params, error) {

	var err error
	var count int64
	result := make([]orm.Params, 0)

	var i int
	for i = 0; i <= retry; i++ {
		count, err = e.Raw(sql).Values(&result)
		if err == nil {
			return count, i, result, nil
		}
	}

	return count, i, result, fmt.Errorf("超过最大允许重试次数：%v 查询失败", retry)
}

func (e *Executor) Modify(t *utils.TransactionInfo, runner *utils.Runner) error {

	if t.Retry < runner.Retry {
		return ERROUTRETRYTIME
	}

	runner.Retry++

	tx, err := e.Begin()
	if err != nil {
		runner.ErrMsg = err.Error()
		return err
	}

	for i, sqlInfo := range t.Sqls {
		result, err := e.Raw(sqlInfo.Sql).Exec()
		if err != nil {
			_ = tx.Rollback()
			runner.SqlInfo[i].ErrMsg = err.Error()
			return err
		}
		if result != nil {
			count, err := result.RowsAffected()
			if err != nil {
				// 这里不一定是错误，可能是没有行被修改（noLows）, 由不同的数据库驱动实现决定
				runner.SqlInfo[i].ErrMsg = err.Error()
			}
			if count > 0 {
				runner.SqlInfo[i].Count = count
			}
		}

	}

	runner.ErrMsg = "Transaction execute successfully"

	err = tx.Commit()
	if err != nil {
		runner.ErrMsg = err.Error()

		return err
	}

	return nil
}
