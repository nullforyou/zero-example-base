package response

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-base/utils/xerr"
	"google.golang.org/grpc/status"
	"net/http"
)

type Bean struct {
	Code    uint32      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Response 聚合应答
func Response(r *http.Request, w http.ResponseWriter, resp interface{}, err error) {
	if err != nil {
		FailedResponse(r,w, resp, err)
	} else {
		SuccessResponse(r, w, resp)
	}
}

// SuccessResponse 成功的应答
func SuccessResponse(r *http.Request, w http.ResponseWriter, resp interface{}) {
	httpx.WriteJson(w, http.StatusOK, &Bean{
		Code:    xerr.SUCCESS,
		Message: xerr.MapMsg(xerr.SUCCESS),
		Data:    resp,
	})
}

// FailedResponse 失败的应答
func FailedResponse(r *http.Request, w http.ResponseWriter, resp interface{}, err error) {
	failedBean := Bean{}
	causeErr := errors.Cause(err)
	if e, ok := causeErr.(*xerr.BusinessError); ok { //自定义错误类型
		failedBean.Code = e.GetErrCode()
		failedBean.Message = e.GetErrMsg()
		if failedBean.Message == "" {
			failedBean.Message = xerr.MapMsg(failedBean.Code)
		}
	} else {
		if grpcStatus, ok := status.FromError(causeErr); ok { //grpc错误
			grpcCode := uint32(grpcStatus.Code())
			if xerr.IsCodeMsg(grpcCode) { //区分自定义错误跟系统底层、db等错误，底层、db错误不能返回给前端
				failedBean.Code = grpcCode
				failedBean.Message = grpcStatus.Message()
			} else {
				failedBean.Code = xerr.ErrorRpcOther
				failedBean.Message = xerr.MapMsg(failedBean.Code)
			}
		} else {
			//其他错误
		}
	}

	logx.WithContext(r.Context()).Errorf("【API-ERR】 : %+v ", err)

	switch failedBean.Code {
	case xerr.ErrorTokenExpire:
		httpx.WriteJson(w, http.StatusUnauthorized, failedBean)
	case xerr.ErrorNotFound:
		httpx.WriteJson(w, http.StatusNotFound, failedBean)
	default:
		httpx.WriteJson(w, http.StatusBadRequest, failedBean)
	}
}

// ValidateErrOrResponse 验证错误的应答
func ValidateErrOrResponse(r *http.Request, w http.ResponseWriter, err error, trans ut.Translator) {
	var msg string
	var data interface{}
	causeErr := errors.Cause(err)
	if _, ok := causeErr.(validator.ValidationErrors); ok {
		validateErrs := err.(validator.ValidationErrors)
		for _, validateErr := range validateErrs {
			msg = validateErr.Translate(trans)
			break
		}
		data = validateErrs.Translate(trans)
	} else {
		msg = err.Error()
	}

	httpx.WriteJson(w, http.StatusUnprocessableEntity, &Bean{
		Code:    xerr.ErrorValidation,
		Message: msg,
		Data: data,
	})
}