package utils

import "time"

// RequestBody 用于反序列化Modify接口请求的body
type RequestBody struct {
	Transactions []*TransactionParamInfo `json:"transactions"` // 事务列表信息参数列表
}

// TransactionParamInfo 存储修改接口传入的事务信息参数
type TransactionParamInfo struct {
	ID      int           `json:"id,omitempty"`      // 事务ID参数
	Retry   int           `json:"retry,omitempty"`   // 允许重试次数参数
	Timeout time.Duration `json:"timeout,omitempty"` // 该事务的超时时间，以秒为单位
	Name    string        `json:"name,omitempty"`    // 事务名称参数
	Sqls    []SqlInfo     `json:"sqls"`              // 事务SQL列表参数
}

// SqlInfo 执行信息参数
type SqlInfo struct {
	ID   int    `json:"id,omitempty"`   // SQL ID参数
	Name string `json:"name,omitempty"` // SQL 名称参数
	Sql  string `json:"sql"`            // SQL 语句参数
}

// ModifyParamErrorJson 用来序列化请求参数异常信息
type ModifyParamErrorJson struct {
	Code   int                     `json:"code"`            // 业务状态码
	Items  []TransactionParamError `json:"items,omitempty"` // 带有不合法SQL的事务列表
	Count  int                     `json:"count"`           // 带有不合法SQL的事务数量
	ErrMsg string                  `json:"err_msg"`         // 错误消息
}

// TransactionParamError 带有SQL语法错误的事务信息
type TransactionParamError struct {
	ID           int            `json:"id,omitempty"`      // 事务ID
	Count        int64          `json:"count"`             // 错误的SQL语句数量
	Timeout      time.Duration  `json:"timeout,omitempty"` // 事务超时时间，以秒为单位
	Name         string         `json:"name,omitempty"`    // 事务名称
	ErrMsg       string         `json:"err_msg"`           // 错误消息
	SqlErrorInfo []SqlErrorInfo `json:"items,omitempty"`   // 该事务中有语法错误的SQL列表
}

// SqlErrorInfo SQL语法错误信息
type SqlErrorInfo struct {
	ID     int    `json:"id,omitempty"`   // SQL ID
	Name   string `json:"name,omitempty"` // SQL 名称
	Sql    string `json:"sql"`            // 有语法错误的SQL 语句
	ErrMsg string `json:"err_msg"`        // SQL语法错误信息
}

// ModifyJson 存储Modify接口执行结果
type ModifyJson struct {
	Code   int      `json:"code"`    // 业务状态码
	Items  []Runner `json:"items"`   // 事务执行结果信息，以列表形式存储
	Count  int      `json:"count"`   // 事务数量
	ErrMsg string   `json:"err_msg"` // 事务执行信息
}

// Runner 存储事务执行信息
type Runner struct {
	ID          int           `json:"id,omitempty"`      // 事务ID
	Retry       int           `json:"retry"`             // 重试次数
	Count       int64         `json:"count"`             // 该事务中执行的SQL数量
	Name        string        `json:"name,omitempty"`    // 事务名称
	ErrMsg      string        `json:"err_msg,omitempty"` // 事务运行消息
	SqlExecInfo []SqlExecInfo `json:"items"`             // 该事务中的SQL运行情况列表
	Timeout     time.Duration `json:"timeout"`           // 该事务的超时数据，以秒为单位
}

// SqlExecInfo sql语句信息
type SqlExecInfo struct {
	ID     int    `json:"id,omitempty"`   // SQL ID
	Name   string `json:"name,omitempty"` // SQL名称
	Sql    string `json:"sql"`            // SQL语句
	ErrMsg string `json:"err_msg"`        // 事务消息
	Count  int64  `json:"count"`          // sql作用条数
}

// ModifyParamError 用于序列化查询接口异常信息
type ModifyParamError struct {
	Code   int    `json:"code"`    // 业务状态码
	ErrMsg string `json:"err_msg"` // 异常信息
}
