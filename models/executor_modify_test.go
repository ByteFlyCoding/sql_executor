package models

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/beego/beego/v2/client/orm"
	utils2 "github.com/beego/beego/v2/core/utils"
	"github.com/pkg/errors"
	"sql_executor/utils"
)

func TestExecutor_Modify(t *testing.T) {

	tests := []struct {
		name       string
		endpoint   *mockModifyEndpoint
		wantRunner *utils.Runner
		wantErr    bool
	}{
		{
			// 一个事务，里面只有一个SQL，commit成功，无需重试
			name: "test1",
			endpoint: &mockModifyEndpoint{
				factRetry:       0,
				id:              1,
				isOpenTxSuccess: true,
				isSuccess:       true,
				retry:           0,
				timeout:         10,
				name:            "test 1",
				count:           1,
				execHistory: []sqlsExecHistory{
					{
						id:            1,
						name:          "test 1-1",
						sql:           "DELETE FROM students WHERE graduation_year = 1900;",
						isExecErr:     false,
						isRollback:    false,
						effectRow:     1,
						resultIsNil:   false,
						isRowEffect:   true,
						isNoRowsError: true,
					},
				},
				sqlInfo: []utils.SqlInfo{
					{
						ID:   1,
						Name: "test 1-1",
						Sql:  "DELETE FROM students WHERE graduation_year = 1900;",
					},
				},
			},
			wantRunner: &utils.Runner{
				ID:      1,
				Retry:   0,
				Timeout: 10,
				Name:    "test 1",
				ErrMsg:  "事务提交成功",
				SqlExecInfo: []utils.SqlExecInfo{
					{
						ID:     1,
						Name:   "test 1-1",
						Sql:    "DELETE FROM students WHERE graduation_year = 1900;",
						ErrMsg: "该SQL执行成功，等待事务提交",
						Count:  1,
					},
				},
			},
			wantErr: false,
		},
		{
			// 一个事务，里面有多条SQL，commit成功，无需重试
			name: "test2",
			endpoint: &mockModifyEndpoint{
				factRetry:       0,
				id:              1,
				isOpenTxSuccess: true,
				isSuccess:       true,
				retry:           0,
				timeout:         10,
				name:            "test 2",
				count:           3,
				execHistory: []sqlsExecHistory{
					{
						id:            1,
						name:          "test 2-1",
						sql:           "DELETE FROM students WHERE graduation_year = 1900;",
						isExecErr:     false,
						isRollback:    false,
						effectRow:     1,
						resultIsNil:   false,
						isRowEffect:   true,
						isNoRowsError: true,
					},
					{
						id:            2,
						name:          "test 2-2",
						sql:           "DELETE FROM students WHERE graduation_year = 1901;",
						isExecErr:     false,
						isRollback:    false,
						effectRow:     2,
						resultIsNil:   false,
						isRowEffect:   true,
						isNoRowsError: true,
					},
					{
						id:            3,
						name:          "test 2-3",
						sql:           "DELETE FROM students WHERE graduation_year = 1902;",
						isExecErr:     false,
						isRollback:    false,
						effectRow:     3,
						resultIsNil:   false,
						isRowEffect:   true,
						isNoRowsError: true,
					},
				},
				sqlInfo: []utils.SqlInfo{
					{
						ID:   1,
						Name: "test 2-1",
						Sql:  "DELETE FROM students WHERE graduation_year = 1900;",
					},
					{
						ID:   2,
						Name: "test 2-2",
						Sql:  "DELETE FROM students WHERE graduation_year = 1901;",
					},
					{
						ID:   3,
						Name: "test 2-3",
						Sql:  "DELETE FROM students WHERE graduation_year = 1902;",
					},
				},
			},
			wantRunner: &utils.Runner{
				ID:      1,
				Retry:   0,
				Timeout: 10,
				Name:    "test 2",
				ErrMsg:  "事务提交成功",
				SqlExecInfo: []utils.SqlExecInfo{
					{
						ID:     1,
						Name:   "test 2-1",
						Sql:    "DELETE FROM students WHERE graduation_year = 1900;",
						ErrMsg: "该SQL执行成功，等待事务提交",
						Count:  1,
					},
					{
						ID:     2,
						Name:   "test 2-2",
						Sql:    "DELETE FROM students WHERE graduation_year = 1901;",
						ErrMsg: "该SQL执行成功，等待事务提交",
						Count:  2,
					},
					{
						ID:     3,
						Name:   "test 2-3",
						Sql:    "DELETE FROM students WHERE graduation_year = 1902;",
						ErrMsg: "该SQL执行成功，等待事务提交",
						Count:  3,
					},
				},
			},
			wantErr: false,
		},
		{
			// 一个事务，里面只有一条SQL，开启事务失败，最大重试次数为0
			name: "test3",
			endpoint: &mockModifyEndpoint{
				factRetry:       0,
				id:              1,
				isOpenTxSuccess: false,
				isSuccess:       false,
				retry:           0,
				timeout:         10,
				name:            "test 3",
				count:           1,
				execHistory:     []sqlsExecHistory{},
			},
			wantRunner: &utils.Runner{
				ID:          1,
				Retry:       0,
				Timeout:     10,
				Name:        "test 3",
				ErrMsg:      "开启事务失败，",
				SqlExecInfo: []utils.SqlExecInfo{},
			},
			wantErr: true,
		},
		{
			// 一个事务，里面只有一条SQL，开启事务失败，重试一次后成功提交
			name: "test4",
			endpoint: &mockModifyEndpoint{
				factRetry:       1,
				id:              1,
				isOpenTxSuccess: true,
				isSuccess:       true,
				retry:           1,
				timeout:         10,
				name:            "test 4",
				count:           1,
				execHistory: []sqlsExecHistory{
					{
						id:   1,
						name: "test 4-1",
						sql:  "DELETE FROM students WHERE graduation_year = 1900;",
						// errMsg:      "该SQL执行成功，等待事务提交",
						isExecErr:     false,
						isRollback:    false,
						effectRow:     1,
						resultIsNil:   false,
						isRowEffect:   true,
						isNoRowsError: true,
					},
				},
				sqlInfo: []utils.SqlInfo{
					{
						ID:   1,
						Name: "test 4-1",
						Sql:  "DELETE FROM students WHERE graduation_year = 1900;",
					},
				},
			},
			wantRunner: &utils.Runner{
				ID:      1,
				Retry:   1,
				Timeout: 10,
				Name:    "test 4",
				ErrMsg:  "事务提交成功",
				SqlExecInfo: []utils.SqlExecInfo{
					{
						ID:     1,
						Name:   "test 4-1",
						Sql:    "DELETE FROM students WHERE graduation_year = 1900;",
						ErrMsg: "该SQL执行成功，等待事务提交",
						Count:  1,
					},
				},
			},
			wantErr: false,
		},
		{
			// 一个事务，里面只有一条SQL，开启事务失败，最大重试次数为1，重试1次后依旧失败
			name: "test5",
			endpoint: &mockModifyEndpoint{
				factRetry:       1,
				id:              1,
				isOpenTxSuccess: false,
				isSuccess:       true,
				retry:           1,
				timeout:         10,
				name:            "test 5",
				count:           1,
				execHistory: []sqlsExecHistory{
					{
						id:          1,
						name:        "test 5-1",
						sql:         "DELETE FROM students WHERE graduation_year = 1900;",
						isExecErr:   true,
						isRollback:  false,
						effectRow:   0,
						resultIsNil: false,
						isRowEffect: true,
					},
				},
				sqlInfo: []utils.SqlInfo{
					{
						ID:   1,
						Name: "test 5-1",
						Sql:  "DELETE FROM students WHERE graduation_year = 1900;",
					},
				},
			},
			wantRunner: &utils.Runner{
				ID:          1,
				Retry:       1,
				Timeout:     10,
				Name:        "test 5",
				ErrMsg:      "开启事务失败",
				SqlExecInfo: []utils.SqlExecInfo{},
			},
			wantErr: true,
		},
		{
			// 一个事务，里面有多条SQL, 开启事务失败，最大重试次数为1，重试后成功
			name: "test6",
			endpoint: &mockModifyEndpoint{
				factRetry:       1,
				id:              1,
				isOpenTxSuccess: true,
				isSuccess:       true,
				retry:           1,
				timeout:         10,
				name:            "test 6",
				count:           3,
				execHistory: []sqlsExecHistory{
					{
						id:            1,
						name:          "test 6-1",
						sql:           "DELETE FROM students WHERE graduation_year = 1900;",
						isExecErr:     false,
						isRollback:    false,
						effectRow:     1,
						resultIsNil:   false,
						isRowEffect:   true,
						isNoRowsError: true,
					},
					{
						id:            2,
						name:          "test 6-2",
						sql:           "DELETE FROM students WHERE graduation_year = 1900;",
						isExecErr:     false,
						isRollback:    false,
						effectRow:     2,
						resultIsNil:   false,
						isRowEffect:   true,
						isNoRowsError: true,
					},
					{
						id:            3,
						name:          "test 6-3",
						sql:           "DELETE FROM students WHERE graduation_year = 1900;",
						isExecErr:     false,
						isRollback:    false,
						effectRow:     3,
						resultIsNil:   false,
						isRowEffect:   true,
						isNoRowsError: true,
					},
				},
				sqlInfo: []utils.SqlInfo{
					{
						ID:   1,
						Name: "test 6-1",
						Sql:  "DELETE FROM students WHERE graduation_year = 1900;",
					},
					{
						ID:   2,
						Name: "test 6-2",
						Sql:  "DELETE FROM students WHERE graduation_year = 1900;",
					},
					{
						ID:   3,
						Name: "test 6-3",
						Sql:  "DELETE FROM students WHERE graduation_year = 1900;",
					},
				},
			},
			wantRunner: &utils.Runner{
				ID:      1,
				Retry:   1,
				Timeout: 10,
				Name:    "test 6",
				ErrMsg:  "事务提交成功",
				SqlExecInfo: []utils.SqlExecInfo{
					{
						ID:     1,
						Name:   "test 6-1",
						Sql:    "DELETE FROM students WHERE graduation_year = 1900;",
						ErrMsg: "该SQL执行成功，等待事务提交",
						Count:  1,
					},
					{
						ID:     2,
						Name:   "test 6-2",
						Sql:    "DELETE FROM students WHERE graduation_year = 1900;",
						ErrMsg: "该SQL执行成功，等待事务提交",
						Count:  2,
					},
					{
						ID:     3,
						Name:   "test 6-3",
						Sql:    "DELETE FROM students WHERE graduation_year = 1900;",
						ErrMsg: "该SQL执行成功，等待事务提交",
						Count:  3,
					},
				},
			},
			wantErr: false,
		},
		{
			// 一个事务，里面有多条SQL, 开启事务失败，最大重试次数为2，重试2次后成功
			name: "test7",
			endpoint: &mockModifyEndpoint{
				factRetry:       2,
				id:              1,
				isOpenTxSuccess: true,
				isSuccess:       true,
				retry:           2,
				timeout:         10,
				name:            "test 7",
				count:           3,
				execHistory: []sqlsExecHistory{
					{
						id:            1,
						name:          "test 7-1",
						sql:           "DELETE FROM students WHERE graduation_year = 1900;",
						isExecErr:     false,
						isRollback:    false,
						effectRow:     1,
						resultIsNil:   false,
						isRowEffect:   true,
						isNoRowsError: true,
					},
					{
						id:            2,
						name:          "test 7-2",
						sql:           "DELETE FROM students WHERE graduation_year = 1900;",
						isExecErr:     false,
						isRollback:    false,
						effectRow:     2,
						resultIsNil:   false,
						isRowEffect:   true,
						isNoRowsError: true,
					},
					{
						id:            3,
						name:          "test 7-3",
						sql:           "DELETE FROM students WHERE graduation_year = 1900;",
						isExecErr:     false,
						isRollback:    false,
						effectRow:     3,
						resultIsNil:   false,
						isRowEffect:   true,
						isNoRowsError: true,
					},
				},
				sqlInfo: []utils.SqlInfo{
					{
						ID:   1,
						Name: "test 7-1",
						Sql:  "DELETE FROM students WHERE graduation_year = 1900;",
					},
					{
						ID:   2,
						Name: "test 7-2",
						Sql:  "DELETE FROM students WHERE graduation_year = 1900;",
					},
					{
						ID:   3,
						Name: "test 7-3",
						Sql:  "DELETE FROM students WHERE graduation_year = 1900;",
					},
				},
			},
			wantRunner: &utils.Runner{
				ID:      1,
				Retry:   2,
				Timeout: 10,
				Name:    "test 7",
				ErrMsg:  "事务提交成功",
				SqlExecInfo: []utils.SqlExecInfo{
					{
						ID:     1,
						Name:   "test 7-1",
						Sql:    "DELETE FROM students WHERE graduation_year = 1900;",
						ErrMsg: "该SQL执行成功，等待事务提交",
						Count:  1,
					},
					{
						ID:     2,
						Name:   "test 7-2",
						Sql:    "DELETE FROM students WHERE graduation_year = 1900;",
						ErrMsg: "该SQL执行成功，等待事务提交",
						Count:  2,
					},
					{
						ID:     3,
						Name:   "test 7-3",
						Sql:    "DELETE FROM students WHERE graduation_year = 1900;",
						ErrMsg: "该SQL执行成功，等待事务提交",
						Count:  3,
					},
				},
			},
			wantErr: false,
		},
		{
			// 一个事务，里面有多条SQL, 开启事务失败，最大重试次数为2，重试2次依旧失败
			name: "test8",
			endpoint: &mockModifyEndpoint{
				factRetry:       2,
				id:              1,
				isOpenTxSuccess: false,
				isSuccess:       false,
				retry:           2,
				timeout:         10,
				name:            "test 8",
				count:           3,
				execHistory: []sqlsExecHistory{
					{
						id:          1,
						name:        "test 8-1",
						sql:         "DELETE FROM students WHERE graduation_year = 1900;",
						isExecErr:   false,
						isRollback:  false,
						effectRow:   1,
						resultIsNil: false,
						isRowEffect: true,
					},
					{
						id:          2,
						name:        "test 8-2",
						sql:         "DELETE FROM students WHERE graduation_year = 1900;",
						isExecErr:   false,
						isRollback:  false,
						effectRow:   2,
						resultIsNil: false,
						isRowEffect: true,
					},
					{
						id:          3,
						name:        "test 8-3",
						sql:         "DELETE FROM students WHERE graduation_year = 1900;",
						isExecErr:   false,
						isRollback:  false,
						effectRow:   3,
						resultIsNil: false,
						isRowEffect: true,
					},
				},
				sqlInfo: []utils.SqlInfo{
					{
						ID:   1,
						Name: "test 8-1",
						Sql:  "DELETE FROM students WHERE graduation_year = 1900;",
					},
					{
						ID:   2,
						Name: "test 8-2",
						Sql:  "DELETE FROM students WHERE graduation_year = 1900;",
					},
					{
						ID:   3,
						Name: "test 8-3",
						Sql:  "DELETE FROM students WHERE graduation_year = 1900;",
					},
				},
			},
			wantRunner: &utils.Runner{
				ID:          1,
				Retry:       2,
				Timeout:     10,
				Name:        "test 8",
				ErrMsg:      "开启事务失败",
				SqlExecInfo: []utils.SqlExecInfo{},
			},
			wantErr: true,
		},
		{
			// 一个事务，里面只有一条SQL，开启事务成功，执行SQL失败，回滚失败
			name: "test9",
			endpoint: &mockModifyEndpoint{
				factRetry:       0,
				id:              1,
				isOpenTxSuccess: true,
				isSuccess:       false,
				retry:           0,
				timeout:         10,
				name:            "test 9",
				count:           1,
				execHistory: []sqlsExecHistory{
					{
						id:   1,
						name: "test 9-1",
						sql:  "DELETE FROM students WHERE graduation_year = 1900;",
						// errMsg:      "事务执行失败，已回滚：",
						isExecErr:   true,
						isRollback:  true,
						effectRow:   1,
						resultIsNil: false,
						isRowEffect: true,
					},
				},
				sqlInfo: []utils.SqlInfo{
					{
						ID:   1,
						Name: "test 9-1",
						Sql:  "DELETE FROM students WHERE graduation_year = 1900;",
					},
				},
			},
			wantRunner: &utils.Runner{
				ID:      1,
				Retry:   0,
				Timeout: 10,
				Name:    "test 9",
				ErrMsg:  "事务执行失败，已回滚：",
				SqlExecInfo: []utils.SqlExecInfo{
					{
						ID:     1,
						Name:   "test 9-1",
						Sql:    "DELETE FROM students WHERE graduation_year = 1900;",
						ErrMsg: "SQL执行失败，等待回滚mock：SQL执行失败，事务回滚成功",
						Count:  0,
					},
				},
			},
			wantErr: true,
		},
		{
			// 一个事务，里面只有一条SQL，开启事务成功，执行SQL失败，回滚失败
			name: "test10",
			endpoint: &mockModifyEndpoint{
				factRetry:       0,
				id:              1,
				isOpenTxSuccess: true,
				isSuccess:       false,
				retry:           0,
				timeout:         10,
				name:            "test 10",
				count:           1,
				execHistory: []sqlsExecHistory{
					{
						id:          1,
						name:        "test 10-1",
						sql:         "DELETE FROM students WHERE graduation_year = 1900;",
						isExecErr:   true,
						isRollback:  false,
						effectRow:   1,
						resultIsNil: false,
						isRowEffect: true,
					},
				},
				sqlInfo: []utils.SqlInfo{
					{
						ID:   1,
						Name: "test 10-1",
						Sql:  "DELETE FROM students WHERE graduation_year = 1900;",
					},
				},
			},
			wantRunner: &utils.Runner{
				ID:      1,
				Retry:   0,
				Timeout: 10,
				Name:    "test 10",
				ErrMsg:  "事务回滚失败，等待自动回滚：",
				SqlExecInfo: []utils.SqlExecInfo{
					{
						ID:     1,
						Name:   "test 10-1",
						Sql:    "DELETE FROM students WHERE graduation_year = 1900;",
						ErrMsg: "SQL执行失败，等待回滚mock：SQL执行失败，事务回滚失败，等待自动回滚",
						Count:  0,
					},
				},
			},
			wantErr: true,
		},
		{
			// 一个事务，里面只有一条SQL，开启事务成功，执行SQL成功，result为nil，提交成功
			name: "test11",
			endpoint: &mockModifyEndpoint{
				factRetry:       0,
				id:              1,
				isOpenTxSuccess: true,
				isSuccess:       true,
				retry:           0,
				timeout:         10,
				name:            "test 11",
				count:           1,
				execHistory: []sqlsExecHistory{
					{
						id:            1,
						name:          "test 11-1",
						sql:           "DELETE FROM students WHERE graduation_year = 1900;",
						isExecErr:     false,
						isRollback:    false,
						effectRow:     1,
						resultIsNil:   false,
						isRowEffect:   true,
						isNoRowsError: true,
					},
				},
				sqlInfo: []utils.SqlInfo{
					{
						ID:   1,
						Name: "test 11-1",
						Sql:  "DELETE FROM students WHERE graduation_year = 1900;",
					},
				},
			},
			wantRunner: &utils.Runner{
				ID:      1,
				Retry:   0,
				Timeout: 10,
				Name:    "test 11",
				ErrMsg:  "事务提交成功",
				SqlExecInfo: []utils.SqlExecInfo{
					{
						ID:     1,
						Name:   "test 11-1",
						Sql:    "DELETE FROM students WHERE graduation_year = 1900;",
						ErrMsg: "该SQL执行成功，等待事务提交",
						Count:  1,
					},
				},
			},
			wantErr: false,
		},
		{
			// 一个事务，里面只有一条SQL，开启事务成功，执行SQL成功，result为nil，提交失败
			name: "test12",
			endpoint: &mockModifyEndpoint{
				factRetry:       0,
				id:              1,
				isOpenTxSuccess: true,
				isSuccess:       false,
				retry:           0,
				timeout:         10,
				name:            "test 12",
				count:           1,
				execHistory: []sqlsExecHistory{
					{
						id:          1,
						name:        "test 12-1",
						sql:         "DELETE FROM students WHERE graduation_year = 1900;",
						isExecErr:   false,
						isRollback:  false,
						effectRow:   0,
						resultIsNil: true,
					},
				},
				sqlInfo: []utils.SqlInfo{
					{
						ID:   1,
						Name: "test 12-1",
						Sql:  "DELETE FROM students WHERE graduation_year = 1900;",
					},
				},
			},
			wantRunner: &utils.Runner{
				ID:          1,
				Retry:       0,
				Timeout:     10,
				Name:        "test 12",
				ErrMsg:      "该事务提交失败，已回滚：",
				SqlExecInfo: []utils.SqlExecInfo{},
			},
			wantErr: true,
		},
		{
			// 一个事务，里面只有一条SQL，开启事务成功，执行SQL成功，result不为nil，提交成功
			name: "test13",
			endpoint: &mockModifyEndpoint{
				factRetry:       0,
				id:              1,
				isOpenTxSuccess: true,
				isSuccess:       true,
				retry:           0,
				timeout:         10,
				name:            "test 13",
				count:           1,
				execHistory: []sqlsExecHistory{
					{
						id:            1,
						name:          "test 13-1",
						sql:           "DELETE FROM students WHERE graduation_year = 1900;",
						isExecErr:     false,
						isRollback:    false,
						effectRow:     1,
						resultIsNil:   false,
						isRowEffect:   true,
						isNoRowsError: true,
					},
				},
				sqlInfo: []utils.SqlInfo{
					{
						ID:   1,
						Name: "test 13-1",
						Sql:  "DELETE FROM students WHERE graduation_year = 1900;",
					},
				},
			},
			wantRunner: &utils.Runner{
				ID:      1,
				Retry:   0,
				Timeout: 10,
				Name:    "test 13",
				ErrMsg:  "事务提交成功",
				SqlExecInfo: []utils.SqlExecInfo{
					{
						ID:     1,
						Name:   "test 13-1",
						Sql:    "DELETE FROM students WHERE graduation_year = 1900;",
						ErrMsg: "该SQL执行成功，等待事务提交",
						Count:  1,
					},
				},
			},
			wantErr: false,
		},
		{
			// 一个事务，里面只有一条SQL，开启事务成功，执行SQL成功，result不为nil，没有出现空行错误，提交成功
			name: "test14",
			endpoint: &mockModifyEndpoint{
				factRetry:       0,
				id:              1,
				isOpenTxSuccess: true,
				isSuccess:       true,
				retry:           0,
				timeout:         10,
				name:            "test 14",
				count:           1,
				execHistory: []sqlsExecHistory{
					{
						id:            1,
						name:          "test 14-1",
						sql:           "DELETE FROM students WHERE graduation_year = 1900;",
						isExecErr:     false,
						isRollback:    false,
						effectRow:     10,
						resultIsNil:   false,
						isRowEffect:   true,
						isNoRowsError: true,
					},
				},
				sqlInfo: []utils.SqlInfo{
					{
						ID:   1,
						Name: "test 14-1",
						Sql:  "DELETE FROM students WHERE graduation_year = 1900;",
					},
				},
			},
			wantRunner: &utils.Runner{
				ID:      1,
				Retry:   0,
				Timeout: 10,
				Name:    "test 14",
				ErrMsg:  "事务提交成功",
				SqlExecInfo: []utils.SqlExecInfo{
					{
						ID:     1,
						Name:   "test 14-1",
						Sql:    "DELETE FROM students WHERE graduation_year = 1900;",
						ErrMsg: "该SQL执行成功，等待事务提交",
						Count:  10,
					},
				},
			},
			wantErr: false,
		},
		{
			// 一个事务，里面只有一条SQL，开启事务成功，执行SQL成功，result不为nil，出现空行错误，提交成功
			name: "test14",
			endpoint: &mockModifyEndpoint{
				factRetry:       0,
				id:              1,
				isOpenTxSuccess: true,
				isSuccess:       true,
				retry:           0,
				timeout:         10,
				name:            "test 15",
				count:           1,
				execHistory: []sqlsExecHistory{
					{
						id:            1,
						name:          "test 15-1",
						sql:           "DELETE FROM students WHERE graduation_year = 1900;",
						isExecErr:     false,
						isRollback:    false,
						effectRow:     0,
						resultIsNil:   false,
						isRowEffect:   false,
						isNoRowsError: false,
					},
				},
				sqlInfo: []utils.SqlInfo{
					{
						ID:   1,
						Name: "test 15-1",
						Sql:  "DELETE FROM students WHERE graduation_year = 1900;",
					},
				},
			},
			wantRunner: &utils.Runner{
				ID:      1,
				Retry:   0,
				Timeout: 10,
				Name:    "test 15",
				ErrMsg:  "事务提交成功",
				SqlExecInfo: []utils.SqlExecInfo{
					{
						ID:     1,
						Name:   "test 15-1",
						Sql:    "DELETE FROM students WHERE graduation_year = 1900;",
						ErrMsg: "出现空行错误",
						Count:  0,
					},
				},
			},
			wantErr: false,
		},
		{
			// 一个事务，里面有多条SQL，开启事务成功，执行SQL成功，result不为nil，出现空行错误，提交成功
			name: "test16",
			endpoint: &mockModifyEndpoint{
				factRetry:       0,
				id:              1,
				isOpenTxSuccess: true,
				isSuccess:       true,
				retry:           0,
				timeout:         10,
				name:            "test 16",
				count:           1,
				execHistory: []sqlsExecHistory{
					{
						id:            1,
						name:          "test 16-1",
						sql:           "DELETE FROM students WHERE graduation_year = 1900;",
						isExecErr:     false,
						isRollback:    false,
						effectRow:     0,
						resultIsNil:   false,
						isRowEffect:   true,
						isNoRowsError: false,
					},
					{
						id:            1,
						name:          "test 16-2",
						sql:           "DELETE FROM students WHERE graduation_year = 1900;",
						isExecErr:     false,
						isRollback:    false,
						effectRow:     0,
						resultIsNil:   false,
						isRowEffect:   true,
						isNoRowsError: false,
					},
					{
						id:            1,
						name:          "test 16-3",
						sql:           "DELETE FROM students WHERE graduation_year = 1900;",
						isExecErr:     false,
						isRollback:    false,
						effectRow:     0,
						resultIsNil:   false,
						isRowEffect:   true,
						isNoRowsError: false,
					},
				},
				sqlInfo: []utils.SqlInfo{
					{
						ID:   1,
						Name: "test 16-1",
						Sql:  "DELETE FROM students WHERE graduation_year = 1900;",
					},
					{
						ID:   2,
						Name: "test 16-2",
						Sql:  "DELETE FROM students WHERE graduation_year = 1900;",
					},
					{
						ID:   3,
						Name: "test 16-3",
						Sql:  "DELETE FROM students WHERE graduation_year = 1900;",
					},
				},
			},
			wantRunner: &utils.Runner{
				ID:      1,
				Retry:   0,
				Timeout: 10,
				Name:    "test 16",
				ErrMsg:  "事务提交成功",
				SqlExecInfo: []utils.SqlExecInfo{
					{
						ID:     1,
						Name:   "test 16-1",
						Sql:    "DELETE FROM students WHERE graduation_year = 1900;",
						ErrMsg: "出现空行错误",
						Count:  0,
					}, {
						ID:     2,
						Name:   "test 16-2",
						Sql:    "DELETE FROM students WHERE graduation_year = 1900;",
						ErrMsg: "出现空行错误",
						Count:  0,
					}, {
						ID:     3,
						Name:   "test 16-3",
						Sql:    "DELETE FROM students WHERE graduation_year = 1900;",
						ErrMsg: "出现空行错误",
						Count:  0,
					},
				},
			},
			wantErr: false,
		},
		{
			// 一个事务，里面有多条SQL，开启事务成功，执行第二条SQL失败 回滚成功
			name: "test17",
			endpoint: &mockModifyEndpoint{
				factRetry:       0,
				id:              1,
				isOpenTxSuccess: true,
				isSuccess:       false,
				retry:           0,
				timeout:         10,
				name:            "test 16",
				count:           1,
				execHistory: []sqlsExecHistory{
					{
						id:            1,
						name:          "test 16-1",
						sql:           "DELETE FROM students WHERE graduation_year = 1900;",
						isExecErr:     false,
						isRollback:    false,
						effectRow:     10,
						resultIsNil:   false,
						isRowEffect:   true,
						isNoRowsError: true,
					},
					{
						id:            1,
						name:          "test 16-2",
						sql:           "DELETE FROM students WHERE graduation_year = 1900;",
						isExecErr:     true,
						isRollback:    true,
						effectRow:     0,
						resultIsNil:   false,
						isRowEffect:   true,
						isNoRowsError: false,
					},
				},
				sqlInfo: []utils.SqlInfo{
					{
						ID:   1,
						Name: "test 16-1",
						Sql:  "DELETE FROM students WHERE graduation_year = 1900;",
					},
					{
						ID:   2,
						Name: "test 16-2",
						Sql:  "DELETE FROM students WHERE graduation_year = 1900;",
					},
					{
						ID:   3,
						Name: "test 16-3",
						Sql:  "DELETE FROM students WHERE graduation_year = 1900;",
					},
				},
			},
			wantRunner: &utils.Runner{
				ID:      1,
				Retry:   0,
				Timeout: 10,
				Name:    "test 16",
				ErrMsg:  "事务执行失败，已回滚",
				SqlExecInfo: []utils.SqlExecInfo{
					{
						ID:     1,
						Name:   "test 16-1",
						Sql:    "DELETE FROM students WHERE graduation_year = 1900;",
						ErrMsg: "该SQL执行成功，等待事务提交",
						Count:  10,
					},
					{
						ID:     2,
						Name:   "test 16-2",
						Sql:    "DELETE FROM students WHERE graduation_year = 1900;",
						ErrMsg: "SQL执行失败，等待回滚",
						Count:  0,
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// get runner
			runner := &utils.Runner{}

			// mock modify orm
			mockORM := new(mockModifyOrmer)
			mockORM.mockData = tt.endpoint
			mockORM.mockRunner = runner
			e := &Executor{
				Ormer: mockORM,
			}

			inputParams := &utils.TransactionParamInfo{
				ID:      tt.endpoint.id,
				Retry:   tt.endpoint.id,
				Timeout: tt.endpoint.timeout,
				Name:    tt.endpoint.name,
				Sqls:    tt.endpoint.sqlInfo,
			}

			err := e.Modify(inputParams, runner)
			if (err != nil) != tt.wantErr {
				t.Errorf("Modify() error = %v, wantErr %v", err, tt.wantErr)
			}
			// 对比 getrunner 和 wantrunner
			if runner.ID != tt.wantRunner.ID {
				t.Errorf("Modify() runner.ID = %v, want %v", runner.ID, tt.wantRunner.ID)
			}
			if runner.Name != tt.wantRunner.Name {
				t.Errorf("Modify() runner.Name = %v, want %v", runner.Name, tt.wantRunner.Name)
			}
			if runner.Retry != tt.wantRunner.Retry {
				t.Errorf("Modify() runner.Retry = %v, want %v", runner.Retry, tt.wantRunner.Retry)
			}
			if runner.Timeout != tt.wantRunner.Timeout {
				t.Errorf("Modify() runner.Timeout = %v, want %v", runner.Timeout, tt.wantRunner.Timeout)
			}
			if !strings.Contains(runner.ErrMsg, tt.wantRunner.ErrMsg) {
				t.Errorf("Modify() runner.ErrMsg = %v, want %v", runner.ErrMsg, tt.wantRunner.ErrMsg)
			}
			if len(runner.SqlExecInfo) != len(tt.wantRunner.SqlExecInfo) {
				t.Errorf("Modify() runner.SqlExecInfo length = %v, want %v", len(runner.SqlExecInfo), len(tt.wantRunner.SqlExecInfo))
				return
			}
			for i, v := range runner.SqlExecInfo {
				if v.ID != tt.wantRunner.SqlExecInfo[i].ID {
					t.Errorf("Modify() runner.ID  = %v, want %v", v.ID, tt.wantRunner.SqlExecInfo[i].ID)
				}
				if v.Name != tt.wantRunner.SqlExecInfo[i].Name {
					t.Errorf("Modify() runner.Name  = %v, want %v", v.Name, tt.wantRunner.SqlExecInfo[i].Name)
				}
				if v.Sql != tt.wantRunner.SqlExecInfo[i].Sql {
					t.Errorf("Modify() runner.  = %v, want %v", v.Sql, tt.wantRunner.SqlExecInfo[i].Sql)
				}
				if v.Count != tt.wantRunner.SqlExecInfo[i].Count {
					t.Errorf("Modify() runner.  = %v, want %v", v.Count, tt.wantRunner.SqlExecInfo[i].Count)
				}
				if !strings.Contains(v.ErrMsg, tt.wantRunner.SqlExecInfo[i].ErrMsg) {
					t.Errorf("Modify() runner.  = %v, want contains %v", v.ErrMsg, tt.wantRunner.SqlExecInfo[i].ErrMsg)
				}
			}
		})
	}
}

type sqlsExecHistory struct {
	id            int    // 输入参数
	name          string // 输入参数
	sql           string // 输入参数sql
	isExecErr     bool   // SQL语句执行是否出错
	isRollback    bool   // 执行出错的时候是否成功回滚
	effectRow     int64  // 生效的行数
	resultIsNil   bool   // 执行结果是否为nil
	isRowEffect   bool   // 是否有行生效
	isNoRowsError bool   // 是否出现空行错误
}

type mockModifyEndpoint struct {
	factRetry       int               // 退出时的重试次数
	id              int               // 输入参数 id
	retry           int               // 输入参数 retry
	isOpenTxSuccess bool              // 是否打开事务成功
	isSuccess       bool              // 是否提交成功
	timeout         time.Duration     // 输入参数 超时时间
	name            string            // 输入参数 事务名
	count           int               // sql数量
	execHistory     []sqlsExecHistory // 事务中的SQL历史执行信息
	sqlInfo         []utils.SqlInfo   // 事务中的SQL历史执行信息
}

type option func(runner *utils.Runner)

func buildMockRunner(runner *utils.Runner, fn ...option) {
	for _, f := range fn {
		f(runner)
	}
}

func setMockRunnerId(mockData *mockModifyEndpoint) option {

	return func(runner *utils.Runner) {
		runner.ID = mockData.id
	}
}

func setMockRunnerName(mockData *mockModifyEndpoint) option {

	return func(runner *utils.Runner) {
		runner.Name = mockData.name
	}
}

func setMockRunnerRetry(mockData *mockModifyEndpoint) option {

	return func(runner *utils.Runner) {
		runner.Retry = mockData.factRetry
	}
}

func setMockRunnerCount(mockData *mockModifyEndpoint) option {

	return func(runner *utils.Runner) {
		runner.Count = mockData.count
	}
}

func setMockRunnerTimeout(mockData *mockModifyEndpoint) option {

	return func(runner *utils.Runner) {
		runner.Timeout = mockData.timeout
	}
}

var _ orm.Ormer = (*mockModifyOrmer)(nil)

type mockModifyOrmer struct {
	mockRunner *utils.Runner
	mockData   *mockModifyEndpoint
}

func (m *mockModifyOrmer) BeginWithCtx(ctx context.Context) (orm.TxOrmer, error) {
	buildFunc := []option{
		setMockRunnerId(m.mockData),
		setMockRunnerName(m.mockData),
		setMockRunnerCount(m.mockData),
		setMockRunnerRetry(m.mockData),
		setMockRunnerTimeout(m.mockData),
		setMockRunnerTimeout(m.mockData),
	}

	buildMockRunner(m.mockRunner, buildFunc...)

	if !m.mockData.isOpenTxSuccess {
		return nil, fmt.Errorf("mock 开启事务失败")
	}

	return &mockModifyTxOrmer{
		mockData: m.mockData,
	}, nil
}

var _ orm.TxOrmer = (*mockModifyTxOrmer)(nil)

type mockModifyTxOrmer struct {
	mockData           *mockModifyEndpoint
	mockExecHistoryNum int // 正在mock的事务的记录数
}

func (m *mockModifyTxOrmer) Raw(query string, args ...interface{}) orm.RawSeter {

	m.mockExecHistoryNum++

	return &modifyRawSeter{
		execHistroy:        &m.mockData.execHistory[m.mockExecHistoryNum-1],
		mockExecHistoryNum: m.mockExecHistoryNum - 1,
	}
}

func (m *mockModifyTxOrmer) Commit() error {
	if !m.mockData.isSuccess {
		return fmt.Errorf("mock 事务提交失败")
	}

	return nil
}

func (m *mockModifyTxOrmer) Rollback() error {
	if !m.mockData.execHistory[m.mockExecHistoryNum-1].isRollback {
		return errors.New("mock 事务回滚失败")
	}

	return nil
}

var _ orm.RawSeter = (*modifyRawSeter)(nil)

type modifyRawSeter struct {
	execHistroy        *sqlsExecHistory
	mockExecHistoryNum int
}

func (m *modifyRawSeter) Exec() (sql.Result, error) {

	if m.execHistroy.isExecErr && m.execHistroy.isRollback {
		return nil, fmt.Errorf("mock：SQL执行失败，事务回滚成功")
	}

	if m.execHistroy.isExecErr && !m.execHistroy.isRollback {
		return nil, fmt.Errorf("mock：SQL执行失败，事务回滚失败，等待自动回滚")
	}

	result := &modifyExecResult{
		mockData: m.execHistroy,
	}

	return result, nil
}

var _ sql.Result = (*modifyExecResult)(nil)

type modifyExecResult struct {
	mockData *sqlsExecHistory
}

func (m *modifyExecResult) RowsAffected() (int64, error) {

	if m.mockData.resultIsNil {
		return 0, nil
	}

	if !m.mockData.isNoRowsError {
		return 0, errors.New("出现空行错误")
	}

	if !m.mockData.isRowEffect {
		return 0, nil
	}

	return m.mockData.effectRow, nil
}

func (m *mockModifyOrmer) Read(md interface{}, cols ...string) error {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyOrmer) ReadWithCtx(ctx context.Context, md interface{}, cols ...string) error {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyOrmer) ReadForUpdate(md interface{}, cols ...string) error {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyOrmer) ReadForUpdateWithCtx(ctx context.Context, md interface{}, cols ...string) error {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyOrmer) ReadOrCreate(md interface{}, col1 string, cols ...string) (bool, int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyOrmer) ReadOrCreateWithCtx(ctx context.Context, md interface{}, col1 string, cols ...string) (bool, int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyOrmer) LoadRelated(md interface{}, name string, args ...utils2.KV) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyOrmer) LoadRelatedWithCtx(ctx context.Context, md interface{}, name string, args ...utils2.KV) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyOrmer) QueryM2M(md interface{}, name string) orm.QueryM2Mer {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyOrmer) QueryM2MWithCtx(ctx context.Context, md interface{}, name string) orm.QueryM2Mer {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyOrmer) QueryTable(ptrStructOrTableName interface{}) orm.QuerySeter {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyOrmer) QueryTableWithCtx(ctx context.Context, ptrStructOrTableName interface{}) orm.QuerySeter {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyOrmer) DBStats() *sql.DBStats {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyOrmer) Insert(md interface{}) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyOrmer) InsertWithCtx(ctx context.Context, md interface{}) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyOrmer) InsertOrUpdate(md interface{}, colConflitAndArgs ...string) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyOrmer) InsertOrUpdateWithCtx(ctx context.Context, md interface{}, colConflitAndArgs ...string) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyOrmer) InsertMulti(bulk int, mds interface{}) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyOrmer) InsertMultiWithCtx(ctx context.Context, bulk int, mds interface{}) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyOrmer) Update(md interface{}, cols ...string) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyOrmer) UpdateWithCtx(ctx context.Context, md interface{}, cols ...string) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyOrmer) Delete(md interface{}, cols ...string) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyOrmer) DeleteWithCtx(ctx context.Context, md interface{}, cols ...string) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyOrmer) Raw(query string, args ...interface{}) orm.RawSeter {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyOrmer) RawWithCtx(ctx context.Context, query string, args ...interface{}) orm.RawSeter {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyOrmer) Driver() orm.Driver {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyOrmer) Begin() (orm.TxOrmer, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyOrmer) BeginWithOpts(opts *sql.TxOptions) (orm.TxOrmer, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyOrmer) BeginWithCtxAndOpts(ctx context.Context, opts *sql.TxOptions) (orm.TxOrmer, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyOrmer) DoTx(task func(ctx context.Context, txOrm orm.TxOrmer) error) error {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyOrmer) DoTxWithCtx(ctx context.Context, task func(ctx context.Context, txOrm orm.TxOrmer) error) error {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyOrmer) DoTxWithOpts(opts *sql.TxOptions, task func(ctx context.Context, txOrm orm.TxOrmer) error) error {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyOrmer) DoTxWithCtxAndOpts(ctx context.Context, opts *sql.TxOptions, task func(ctx context.Context, txOrm orm.TxOrmer) error) error {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyTxOrmer) Read(md interface{}, cols ...string) error {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyTxOrmer) ReadWithCtx(ctx context.Context, md interface{}, cols ...string) error {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyTxOrmer) ReadForUpdate(md interface{}, cols ...string) error {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyTxOrmer) ReadForUpdateWithCtx(ctx context.Context, md interface{}, cols ...string) error {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyTxOrmer) ReadOrCreate(md interface{}, col1 string, cols ...string) (bool, int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyTxOrmer) ReadOrCreateWithCtx(ctx context.Context, md interface{}, col1 string, cols ...string) (bool, int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyTxOrmer) LoadRelated(md interface{}, name string, args ...utils2.KV) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyTxOrmer) LoadRelatedWithCtx(ctx context.Context, md interface{}, name string, args ...utils2.KV) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyTxOrmer) QueryM2M(md interface{}, name string) orm.QueryM2Mer {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyTxOrmer) QueryM2MWithCtx(ctx context.Context, md interface{}, name string) orm.QueryM2Mer {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyTxOrmer) QueryTable(ptrStructOrTableName interface{}) orm.QuerySeter {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyTxOrmer) QueryTableWithCtx(ctx context.Context, ptrStructOrTableName interface{}) orm.QuerySeter {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyTxOrmer) DBStats() *sql.DBStats {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyTxOrmer) Insert(md interface{}) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyTxOrmer) InsertWithCtx(ctx context.Context, md interface{}) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyTxOrmer) InsertOrUpdate(md interface{}, colConflitAndArgs ...string) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyTxOrmer) InsertOrUpdateWithCtx(ctx context.Context, md interface{}, colConflitAndArgs ...string) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyTxOrmer) InsertMulti(bulk int, mds interface{}) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyTxOrmer) InsertMultiWithCtx(ctx context.Context, bulk int, mds interface{}) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyTxOrmer) Update(md interface{}, cols ...string) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyTxOrmer) UpdateWithCtx(ctx context.Context, md interface{}, cols ...string) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyTxOrmer) Delete(md interface{}, cols ...string) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyTxOrmer) DeleteWithCtx(ctx context.Context, md interface{}, cols ...string) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyTxOrmer) RawWithCtx(ctx context.Context, query string, args ...interface{}) orm.RawSeter {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyTxOrmer) Driver() orm.Driver {
	// TODO implement me
	panic("implement me")
}

func (m *mockModifyTxOrmer) RollbackUnlessCommit() error {
	// TODO implement me
	panic("implement me")
}

func (m *modifyRawSeter) QueryRow(containers ...interface{}) error {
	// TODO implement me
	panic("implement me")
}

func (m *modifyRawSeter) QueryRows(containers ...interface{}) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *modifyRawSeter) SetArgs(i ...interface{}) orm.RawSeter {
	// TODO implement me
	panic("implement me")
}

func (m *modifyRawSeter) Values(container *[]orm.Params, cols ...string) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *modifyRawSeter) ValuesList(container *[]orm.ParamsList, cols ...string) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *modifyRawSeter) ValuesFlat(container *orm.ParamsList, cols ...string) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *modifyRawSeter) RowsToMap(result *orm.Params, keyCol, valueCol string) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *modifyRawSeter) RowsToStruct(ptrStruct interface{}, keyCol, valueCol string) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *modifyRawSeter) Prepare() (orm.RawPreparer, error) {
	// TODO implement me
	panic("implement me")
}

func (m *modifyExecResult) LastInsertId() (int64, error) {
	// TODO implement me
	panic("implement me")
}
