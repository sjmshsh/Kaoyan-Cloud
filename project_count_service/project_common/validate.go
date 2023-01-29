package project_common

import "regexp"

// VerifyMobile 验证手机合法性
func VerifyMobile (mobile string) bool {
	if mobile == "" {
		return false
	}
	regular := "^(13[0-9]|14[579]|15[0-3,5-9]|16[6]|17[0135678]|18[0-9]|19[89])\\d{8}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobile)
}
