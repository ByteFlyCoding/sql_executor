package utils

const (
	SUCCESSQUERY    = iota // 查询接口执行任务成功
	FAILQUERY              // 查询接口执行任务失败
	SUCCESSMODIFY          // 修改接口执行任务成功
	FAILMODIFYEXIST        // 修改接口存在运行异常的子任务
	PARAMETERERROR         // 传入参数存在异常
)
