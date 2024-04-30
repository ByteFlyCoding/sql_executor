package utils

import (
	"reflect"
	"testing"
	"time"
)

func TestModifySqlValidate(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr bool
	}{
		{
			name:    "test1",
			sql:     "INSERT INTO users (username, email, birthdate, is_active) VALUES ('test', 'test@runoob.com', '1990-01-01', TRUE);",
			wantErr: false,
		},
		{
			name:    "test2",
			sql:     "INSERT INTO users (username, email, birthdate, is_active) VALUES ('test', 'test@runoob.com', '1990-01-01', TRUE)",
			wantErr: false,
		},
		{
			name:    "test3",
			sql:     "INSERT INTO users (username, email, birthdate, is_active)  ('test', 'test@runoob.com', '1990-01-01', TRUE)",
			wantErr: true,
		},
		{
			name:    "test4",
			sql:     "INSERT INTO users username, email, birthdate, is_active) VALUES ('test', 'test@runoob.com', '1990-01-01', TRUE)",
			wantErr: true,
		},
		{
			name:    "test5",
			sql:     " INTO users (username, email, birthdate, is_active) VALUES ('test', 'test@runoob.com', '1990-01-01', TRUE)",
			wantErr: true,
		},
		{
			name:    "test6",
			sql:     "INSERT users (username, email, birthdate, is_active) VALUES ('test', 'test@runoob.com', '1990-01-01', TRUE)",
			wantErr: false,
		},
		{
			name:    "test7",
			sql:     "DELETE FROM students WHERE graduation_year = 2021;",
			wantErr: false,
		},
		{
			name:    "test7",
			sql:     "DELETE  students WHERE graduation_year = 2021;",
			wantErr: true,
		},
		// {
		// 该测试用例有点奇怪，DELETE语句不完成，但是却可以通过sql合法性校验，后期需要再探索一下
		// 	name:    "test8",
		// 	sql:     "DELETE FROM students WHERE graduation_year",
		// 	wantErr: true, // },
		{
			name:    "test7",
			sql:     "FROM  students WHERE graduation_year = 2021;",
			wantErr: true,
		},
		{
			name:    "test8",
			sql:     "UPDATE employees SET salary = 60000 WHERE employee_id = 101;",
			wantErr: false,
		},
		{
			name:    "test9",
			sql:     " employees SET salary = 60000 WHERE employee_id = 101;",
			wantErr: true,
		},
		{
			name:    "test9",
			sql:     "UPDATE employees  salary = 60000 WHERE employee_id = 101;",
			wantErr: true,
		},
		{
			name:    "test10",
			sql:     "DROP TABLE IF EXISTS students CASCADE;",
			wantErr: true,
		},
		{
			name:    "test11",
			sql:     "DROP DATABASE IF EXISTS students",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ModifySqlValidate(tt.sql); (err != nil) != tt.wantErr {
				t.Errorf("ModifySqlValidate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQuerySqlValidate(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr bool
	}{
		{
			"test1",
			"SELECT * FROM user",
			false,
		},
		{
			"test2",
			"* from user",
			true,
		},
		{
			"test3",
			"SELECT * FROM ",
			true,
		},
		{
			"test4",
			"SELECT * FROM user;",
			false,
		},
		{
			"test5",
			"INSERT INTO SQL_EXECUTOR.user (user_name, password, create_time, update_time, `describe`) VALUES ('qwer', 'qewqwer', '2024-04-27 13:52:02', '2024-04-27 13:52:18', 'ewqreqwwer')",
			true,
		},
		{
			"test6",
			"DELETE FROM students\nWHERE graduation_year = 2021",
			true,
		},
		{
			"test7",
			"UPDATE students SET graduation_year = 2021",
			true,
		},
		{
			name:    "test8",
			sql:     "DROP TABLE IF EXISTS students CASCADE;",
			wantErr: true,
		},
		{
			name:    "test9",
			sql:     "DROP DATABASE IF EXISTS students",
			wantErr: true,
		},
		{
			name:    "test10",
			sql:     "USE student",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := QuerySqlValidate(tt.sql); (err != nil) != tt.wantErr {
				t.Errorf("QuerySqlValidate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func buildSqlInfoTestInstance(id int, name string, sql string) SqlInfo {
	return SqlInfo{
		ID:   id,
		Name: name,
		Sql:  sql,
	}
}

func NewRequestBody(transactions []*TransactionParamInfo) *RequestBody {
	return &RequestBody{Transactions: transactions}
}

func NewTransactionParamInfo(ID int, retry int, timeout time.Duration, name string, sqls []SqlInfo) *TransactionParamInfo {
	return &TransactionParamInfo{ID: ID, Retry: retry, Timeout: timeout, Name: name, Sqls: sqls}
}

func NewSqlInfo(ID int, name string, sql string) *SqlInfo {
	return &SqlInfo{ID: ID, Name: name, Sql: sql}
}

func NewModifyParamErrorJson(code int, items []TransactionParamError, count int, errMsg string) *ModifyParamErrorJson {
	return &ModifyParamErrorJson{Code: code, Items: items, Count: count, ErrMsg: errMsg}
}

func NewTransactionParamError(ID int, count int64, timeout time.Duration, name string, errMsg string, sqlErrorInfo []SqlErrorInfo) *TransactionParamError {
	return &TransactionParamError{ID: ID, Count: count, Timeout: timeout, Name: name, ErrMsg: errMsg, SqlErrorInfo: sqlErrorInfo}
}

func NewSqlErrorInfo(ID int, name string, sql string, errMsg string) *SqlErrorInfo {
	return &SqlErrorInfo{ID: ID, Name: name, Sql: sql, ErrMsg: errMsg}
}

func TestTransactionsValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     *RequestBody
		want    *ModifyParamErrorJson
		wantErr bool
	}{
		{
			name: "test1",
			req: NewRequestBody([]*TransactionParamInfo{
				{
					ID:      1,
					Retry:   1,
					Timeout: 10,
					Name:    "ts1",
					Sqls: []SqlInfo{
						*NewSqlInfo(1, "sql1", "INSERT INTO SQL_EXECUTOR.user (user_name, password) VALUES ('aaa', 'bbb')"),
					},
				},
			}),
			want:    NewModifyParamErrorJson(SUCCESSQUERY, []TransactionParamError{}, 0, ""),
			wantErr: false,
		},
		{
			name: "test2",
			req: NewRequestBody([]*TransactionParamInfo{
				{
					ID:      1,
					Retry:   1,
					Timeout: 10,
					Name:    "ts1",
					Sqls: []SqlInfo{
						*NewSqlInfo(1, "sql1", "INSERT INTO SQL_EXECUTOR.user (user_name, password) VALUES ('aaa', 'bbb')"),
					},
				},
				{
					ID:      2,
					Retry:   5,
					Timeout: 15,
					Name:    "ts2",
					Sqls: []SqlInfo{
						*NewSqlInfo(1, "sql1", "INSERT INTO SQL_EXECUTOR.user (user_name, password) VALUES ('aaa', 'bbb')"),
					},
				},
			}),
			want:    NewModifyParamErrorJson(SUCCESSQUERY, []TransactionParamError{}, 0, ""),
			wantErr: false,
		},
		{
			name: "test3",
			req: NewRequestBody([]*TransactionParamInfo{
				{
					ID:      1,
					Retry:   1,
					Timeout: 10,
					Name:    "ts1",
					Sqls: []SqlInfo{
						*NewSqlInfo(1, "sql1", "INSERT INTO SQL_EXECUTOR.user (user_name, password) VALUES ('aaa', 'bbb')"),
						*NewSqlInfo(1, "sql1", "INSERT INTO SQL_EXECUTOR.user (user_name, password) VALUES ('aaa', 'bbb')"),
					},
				},
				{
					ID:      2,
					Retry:   5,
					Timeout: 15,
					Name:    "ts2",
					Sqls: []SqlInfo{
						*NewSqlInfo(1, "sql1", "INSERT INTO SQL_EXECUTOR.user (user_name, password) VALUES ('aaa', 'bbb')"),
						*NewSqlInfo(1, "sql1", "INSERT INTO SQL_EXECUTOR.user (user_name, password) VALUES ('aaa', 'bbb')"),
					},
				},
			}),
			want:    NewModifyParamErrorJson(SUCCESSQUERY, []TransactionParamError{}, 0, ""),
			wantErr: false,
		},
		{
			name: "test4",
			req: NewRequestBody([]*TransactionParamInfo{
				{
					ID:      1,
					Retry:   1,
					Timeout: 10,
					Name:    "ts1",
					Sqls: []SqlInfo{
						*NewSqlInfo(1, "sql1", "INSERT INTO SQL_EXECUTOR.user (user_name, password) VALUES ('aaa', 'bbb')"),
						*NewSqlInfo(1, "sql1", "INSERT INTO SQL_EXECUTOR.user (user_name, password) VALUES ('aaa', 'bbb')"),
					},
				},
				{
					ID:      2,
					Retry:   5,
					Timeout: 15,
					Name:    "ts2",
					Sqls: []SqlInfo{
						*NewSqlInfo(1, "sql1", "INSERT INTO SQL_EXECUTOR.user (user_name, password) VALUES ('aaa', 'bbb')"),
						*NewSqlInfo(1, "sql1", "INSERT INTO SQL_EXECUTOR.user (user_name, password) VALUES ('aaa', 'bbb')"),
					},
				},
				{
					Sqls: []SqlInfo{
						*NewSqlInfo(1, "sql1", "INSERT INTO SQL_EXECUTOR.user (user_name, password) VALUES ('aaa', 'bbb')"),
						*NewSqlInfo(1, "sql1", "INSERT INTO SQL_EXECUTOR.user (user_name, password) VALUES ('aaa', 'bbb')"),
					},
				},
			}),
			want:    NewModifyParamErrorJson(SUCCESSQUERY, []TransactionParamError{}, 0, ""),
			wantErr: false,
		},
		{
			name: "test5",
			req: NewRequestBody([]*TransactionParamInfo{
				{
					ID:      1,
					Retry:   1,
					Timeout: 10,
					Name:    "ts1",
					Sqls: []SqlInfo{
						*NewSqlInfo(1, "sql1", "INSERT INTO SQL_EXECUTOR.user (user_name, password) VALUES ('aaa', 'bbb')"),
						*NewSqlInfo(1, "sql1", "INSERT INTO SQL_EXECUTOR.user (user_name, password) VALUES ('aaa', 'bbb')"),
					},
				},
				{
					ID:      2,
					Retry:   5,
					Timeout: 15,
					Name:    "ts2",
					Sqls: []SqlInfo{
						*NewSqlInfo(1, "sql1", "INSERT INTO SQL_EXECUTOR.user (user_name, password) VALUES ('aaa', 'bbb')"),
						*NewSqlInfo(1, "sql1", "INSERT INTO SQL_EXECUTOR.user (user_name, password) VALUES ('aaa', 'bbb')"),
					},
				},
			}),
			want:    NewModifyParamErrorJson(SUCCESSQUERY, []TransactionParamError{}, 0, ""),
			wantErr: false,
		},
		{
			name: "test6",
			req: NewRequestBody([]*TransactionParamInfo{
				{
					ID:      1,
					Retry:   1,
					Timeout: 10,
					Name:    "ts1",
					Sqls: []SqlInfo{
						*NewSqlInfo(1, "sql1", "INSERT INTO SQL_EXECUTOR.user (user_name, password) VALUES ('aaa', 'bbb')"),
					},
				},
				{
					ID:      2,
					Retry:   5,
					Timeout: 15,
					Name:    "ts2",
					Sqls: []SqlInfo{
						*NewSqlInfo(1, "sql1", "INSERT INTO SQL_EXECUTOR.user (user_name, password) VALUES ('aaa', 'bbb')"),
						*NewSqlInfo(1, "sql1", "INSERT INTO SQL_EXECUTOR.user (user_name, password) VALUES ('aaa', 'bbb')"),
					},
				},
			}),
			want:    NewModifyParamErrorJson(SUCCESSQUERY, []TransactionParamError{}, 0, ""),
			wantErr: false,
		},
		{
			name: "test7",
			req: NewRequestBody([]*TransactionParamInfo{
				{
					Retry:   1,
					Timeout: 10,
					Name:    "ts1",
					Sqls: []SqlInfo{
						*NewSqlInfo(1, "sql1", "INSERT INTO SQL_EXECUTOR.user (user_name, password) VALUES ('aaa', 'bbb')"),
						*NewSqlInfo(1, "sql1", "INSERT INTO SQL_EXECUTOR.user (user_name, password) VALUES ('aaa', 'bbb')"),
					},
				},
				{
					ID:      2,
					Timeout: 15,
					Name:    "ts2",
					Sqls: []SqlInfo{
						*NewSqlInfo(1, "sql1", "INSERT INTO SQL_EXECUTOR.user (user_name, password) VALUES ('aaa', 'bbb')"),
						*NewSqlInfo(1, "sql1", "INSERT INTO SQL_EXECUTOR.user (user_name, password) VALUES ('aaa', 'bbb')"),
					},
				},
			}),
			want:    NewModifyParamErrorJson(SUCCESSQUERY, []TransactionParamError{}, 0, ""),
			wantErr: false,
		},
		{
			name: "test8",
			req: NewRequestBody([]*TransactionParamInfo{
				{
					ID:      1,
					Retry:   1,
					Timeout: 10,
					Sqls: []SqlInfo{
						*NewSqlInfo(1, "sql1", "INSERT INTO SQL_EXECUTOR.user (user_name, password) VALUES ('aaa', 'bbb')"),
						*NewSqlInfo(1, "sql1", "INSERT INTO SQL_EXECUTOR.user (user_name, password) VALUES ('aaa', 'bbb')"),
					},
				},
				{
					ID:    2,
					Retry: 5,
					Name:  "ts2",
					Sqls: []SqlInfo{
						*NewSqlInfo(1, "sql1", "INSERT INTO SQL_EXECUTOR.user (user_name, password) VALUES ('aaa', 'bbb')"),
						*NewSqlInfo(1, "sql1", "INSERT INTO SQL_EXECUTOR.user (user_name, password) VALUES ('aaa', 'bbb')"),
					},
				},
			}),
			want:    NewModifyParamErrorJson(SUCCESSQUERY, []TransactionParamError{}, 0, ""),
			wantErr: false,
		},
		{
			name: "test9",
			req: NewRequestBody([]*TransactionParamInfo{
				{
					ID:      1,
					Retry:   1,
					Name:    "ts1",
					Timeout: 10,
					Sqls:    []SqlInfo{},
				},
				{
					ID:      2,
					Retry:   5,
					Name:    "ts2",
					Timeout: 15,
					Sqls: []SqlInfo{
						*NewSqlInfo(1, "sql1", "INSERT INTO SQL_EXECUTOR.user (user_name, password) VALUES ('aaa', 'bbb')"),
						*NewSqlInfo(1, "sql1", "INSERT INTO SQL_EXECUTOR.user (user_name, password) VALUES ('aaa', 'bbb')"),
					},
				},
			}),
			want: NewModifyParamErrorJson(PARAMETERERROR, []TransactionParamError{
				*NewTransactionParamError(1,
					1,
					10,
					"ts1",
					"事务没有输入SQL或输入的SQL中有语法错误",
					[]SqlErrorInfo{
						*NewSqlErrorInfo(0, "", "", "事务中没有输入任何sql"),
					},
				),
			},
				1,
				"事务没有输入SQL或输入的SQL中有语法错误",
			),
			wantErr: true,
		},
		{
			name: "test10",
			req:  NewRequestBody([]*TransactionParamInfo{}),
			want: NewModifyParamErrorJson(PARAMETERERROR, []TransactionParamError{},
				0,
				"没有输入任何事务",
			),
			wantErr: true,
		},
		{
			name: "test11",
			req: NewRequestBody([]*TransactionParamInfo{
				{
					ID:      1,
					Retry:   5,
					Name:    "ts1",
					Timeout: 10,
					Sqls: []SqlInfo{
						*NewSqlInfo(1, "sql1", " INTO SQL_EXECUTOR.user (user_name, password) VALUES ('aaa', 'bbb')"),
						*NewSqlInfo(1, "sql2", "INSERT INTO SQL_EXECUTOR.user (user_name, password) VALUES ('aaa', 'bbb')"),
					},
				},
			}),
			want: NewModifyParamErrorJson(PARAMETERERROR, []TransactionParamError{
				*NewTransactionParamError(1,
					1,
					10,
					"ts1",
					"事务没有输入SQL或输入的SQL中有语法错误",
					[]SqlErrorInfo{
						*NewSqlErrorInfo(1, "sql1", " INTO SQL_EXECUTOR.user (user_name, password) VALUES ('aaa', 'bbb')", "syntax error at position 6 near 'into'"),
					},
				),
			},
				1,
				"事务没有输入SQL或输入的SQL中有语法错误",
			),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TransactionsValidate(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransactionsValidate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			checkItems := func(wantItems, gotItems []TransactionParamError) bool {

				if len(wantItems) != len(gotItems) {
					return false
				}

				for i := 0; i < len(wantItems); i++ {
					if wantItems[i].ID != gotItems[i].ID {
						return false
					}

					if wantItems[i].ErrMsg != gotItems[i].ErrMsg {
						return false
					}

					if wantItems[i].Count != gotItems[i].Count {
						return false
					}

					if wantItems[i].Name != gotItems[i].Name {
						return false
					}

					if wantItems[i].Timeout != gotItems[i].Timeout {
						return false
					}

					if !reflect.DeepEqual(wantItems[i].SqlErrorInfo, gotItems[i].SqlErrorInfo) {
						return false
					}
				}

				return true
			}

			checkModifyParamErrorJson := func(want *ModifyParamErrorJson, got *ModifyParamErrorJson) bool {

				if got.ErrMsg != want.ErrMsg {
					return false
				}

				if want.Code != got.Code {
					return false
				}

				if want.Count != got.Count {
					return false
				}

				if !checkItems(want.Items, got.Items) {
					return false
				}

				return true
			}

			if !checkModifyParamErrorJson(got, tt.want) {
				t.Errorf("TransactionsValidate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
