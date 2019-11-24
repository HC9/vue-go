package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

// 测试用户登录
func TestUserLogin(t *testing.T) {
	client := http.Client{}
	// 测试用例
	loginTest := []struct {
		status   int    // 登录成功或失败的状态
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		{0, "admin", "password"},
		{1, "abc", "password"},
	}
	respTest := []string{
		`{"code":20000,"data":{"token":"admin"},"msg":"登录成功","error":""}`,
		`{"code":51003,"data":null,"msg":"登录信息错误","error":""}`}
	for _, value := range loginTest {
		jsonBody, _ := json.Marshal(value)
		body := bytes.NewReader(jsonBody)
		resp, _ := client.Post("http://127.0.0.1:3000/api/v1/user/login", "application/json", body)
		all, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		result := fmt.Sprintf("%s", all)
		switch value.status {
		case 0:
			if result != respTest[0] {
				t.Errorf("\n\tInput:%+v\n\tOutput:%s\n\tExpect:%s\n", value, result, respTest[0])
			}
		case 1:
			if result != respTest[1] {
				t.Errorf("\n\tInput:%+v\n\tOutput:%s\n\tExpect:%s\n", value, result, respTest[1])
			}
		}
	}
}
