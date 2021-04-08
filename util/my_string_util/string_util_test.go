package my_string_util

import (
	"fmt"
	"testing"
)

func TestStringUtil(t *testing.T) {
	num := IsNum("1231afs")
	fmt.Printf("num=%#v\n", num)
}