package my_token

import (
	"fmt"
	"testing"
)

func TestToken(t *testing.T) {
 	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VyIjp7InVzZXJfaWQiOjEzMjQxNCwidXNlcm5hbWUiOiJhZG1pbiIsImVtYWlsIjoiYWRtaW5AYm1mdC5jb20iLCJuaWNrbmFtZSI6IumYv-aWh-aWh-WQliIsImF2YXRhciI6Im5vbmUifSwiZXhwIjoxNjE4NDU3NDM4LCJpYXQiOjE2MTc4NTI2MzgsImlzcyI6ImZ0Y2xvdWQiLCJzdWIiOiJ1c2VyIG15X3Rva2VuIn0.Q0-lvr0CTXsF2hVZly0z8l5OBz33RwkpsYmF8swQw5A"
	user, err := GetUser(token)
	fmt.Printf("user=%#v\n", user)
	fmt.Printf("err=%#v\n", err)
}