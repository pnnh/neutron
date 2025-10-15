package models

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type NECode int

func (c NECode) String() string {
	return fmt.Sprintf("%d", c)
}

// Deprecated: Use WithLocalMessage instead
func (c NECode) WithMessage(message string) *NECommonResult {
	return NENewCommonResult(c, NECodeMessage(LangZh, c)+":"+message, nil)
}

func (c NECode) WithLocalMessage(lang, zhMsg, enMsg string) *NECommonResult {
	fullMsg := NECodeMessage(lang, c)
	if lang == LangZh {
		fullMsg = fmt.Sprintf("%s -> %s", fullMsg, zhMsg)
	} else {
		fullMsg = fmt.Sprintf("%s -> %s", fullMsg, enMsg)
	}
	return NENewCommonResult(c, fullMsg, nil)
}

// Deprecated: Use WithLocalData instead
func (c NECode) WithData(data interface{}) *NECommonResult {
	return NENewCommonResult(c, NECodeMessage(LangZh, c), data)
}

func (c NECode) WithLocalData(lang string, data interface{}) *NECommonResult {
	return NENewCommonResult(c, NECodeMessage(lang, c), data)
}

func (c NECode) WithError(err error) *NECommonResult {
	logrus.Errorf("NECode.WithError [%d] %v", c, err)
	return NENewCommonResult(c, NECodeMessage(LangEn, c), nil)
}

func (c NECode) WithLocalError(lang string, err error, zhMsg, enMsg string) *NECommonResult {
	logrus.Errorf("MCode.WithError [%d] %v", c, err)
	return c.WithLocalMessage(lang, zhMsg, enMsg)
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

func NECodeMessage(lang string, code NECode) string {

	if lang == LangZh {
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
		default:
			return fmt.Sprintf("未知错误：%d", code)
		}
	}
	switch code {
	case NECodeOk:
		return "Success"
	case NECodeNotFound:
		return "Resource not found"
	case NECodeAccountNotExists:
		return "Account does not exist"
	case NECodeNotLogin:
		return "Not logged in"
	case NEStatusAccountExists:
		return "Account already exists"
	default:
		return fmt.Sprintf("Unknown error: %d", code)
	}
}
