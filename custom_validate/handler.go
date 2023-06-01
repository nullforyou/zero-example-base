package custom_validate

import (
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

type Validator struct {
	Validate *validator.Validate //验证器
	Trans ut.Translator //翻译器
}

// InitValidator 初始化验证器
func InitValidator() (val Validator) {
	zh := zh.New()
	uni := ut.New(zh, zh)
	trans, _ := uni.GetTranslator("zh")
	validate := validator.New()
	zh_translations.RegisterDefaultTranslations(validate, trans)
	RegisterValidator(validate, trans) //注册自定义验证
	return Validator{Validate: validate, Trans: trans}
}


// RegisterValidator 注册自定义验证
func RegisterValidator(validate *validator.Validate, trans ut.Translator) {
	RegisterMobileValidator(validate, trans)
}

func RegisterMobileValidator(validate *validator.Validate, trans ut.Translator) {
	validate.RegisterValidation("mobile", MobileTag)
	validate.RegisterTranslation("mobile", trans, func(ut ut.Translator) error {
		return ut.Add("mobile", "{0} 格式不正确!", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("mobile", fe.Field())
		return t
	})
}