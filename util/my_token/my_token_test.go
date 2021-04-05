package my_token

import (
	"fmt"
	"testing"
)

func TestToken(t *testing.T) {
	token, err := GetToken(13414321)
	if err != nil {
		fmt.Printf("err=%#v\n", err)
	}
	fmt.Printf("string=%s\n", token)
	parseToken, obj, err := ParseToken(token)
	fmt.Printf("parseToken=%#v\n", parseToken)
	fmt.Printf("obj=%#v\n", obj.UserId)
	if err != nil {
		fmt.Printf("string=%s\n", "parse token")
	}
}