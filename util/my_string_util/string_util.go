package my_string_util

import (
	"unicode"
)

//判断是否是数字
func IsNum(str string) bool  {
	isNum := true
	//_,r 是 index 和 int32 unicode
	for _, r := range str {
		if !unicode.IsDigit(r) {
			isNum = false
			break
		}
	}
	return isNum
}
