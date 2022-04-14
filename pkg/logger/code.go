package logger

// 定义对象
var err map[Code]Err

type Code uint

const (
	// 定义错误信息
	SUCCESS     Code = 10000 // 操作成功 info
	CHECK       Code = 10001 // 请检查 warn
	LIMITERROR  Code = 10002 // 超限 error
	TIMEOUT     Code = 10003 // 超时 error
	SUPWARN     Code = 10004 // 补充数据 warn
	STATUSERROR Code = 10005 // 非常规状态码 error
	EDERROR     Code = 10006 // 编码或解码失败 error
	RPERROR     Code = 10007 // 运行时参数错误 error
	PVERROR     Code = 10008 // 参数未通过校验 error
	DVERROR     Code = 10009 // 数据结果未通过校验 error
	RWARN       Code = 10010 // 运行时发生异常 warn
	RPWARN      Code = 10011 // 运行时发生错误 error
	CALLTIMEOUT Code = 20001 // 调用超时 error
	CALLERROR   Code = 20002 // 调用错误 error
	INTERRUPT   Code = 30001 // 组件中断 fatal
	// 定义默认值
	DEFAULT_ERROR_CODE = SUCCESS
)

// 初始化
func init() {
	err = map[Code]Err{
		SUCCESS:     {Code: SUCCESS, Level: INFO, Label: "操作成功"},
		CHECK:       {Code: CHECK, Level: WARN, Label: "请检查"},
		LIMITERROR:  {Code: LIMITERROR, Level: ERROR, Label: "超限"},
		TIMEOUT:     {Code: TIMEOUT, Level: ERROR, Label: "超时"},
		SUPWARN:     {Code: SUPWARN, Level: WARN, Label: "补充数据"},
		STATUSERROR: {Code: STATUSERROR, Level: ERROR, Label: "非常规状态码"},
		EDERROR:     {Code: EDERROR, Level: ERROR, Label: "编码或解码错误"},
		RPERROR:     {Code: RPERROR, Level: ERROR, Label: "运行时参数错误"},
		PVERROR:     {Code: PVERROR, Level: ERROR, Label: "参数未通过校验"},
		DVERROR:     {Code: DVERROR, Level: ERROR, Label: "数据结果未通过校验"},
		RWARN:       {Code: RWARN, Level: WARN, Label: "运行时发生异常"},
		RPWARN:      {Code: RPWARN, Level: WARN, Label: "运行时发生错误"},
		CALLTIMEOUT: {Code: CALLTIMEOUT, Level: ERROR, Label: "调用超时"},
		CALLERROR:   {Code: CALLERROR, Level: ERROR, Label: "调用错误"},
		INTERRUPT:   {Code: INTERRUPT, Level: FATAL, Label: "组件中断"},
	}
}
