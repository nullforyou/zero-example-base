package custom_validate

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

// MobileTag 手机号验证
func MobileTag(fl validator.FieldLevel) bool{
	regRuler := "^1\\d{10}$"
	reg := regexp.MustCompile(regRuler)
	return reg.MatchString(fl.Field().String())
}

