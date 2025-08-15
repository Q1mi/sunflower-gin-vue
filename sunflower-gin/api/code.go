package api

// ResCode 定义返回码类型
type ResCode int64

// 定义一些返回码示例,可根据业务需求自定义
const (
	CodeSuccess         ResCode = 0
	CodeInvalidParam    ResCode = 4000
	CodeUserExist       ResCode = 4010
	CodeUserNotExist    ResCode = 4011
	CodeInvalidPassword ResCode = 4020

	CodeNeedLogin    ResCode = 4100
	CodeInvalidToken ResCode = 4200

	CodeServerBusy ResCode = 5000
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "请求参数错误",
	CodeUserExist:       "用户名已存在",
	CodeUserNotExist:    "用户名不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy:      "服务繁忙",

	CodeNeedLogin:    "需要登录",
	CodeInvalidToken: "无效的token",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
