package utils

type SqlInfo struct {
	ID     int    `json:"id"`
	Name   string `json:"name,omitempty"`
	Sql    string `json:"sql"`
	ErrMsg string `json:"err_msg,omitempty"`
	Count  int64  `json:"count,omitempty"` // sql作用条数
}

type TransactionInfo struct {
	ID    int       `json:"id"`
	Retry int       `json:"retry,omitempty"` // 允许重试次数
	Name  string    `json:"name"`
	Sqls  []SqlInfo `json:"sqls"`
}

type RequestBody struct {
	Transactions []*TransactionInfo `json:"transactions"`
}

type ModifySqlErrorJson struct {
	Code   int    `json:"code"`
	ErrMsg string `json:"err_msg"`
	Count  int    `json:"count"`
	// Items  []SyntaxError `json:"items"`
	Items []Runner `json:"items"`
}

type Runner struct {
	ID      int       `json:"id"`
	Retry   int       `json:"retry,omitempty"`
	Count   int64     `json:"count"`
	Name    string    `json:"name,omitempty"`
	ErrMsg  string    `json:"err_msg"`
	SqlInfo []SqlInfo `json:"Sql_info"`
}

type ModifySuccessJson struct {
	Code  int      `json:"code"`
	Items []Runner `json:"items"`
	Count int      `json:"count"`
	Msg   string   `json:"msg"`
}
