package models

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type NECode int

func (c NECode) String() string {
	return fmt.Sprintf("%d", c)
}

func (c NECode) WithMessage(message string) *NECommonResult {
	return NENewCommonResult(c, NECodeMessage(c)+":"+message, nil)
}

func (c NECode) WithData(data interface{}) *NECommonResult {
	return NENewCommonResult(c, NECodeMessage(c), data)
}

func (c NECode) WithError(err error) *NECommonResult {
	logrus.Errorf("NECode.WithError [%d] %v", c, err)
	return NENewCommonResult(c, NECodeMessage(c), nil)
}

const (
	NECodeOk               NECode = 200
	NECodeError            NECode = 500
	NECodeAccountExists    NECode = 600 // 账号已存在
	NECodeAccountNotExists NECode = 601 // 账号不存在
	NECodeNotLogin         NECode = 602
	NECodeInvalidParameter NECode = 603
	NECodeNotFound         NECode = 404
	NEStatusAccountExists  NECode = 607 // 账号已存在
	NECodeInvalidParams    NECode = 609 // 参数无效
	NECodeUnauthorized     NECode = 401 // 未授权
)

func NECodeMessage(code NECode) string {
	switch code {
	case NECodeOk:
		return "成功"
	case NECodeNotFound:
		return "资源未找到"
	case NECodeAccountNotExists:
		return "账号不存在"
	case NECodeNotLogin:
		return "尚未登陆"
	case NEStatusAccountExists:
		return "账号已存在"
	}
	return fmt.Sprintf("未知错误：%d", code)
}
