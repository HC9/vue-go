package test

import (
	"testing"
	"vgo/model"
	"vgo/service"
)

// 测试用于登录
func TestLoginCheck(t *testing.T) {
	NameEmailsTests := []struct {
		input    service.UserRegisterService
		expected service.Response
	}{
		{
			service.UserRegisterService{UserName: "abc", Email: "emmhY3zkU3@qq.com"},
			service.Response{Code: 52001, Msg: "注册校验失败", Error: "emmhY3zkU3@qq.com have registered"},
		},
		{
			service.UserRegisterService{UserName: "emmhY3zkU3", Email: "emmhY3zkU3111@qq.com"},
			service.Response{Code: 52001, Msg: "注册校验失败", Error: "emmhY3zkU3 have registered"},
		},
		{
			service.UserRegisterService{UserName: "efg", Email: "abciowe@ip.com"},
			service.Response{Code: 20000, Msg: "注册校验成功"},
		},
	}

	for _, tt := range NameEmailsTests {
		actual := model.CheckNameAndEmail(&tt.input)
		if *actual != tt.expected {
			t.Errorf("CheckNameAndEmail(%v) = \n\t%v; \nexpected %v", tt.input, actual, tt.expected)
		}
	}

}

// 测试注册时单个字段的检查
func TestNameOrEmail(t *testing.T) {
	tests := []struct {
		input    model.CheckFieldRepeat
		expected service.Response
	}{
		{
			model.CheckFieldRepeat{Name: "email", Value: "I59oSfwILza9@qq.com"},
			service.Response{Code: 20000, Msg: "检查通过"},
		},
		{
			model.CheckFieldRepeat{Name: "email", Value: "I59oSfwILz@qq.com"},
			service.Response{Code: 52001, Error: "I59oSfwILz@qq.com 已使用", Msg: "请更换"},
		},
		{
			model.CheckFieldRepeat{Name: "username", Value: "Y9dYF7lnVa"},
			service.Response{Code: 52001, Error: "Y9dYF7lnVa 已使用", Msg: "请更换"},
		},
		{
			model.CheckFieldRepeat{Name: "username", Value: "hello123"},
			service.Response{Code: 20000, Msg: "检查通过"},
		},
	}

	for _, test := range tests {
		actual := test.input.NameOrEmail()
		if *actual != test.expected {
			t.Errorf("CheckNameAndEmail(%v) = \n\t%v; \nexpected %v", test.input, actual, test.expected)
		}
	}
}
