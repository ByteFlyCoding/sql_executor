package utils

import (
	"errors"
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"github.com/xwb1989/sqlparser"
)

// SqlValidate 这个校验逻辑还可以加上判断sql语句是哪种类型的
func SqlValidate(sql string) error {

	if sql == "" {
		return errors.New("sql is empty")
	}

	_, err := sqlparser.Parse(sql)
	if err != nil {
		return err
	}

	return nil
}

// ModifySqlValidate 修改接口sql数据校验
func ModifySqlValidate(req *RequestBody) (*ModifySqlErrorJson, error) {

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
		}
	}()

	// sql合法性校验
	// 可以再加入类型的判断
	rsp := ModifySqlErrorJson{
		Code:   FAILMODIFY,
		ErrMsg: "syntax error hava exist in these sqls",
		// Items:  make([]SyntaxError, 0),
		Items: make([]Runner, 0),
	}
	var err error
	for _, trsInfo := range req.Transactions {
		transactionError := make([]SqlInfo, 0)
		for _, info := range trsInfo.Sqls {
			err = SqlValidate(info.Sql)
			if err != nil {
				logs.Error("some sql syntax error have exist in the sql: %v", err)
				transactionError = append(transactionError, SqlInfo{
					ID:     info.ID,
					Sql:    info.Sql,
					Name:   info.Name,
					ErrMsg: err.Error(),
				})
			}
		}
		if len(transactionError) > 0 {
			rsp.Items = append(rsp.Items, Runner{
				ID:      trsInfo.ID,
				ErrMsg:  "Syntax error has exist in the Transaction sql",
				Count:   int64(len(transactionError)),
				SqlInfo: transactionError,
			})
		}
	}
	rsp.Count = len(rsp.Items)
	if rsp.Count > 0 {

		return &rsp, fmt.Errorf("syntax error has exist in these sqls")
	}

	return &rsp, err
}
