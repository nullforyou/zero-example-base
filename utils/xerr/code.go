package xerr



const SUCCESS uint32 = 200
const ErrorServerCommon uint32 = 100000 //服务器开小差啦,稍后再来试一试
const ErrorValidation uint32 = 100001 //验证错误
const ErrorTokenExpire uint32 = 100002
const ErrorTokenGenerate uint32 = 100003
const ErrorDb uint32 = 100004


//业务所涉状态大于200000
const ErrorBusiness uint32 = 200000
const ErrorNotFound = 200001
const ErrorRpcOther uint32 = 200002
