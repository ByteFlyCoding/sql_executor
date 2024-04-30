package utils

import (
	"fmt"

	"github.com/beego/beego/v2/client/orm"
)

// QueryErrorJson 用于序列化查询接口异常信息
type QueryErrorJson struct {
	Code   int    `json:"code"` // 业务状态码
	Sql    string `json:"sql"`
	ErrMsg string `json:"err_msg"` // 异常信息
}

// ReturnQueryError 返回查询接口执行任务异常信息
func ReturnQueryError(code int, sql string, err interface{}) *QueryErrorJson {

	var msg string
	switch err.(type) {
	case string:
		msg, _ = err.(string)
	default:
		msg = fmt.Sprintf("%s", err)
	}

	jsonData := QueryErrorJson{Code: code, Sql: sql, ErrMsg: msg}

	return &jsonData
}

// QuerySuccessJson 用于序列化查询接口任务执行成功信息
type QuerySuccessJson struct {
	Code   int          `json:"code"`            // 业务状态码
	Sql    string       `json:"sql"`             // SQL语句
	Count  int64        `json:"count"`           // 查询结果记录数
	Items  []orm.Params `json:"items,omitempty"` // 查询结果，以列表形式存储
	Retry  int          `json:"retry"`           // 重试次数
	ErrMsg string       `json:"err_msg"`         // 查询成功消息
}

// ReturnQuerySuccess 返回查询接口执行任务成功信息
func ReturnQuerySuccess(sql string, msg string, items []orm.Params, count int64, retry int) *QuerySuccessJson {

	jsonData := QuerySuccessJson{
		Code:   SUCCESSQUERY,
		Sql:    sql,
		Count:  count,
		Items:  items,
		Retry:  retry,
		ErrMsg: msg,
	}

	return &jsonData
}
