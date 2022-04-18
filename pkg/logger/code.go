package logger

// 定义对象
var baseErrorIns map[Code]BaseError

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
	baseErrorIns = map[Code]BaseError{
		SUCCESS:     {code: SUCCESS, level: INFO, label: "操作成功"},
		CHECK:       {code: CHECK, level: WARN, label: "请检查"},
		LIMITERROR:  {code: LIMITERROR, level: ERROR, label: "超限"},
		TIMEOUT:     {code: TIMEOUT, level: ERROR, label: "超时"},
		SUPWARN:     {code: SUPWARN, level: WARN, label: "补充数据"},
		STATUSERROR: {code: STATUSERROR, level: ERROR, label: "非常规状态码"},
		EDERROR:     {code: EDERROR, level: ERROR, label: "编码或解码错误"},
		RPERROR:     {code: RPERROR, level: ERROR, label: "运行时参数错误"},
		PVERROR:     {code: PVERROR, level: ERROR, label: "参数未通过校验"},
		DVERROR:     {code: DVERROR, level: ERROR, label: "数据结果未通过校验"},
		RWARN:       {code: RWARN, level: WARN, label: "运行时发生异常"},
		RPWARN:      {code: RPWARN, level: WARN, label: "运行时发生错误"},
		CALLTIMEOUT: {code: CALLTIMEOUT, level: ERROR, label: "调用超时"},
		CALLERROR:   {code: CALLERROR, level: ERROR, label: "调用错误"},
		INTERRUPT:   {code: INTERRUPT, level: FATAL, label: "组件中断"},
	}
}

func GetBaseCode(code Code) BaseError {
	return baseErrorIns[code]
}

func GetDefaultCode() BaseError {
	return baseErrorIns[DEFAULT_ERROR_CODE]
}

func GetLevelBaseCode(code Code) LevelSupport {
	return baseErrorIns[code].level
}
