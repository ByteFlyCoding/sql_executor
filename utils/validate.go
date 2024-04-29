package utils

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"github.com/xwb1989/sqlparser"
)

// QuerySqlValidate 校验SQL语句是否合法，并断言其是否为 SELECT 操作，若不是则抛出error异常
func QuerySqlValidate(sql string) error {

	// Parse 完整解析 SQL 并返回一个 Statement，若SQL不合法则返回err
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		return err
	}

	// 断言该SQL是否为SELECT操作，若不是SELECT操作则抛出异常
	if _, ok := stmt.(*sqlparser.Select); !ok {
		return fmt.Errorf("该SQL不是SELECT操作")
	}

	return err
}

// ModifySqlValidate 校验SQL语句是否合法，并断言其是否为 DELETE、INSERT、UPDATE 操作，若不是则抛出error异常
func ModifySqlValidate(sql string) error {

	// Parse 完整解析 SQL 并返回一个 Statement，若SQL不合法则返回err
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		return err
	}

	// 断言该SQL是否为 DELETE、INSERT、UPDATE 操作，若不是则抛出error异常
	switch stmt.(type) {
	case *sqlparser.Insert:
		return nil
	case *sqlparser.Delete:
		return nil
	case *sqlparser.Update:
		return nil
	default:
		return fmt.Errorf("该sql不是SELECT语句")
	}

}

// ReturnModifyParamError 返回查询接口执行任务异常信息
func ReturnModifyParamError(code int, err interface{}) *ModifyParamError {

	var msg string
	switch err.(type) {
	case string:
		msg, _ = err.(string)
	default:
		msg = fmt.Sprintf("%s", err)
	}

	jsonData := ModifyParamError{Code: code, ErrMsg: msg}

	return &jsonData
}

// TransactionsValidate 修改接口输入参数校验
func TransactionsValidate(req *RequestBody) (*ModifyParamErrorJson, error) {

	defer func() {
		if err := recover(); err != nil {
			// 若发生panic()则捕获异常并打印日志，使程序继续执行而不退出
			logs.Error(err)
		}
	}()

	var err error
	var rsp ModifyParamErrorJson
	if len(req.Transactions) == 0 {
		// 请求中输入的事务列表为空，直接返回参数错误异常
		rsp.Code = PARAMETERERROR
		err = fmt.Errorf("没有输入任何事务")
		rsp.ErrMsg = err.Error()
		return &rsp, err
	}

	// 遍历事务列表
	for _, trsInfo := range req.Transactions {
		sqlErrorInfo := make([]SqlErrorInfo, 0)
		if len(trsInfo.Sqls) == 0 {
			// 该事务中的SQL列表为空，直接返回参数错误
			sqlErrorInfo = append(sqlErrorInfo, SqlErrorInfo{
				ErrMsg: fmt.Sprintf("事务%v中没有输入任何sql", *trsInfo),
			})
			logs.Error("该事务%v中没有输入任何sql语句", *trsInfo)
		}
		for _, info := range trsInfo.Sqls {
			err = ModifySqlValidate(info.Sql)
			if err != nil {
				logs.Error(err)
				sqlErrorInfo = append(sqlErrorInfo, SqlErrorInfo{
					ID:     info.ID,
					Sql:    info.Sql,
					Name:   info.Name,
					ErrMsg: err.Error(),
				})
			}
		}
		if len(sqlErrorInfo) > 0 {
			// 该事务存在异常，添加进返回列表
			rsp.Items = append(rsp.Items, TransactionParamError{
				ID:           trsInfo.ID,
				ErrMsg:       "事务没有输入SQL或输入的SQL中有语法错误",
				Count:        int64(len(sqlErrorInfo)),
				SqlErrorInfo: sqlErrorInfo,
			})
		}
	}
	rsp.Count = len(rsp.Items)
	if rsp.Count > 0 {
		// 输入参数存在异常
		rsp.Code = PARAMETERERROR
		err = fmt.Errorf("事务没有输入SQL或输入的SQL中有语法错误")
		rsp.ErrMsg = err.Error()
	}

	return &rsp, err
}
