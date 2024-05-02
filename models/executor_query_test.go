package models

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/client/orm/clauses/order_clause"
	utils2 "github.com/beego/beego/v2/core/utils"
	"github.com/pkg/errors"
)

func TestExecutor_Query(t *testing.T) {

	tests := []struct {
		name           string
		sql            string
		retry          int // 最大允许重试次数
		realRetryTimes int // 重试多少次后返回
		data           []orm.Params
		ormReturnError error
		want           int64
		wantRetryTime  int // 查询的次数
		want2          []orm.Params
		wantErr        bool
	}{
		{
			name:           "查询成功",
			sql:            "SELECT * FROM user",
			retry:          1,
			realRetryTimes: 0,
			data: []orm.Params{
				{
					"a": "a",
					"b": "b",
					"c": "c",
				},
				{
					"a": "a",
					"b": "b",
					"c": "c",
				},
			},
			ormReturnError: nil,
			want:           2,
			wantRetryTime:  0,
			want2: []orm.Params{
				{
					"a": "a",
					"b": "b",
					"c": "c",
				},
				{
					"a": "a",
					"b": "b",
					"c": "c",
				},
			},
			wantErr: false,
		},
		{
			name:           "重试一次后依旧查询失败",
			sql:            "SELECT * FROM user",
			retry:          1,
			realRetryTimes: 1,
			data: []orm.Params{
				{
					"a": "a",
					"b": "b",
					"c": "c",
				},
				{
					"a": "a",
					"b": "b",
					"c": "c",
				},
			},
			ormReturnError: errors.New("mock error"),
			want:           0,
			wantRetryTime:  1,
			want2:          []orm.Params{},
			wantErr:        true,
		},
		{
			name:           "重试两次后依旧查询失败",
			sql:            "SELECT * FROM user",
			retry:          2,
			realRetryTimes: 2,
			data: []orm.Params{
				{
					"a": "a",
					"b": "b",
					"c": "c",
				},
				{
					"a": "a",
					"b": "b",
					"c": "c",
				},
			},
			ormReturnError: errors.New("mock error"),
			want:           0,
			wantRetryTime:  2,
			want2:          []orm.Params{},
			wantErr:        true,
		},
		{
			name:           "重试5次后依旧查询失败",
			sql:            "SELECT * FROM user",
			retry:          5,
			realRetryTimes: 5,
			data: []orm.Params{
				{
					"a": "a",
					"b": "b",
					"c": "c",
				},
				{
					"a": "a",
					"b": "b",
					"c": "c",
				},
			},
			ormReturnError: errors.New("mock error"),
			want:           0,
			wantRetryTime:  5,
			want2:          []orm.Params{},
			wantErr:        true,
		},
		{
			name:           "重试3次后成功",
			sql:            "SELECT * FROM user",
			retry:          5,
			realRetryTimes: 3,
			data: []orm.Params{
				{
					"a": "a",
					"b": "b",
					"c": "c",
				},
				{
					"a": "a",
					"b": "b",
					"c": "c",
				},
			},
			ormReturnError: nil,
			want:           2,
			wantRetryTime:  3,
			want2: []orm.Params{
				{
					"a": "a",
					"b": "b",
					"c": "c",
				},
				{
					"a": "a",
					"b": "b",
					"c": "c",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// orm mock
			expectData := tt.data

			mockORM := new(mockQueryOrmer)
			mockORM.mockData = make(map[string]*[]orm.Params)
			mockORM.mockData[tt.sql] = &expectData
			mockORM.realRetryTimes = tt.realRetryTimes
			mockORM.returnError = tt.ormReturnError
			mockORM.retry = 0

			exec := &Executor{
				mockORM,
			}

			got, got1, got2, err := exec.Query(tt.sql, tt.retry)
			if (err != nil) != tt.wantErr {
				t.Errorf("Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Query() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.wantRetryTime) {
				t.Errorf("Query() got1 = %v, want %v", got1, tt.wantRetryTime)
			}

			if (len(got2) != len(tt.want2)) && (got1 != tt.wantRetryTime) {
				t.Errorf("Query() got2 length = %v, want length %v", len(got2), len(tt.want2))
			}

			if len(got2) == len(tt.want2) {
				for i, v := range got2 {
					if (!reflect.DeepEqual(v, tt.want2[i])) != tt.wantErr {
						t.Errorf("Query() got2[i] = %v, want %v", got2, tt.want2)
						return
					}
				}
			}

		})
	}
}

var _ orm.Ormer = (*mockQueryOrmer)(nil)

type mockQueryOrmer struct {
	retry          int
	realRetryTimes int
	mockData       map[string]*[]orm.Params
	returnError    error
}

func (m *mockQueryOrmer) Raw(query string, args ...interface{}) orm.RawSeter {

	if m.realRetryTimes > m.retry {
		m.retry++
		return &queryRawSeter{
			rawSeterMockData: m.mockData[query],
			returnError:      errors.New("retry mock"),
			retry:            m.retry,
			realRetryTimes:   m.realRetryTimes,
		}
	}

	return &queryRawSeter{
		rawSeterMockData: m.mockData[query],
		returnError:      m.returnError,
		retry:            m.retry,
		realRetryTimes:   m.realRetryTimes,
	}
}

func (m *mockQueryOrmer) Begin() (orm.TxOrmer, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockQueryOrmer) BeginWithCtx(ctx context.Context) (orm.TxOrmer, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockQueryOrmer) BeginWithOpts(opts *sql.TxOptions) (orm.TxOrmer, error) {
	return nil, nil
}

func (m *mockQueryOrmer) BeginWithCtxAndOpts(ctx context.Context, opts *sql.TxOptions) (orm.TxOrmer, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockQueryOrmer) DoTx(task func(ctx context.Context, txOrm orm.TxOrmer) error) error {
	// TODO implement me
	panic("implement me")
}

func (m *mockQueryOrmer) DoTxWithCtx(ctx context.Context, task func(ctx context.Context, txOrm orm.TxOrmer) error) error {
	// TODO implement me
	panic("implement me")
}

func (m *mockQueryOrmer) DoTxWithOpts(opts *sql.TxOptions, task func(ctx context.Context, txOrm orm.TxOrmer) error) error {
	// TODO implement me
	panic("implement me")
}

func (m *mockQueryOrmer) DoTxWithCtxAndOpts(ctx context.Context, opts *sql.TxOptions, task func(ctx context.Context, txOrm orm.TxOrmer) error) error {
	// TODO implement me
	panic("implement me")
}

func (m *mockQueryOrmer) Read(md interface{}, cols ...string) error {
	// TODO implement me
	panic("implement me")
}

func (m *mockQueryOrmer) ReadWithCtx(ctx context.Context, md interface{}, cols ...string) error {
	// TODO implement me
	panic("implement me")
}

func (m *mockQueryOrmer) ReadForUpdate(md interface{}, cols ...string) error {
	return nil
}

func (m *mockQueryOrmer) ReadForUpdateWithCtx(ctx context.Context, md interface{}, cols ...string) error {
	return nil
}

func (m *mockQueryOrmer) ReadOrCreate(md interface{}, col1 string, cols ...string) (bool, int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockQueryOrmer) ReadOrCreateWithCtx(ctx context.Context, md interface{}, col1 string, cols ...string) (bool, int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockQueryOrmer) LoadRelated(md interface{}, name string, args ...utils2.KV) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockQueryOrmer) LoadRelatedWithCtx(ctx context.Context, md interface{}, name string, args ...utils2.KV) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockQueryOrmer) QueryM2M(md interface{}, name string) orm.QueryM2Mer {
	// TODO implement me
	panic("implement me")
}

func (m *mockQueryOrmer) QueryM2MWithCtx(ctx context.Context, md interface{}, name string) orm.QueryM2Mer {
	// TODO implement me
	panic("implement me")
}

func (m *mockQueryOrmer) QueryTable(ptrStructOrTableName interface{}) orm.QuerySeter {
	// TODO implement me
	panic("implement me")
}

func (m *mockQueryOrmer) QueryTableWithCtx(ctx context.Context, ptrStructOrTableName interface{}) orm.QuerySeter {
	// TODO implement me
	panic("implement me")
}

func (m *mockQueryOrmer) DBStats() *sql.DBStats {
	// TODO implement me
	panic("implement me")
}

func (m *mockQueryOrmer) Insert(md interface{}) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockQueryOrmer) InsertWithCtx(ctx context.Context, md interface{}) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockQueryOrmer) InsertOrUpdate(md interface{}, colConflitAndArgs ...string) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockQueryOrmer) InsertOrUpdateWithCtx(ctx context.Context, md interface{}, colConflitAndArgs ...string) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockQueryOrmer) InsertMulti(bulk int, mds interface{}) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockQueryOrmer) InsertMultiWithCtx(ctx context.Context, bulk int, mds interface{}) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockQueryOrmer) Update(md interface{}, cols ...string) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockQueryOrmer) UpdateWithCtx(ctx context.Context, md interface{}, cols ...string) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockQueryOrmer) Delete(md interface{}, cols ...string) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockQueryOrmer) DeleteWithCtx(ctx context.Context, md interface{}, cols ...string) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockQueryOrmer) RawWithCtx(ctx context.Context, query string, args ...interface{}) orm.RawSeter {
	// TODO implement me
	panic("implement me")
}

func (m *mockQueryOrmer) Driver() orm.Driver {
	// TODO implement me
	panic("implement me")
}

var _ orm.RawSeter = (*queryRawSeter)(nil)

type queryRawSeter struct {
	retry            int
	realRetryTimes   int
	rawSeterMockData *[]orm.Params
	returnError      error
}

func (r *queryRawSeter) Values(results *[]orm.Params, exprs ...string) (int64, error) {

	if r.returnError != nil {
		return 0, r.returnError
	}

	*results = *r.rawSeterMockData

	return int64(len(*results)), nil
}

func (r *queryRawSeter) ValuesFlat(container *orm.ParamsList, cols ...string) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) Exec() (sql.Result, error) {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) QueryRow(containers ...interface{}) error {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) QueryRows(containers ...interface{}) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) SetArgs(i ...interface{}) orm.RawSeter {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) Prepare() (orm.RawPreparer, error) {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) Filter(s string, i ...interface{}) orm.QuerySeter {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) FilterRaw(s string, s2 string) orm.QuerySeter {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) Exclude(s string, i ...interface{}) orm.QuerySeter {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) SetCond(condition *orm.Condition) orm.QuerySeter {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) GetCond() *orm.Condition {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) Limit(limit interface{}, args ...interface{}) orm.QuerySeter {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) Offset(offset interface{}) orm.QuerySeter {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) GroupBy(exprs ...string) orm.QuerySeter {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) OrderBy(exprs ...string) orm.QuerySeter {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) OrderClauses(orders ...*order_clause.Order) orm.QuerySeter {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) ForceIndex(indexes ...string) orm.QuerySeter {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) UseIndex(indexes ...string) orm.QuerySeter {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) IgnoreIndex(indexes ...string) orm.QuerySeter {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) RelatedSel(params ...interface{}) orm.QuerySeter {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) Distinct() orm.QuerySeter {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) ForUpdate() orm.QuerySeter {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) Count() (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) CountWithCtx(ctx context.Context) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) Exist() bool {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) ExistWithCtx(ctx context.Context) bool {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) Update(values orm.Params) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) UpdateWithCtx(ctx context.Context, values orm.Params) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) Delete() (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) DeleteWithCtx(ctx context.Context) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) PrepareInsert() (orm.Inserter, error) {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) PrepareInsertWithCtx(ctx context.Context) (orm.Inserter, error) {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) All(container interface{}, cols ...string) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) AllWithCtx(ctx context.Context, container interface{}, cols ...string) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) One(container interface{}, cols ...string) error {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) OneWithCtx(ctx context.Context, container interface{}, cols ...string) error {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) ValuesWithCtx(ctx context.Context, results *[]orm.Params, exprs ...string) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) ValuesList(results *[]orm.ParamsList, exprs ...string) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) ValuesListWithCtx(ctx context.Context, results *[]orm.ParamsList, exprs ...string) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) z(result *orm.ParamsList, expr string) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) ValuesFlatWithCtx(ctx context.Context, result *orm.ParamsList, expr string) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) RowsToMap(result *orm.Params, keyCol, valueCol string) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) RowsToStruct(ptrStruct interface{}, keyCol, valueCol string) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (r *queryRawSeter) Aggregate(s string) orm.QuerySeter {
	// TODO implement me
	panic("implement me")
}
