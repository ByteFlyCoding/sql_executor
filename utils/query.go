package utils

import (
	"fmt"

	"github.com/beego/beego/v2/client/orm"
)

const (
	SUCCESSQUERY = iota
	FAILQUERY
	SUCCESSMODIFY
	FAILMODIFY
	FAILMODIFYEXIST
)

type ReturnQueryErrorJson struct {
	Code   int    `json:"code"`
	ErrMsg string `json:"err_msg"`
}

func ReturnQueryError(code int, err interface{}) *ReturnQueryErrorJson {

	var msg string
	switch err.(type) {
	case string:
		msg, _ = err.(string)
	default:
		msg = fmt.Sprintf("%s", err)
	}

	jsonData := ReturnQueryErrorJson{Code: code, ErrMsg: msg}

	return &jsonData
}

type ReturnQuerySuccessJson struct {
	Code  int          `json:"code"`
	Sql   string       `json:"sql"`
	Count int64        `json:"count"`
	Items []orm.Params `json:"items"`
	Retry int          `json:"retry"`
	Msg   string       `json:"msg"`
}

func ReturnQuerySuccess(sql string, msg string, items []orm.Params, count int64, retry int) *ReturnQuerySuccessJson {

	jsonData := ReturnQuerySuccessJson{
		Code:  SUCCESSQUERY,
		Sql:   sql,
		Count: count,
		Items: items,
		Retry: retry,
		Msg:   msg,
	}

	return &jsonData
}
