package xerr

var message map[uint32]string


func init() {
	message = make(map[uint32]string)
	message[SUCCESS] = "SUCCESS"
	message[ErrorValidation] = "参数错误"
	message[ErrorServerCommon] = "服务器开小差啦,稍后再来试一试"
	message[ErrorTokenExpire] = "token失效，请重新登陆"
	message[ErrorTokenGenerate] = "生成token失败"
	message[ErrorDb] = "数据库繁忙,请稍后再试"
	message[ErrorNotFound] = "未找到对象"
	message[ErrorBusiness] = "业务错误"
}

func MapMsg(code uint32) string {
	if msg, ok := message[code]; ok {
		return msg
	} else {
		return "服务器开小差啦,稍后再来试一试"
	}
}

func IsCodeMsg(code uint32) bool {
	if _, ok := message[code]; ok {
		return true
	} else {
		return false
	}
}