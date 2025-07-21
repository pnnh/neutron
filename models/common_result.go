package models

import "github.com/sirupsen/logrus"

type NECommonResult struct {
	Code    NECode      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NENewCommonResult(code NECode, message string, data interface{}) *NECommonResult {
	return &NECommonResult{Code: code, Message: message, Data: data}
}

type NESelectResponse struct {
	Page  int   `json:"page"`
	Size  int   `json:"size"`
	Count int   `json:"count"`
	Range []any `json:"range"`
}

type NEViewModel interface {
	ToViewModel() interface{}
}

type NESelectResult[T NEViewModel] struct {
	Page  int `json:"page"`
	Size  int `json:"size"`
	Count int `json:"count"`
	Range []T `json:"range"`
}

func NEModelListToViewList[T NEViewModel](models []T) []any {
	result := make([]any, 0)
	for _, model := range models {
		result = append(result, model.ToViewModel())
	}
	return result
}

func NESelectResultToResponse[T NEViewModel](result *NESelectResult[T]) *NESelectResponse {
	if result == nil {
		return nil
	}
	return &NESelectResponse{
		Page:  result.Page,
		Size:  result.Size,
		Count: result.Count,
		Range: NEModelListToViewList(result.Range),
	}
}

func NEParseCommonResult(data interface{}) *NECommonResult {
	if data == nil {
		return nil
	}
	if result, ok := data.(*NECommonResult); ok {
		return result
	}
	return nil
}

func (r *NECommonResult) SetCode(code NECode) *NECommonResult {
	r.Code = code
	return r
}

func (r *NECommonResult) SetMessage(message string) *NECommonResult {
	r.Message = message
	return r
}

func (r *NECommonResult) SetData(data interface{}) *NECommonResult {
	r.Data = data
	return r
}

// 成功时响应
func NESuccessResult(data any) *NECommonResult {
	return NENewCommonResult(NECodeOk, NECodeMessage(NECodeOk), data)
}

// 成功时响应，附加消息
func NESuccessResultMessage(data any, message string) *NECommonResult {
	return NENewCommonResult(NECodeOk, message, data)
}

// 错误时响应
func NEErrorResult(err error) *NECommonResult {
	return NEErrorResultFull(err, NECodeError, "", nil)
}

// 错误时响应，指定错误提示消息
func NEErrorResultMessage(err error, message string) *NECommonResult {
	return NEErrorResultFull(err, NECodeError, message, err)
}

// 错误时响应，指定响应码
func NEErrorResultCode(err error, bizCode NECode) *NECommonResult {
	return NEErrorResultFull(err, bizCode, "", nil)
}

// 错误时响应，附带额外数据
func NEErrorResultFull(err error, bizCode NECode, message string, data any) *NECommonResult {
	logrus.Errorf("NECode.WithError [%d] %v", bizCode, err)
	if message == "" {
		message = NECodeMessage(bizCode)
	}
	return NENewCommonResult(bizCode, message, data)
}
