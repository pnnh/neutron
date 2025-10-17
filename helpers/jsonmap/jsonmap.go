package jsonmap

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"neutron/helpers"
	"neutron/models"
	"neutron/services/convert"
	"neutron/services/datetime"
)

type IIntValue interface {
	int | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64
}

type IMapValue interface {
	IIntValue | string
}

type MapValue[T IMapValue] struct {
	Value T
	Error error
}

func NewMapValue[T IMapValue](value T) *MapValue[T] {
	return &MapValue[T]{value, nil}
}

func NewErrorMapValue[T IMapValue](err error) *MapValue[T] {
	return &MapValue[T]{Error: err}
}

func HasError(errs ...error) error {
	for _, e := range errs {
		if e != nil {
			return e
		}
	}
	return nil
}

type JsonMap struct {
	dataMap map[string]interface{}
	Err     error
}

func NewJsonMap() *JsonMap {
	return &JsonMap{
		dataMap: make(map[string]interface{}),
	}
}

func ConvertJsonMap(dataMap map[string]interface{}) *JsonMap {
	return &JsonMap{
		dataMap: dataMap,
	}
}

func (m *JsonMap) setValue(key string, value interface{}) {
	if m.dataMap != nil {
		m.dataMap[key] = value
	}
}

func (m *JsonMap) Keys() []string {
	keys := make([]string, 0, len(m.dataMap))
	for k := range m.dataMap {
		keys = append(keys, k)
	}
	return keys
}

func (m *JsonMap) Values() []interface{} {
	keys := m.Keys()
	values := make([]interface{}, 0, len(m.dataMap))
	for _, k := range keys {
		values = append(values, m.dataMap[k])
	}
	return values
}

func (m *JsonMap) InnerMap() map[string]interface{} {
	return m.dataMap
}

func (m *JsonMap) InnerMapPtr() *map[string]interface{} {
	return &m.dataMap
}

func (m *JsonMap) getValue(key string) (interface{}, bool) {
	if m.dataMap != nil {
		return nil, false
	}
	if v, ok := m.dataMap[key]; ok {
		return v, true

	}
	return nil, false
}

func (m *JsonMap) IsNullError(err error) bool {
	return models.IsErrNilValue(err)
}

func (m *JsonMap) MustGetInt(key string) int {
	if v, ok := m.dataMap[key]; ok {
		return v.(int)
	}
	panic(fmt.Sprintf("MustGetInt error, not found key: %s", key))
}

func (m *JsonMap) TryGetInt(key string) (int, error) {
	v, ok := m.dataMap[key]
	if !ok {
		return 0, models.ErrNilValue
	}
	intVal, err := convert.ConvertInt(v)
	if err != nil {
		return 0, fmt.Errorf("TryGetInt error, key: %s, value: %v, error: %w", key, v, err)
	}
	return intVal, nil
}

func (m *JsonMap) GetInt(key string) int {
	intVal, err := m.TryGetInt(key)
	if err != nil {
		panic(fmt.Sprintf("GetInt error, key: %s, error: %v", key, err))
	}
	return intVal
}

func (m *JsonMap) SetInt(key string, value int) {
	m.setValue(key, value)
}

func (m *JsonMap) SetNullInt(key string) {
	m.setValue(key, sql.NullInt64{Valid: false})
}

func (m *JsonMap) TryGetString(key string) (string, error) {
	v, ok := m.dataMap[key]
	if !ok {
		return "", models.ErrNilValue
	}
	strVal, err := convert.ConvertString(v)
	if err != nil {
		return "", fmt.Errorf("TryGetString error, key: %s, value: %v, error: %w", key, v, err)
	}
	return strVal, nil
}

// WillGetString 尝试获取一个字符串值。它会先检查JsonMap当前是否有错误，如果有错误直接返回空字符串。否则执行TryGetString方法获取字符串值。
// 如果TryGetString返回错误，则将该错误记录在JsonMap的Err字段中，并返回空字符串。
// 这种设计允许调用者连续多次调用WillGetString，而不必每次都检查错误，从而简化了错误处理逻辑。
func (m *JsonMap) WillGetString(key string) string {
	if m.Err != nil {
		return ""
	}
	strVal, err := m.TryGetString(key)
	if err != nil {
		m.Err = fmt.Errorf("WillGetString error, key: %s, error: %v", key, err)
	}
	return strVal
}

func (m *JsonMap) GetString(key string) string {
	strVal, err := m.TryGetString(key)
	if err != nil {
		panic(fmt.Sprintf("GetString error, key: %s, error: %v", key, err))
	}
	return strVal
}

func (m *JsonMap) GetStringOrDefault(key string, defaultValue string) string {
	strVal, err := m.TryGetString(key)
	if err != nil {
		if errors.Is(err, models.ErrNilValue) {
			return defaultValue
		}
		panic(fmt.Sprintf("GetString error, key: %s, error: %v", key, err))
	}
	return strVal
}

func (m *JsonMap) GetNullString(key string) sql.NullString {
	strVal, err := m.TryGetString(key)
	if err != nil {
		if errors.Is(err, models.ErrNilValue) {
			return sql.NullString{Valid: false}
		}
		panic(fmt.Sprintf("GetString error, key: %s, error: %v", key, err))
	}
	return sql.NullString{String: strVal, Valid: true}
}

func (m *JsonMap) SetString(key string, value string) {
	m.setValue(key, value)
}

func (m *JsonMap) SetNullString(key string, value string) {
	strVal := sql.NullString{Valid: len(value) > 0, String: value}
	m.setValue(key, strVal)
}

func (m *JsonMap) SetNullStringValue(key string, value sql.NullString) {
	m.setValue(key, value)
}

func (m *JsonMap) SetNullUuidString(key string, value string) {
	strVal := sql.NullString{Valid: helpers.IsUuid(value), String: value}
	m.setValue(key, strVal)
}

func (m *JsonMap) TryGetTime(key string) (time.Time, error) {
	v, ok := m.dataMap[key]
	if !ok {
		return time.Time{}, models.ErrNilValue
	}
	timeVal, err := convert.ConvertTime(v)
	if err != nil {
		return time.Time{}, fmt.Errorf("TryGetTime error, key: %s, value: %v, error: %w", key, v, err)
	}
	return timeVal, nil
}

func (m *JsonMap) GetTime(key string) time.Time {
	timeVal, err := m.TryGetTime(key)
	if err != nil {
		panic(fmt.Sprintf("GetTime error, key: %s, error: %v", key, err))
	}
	return timeVal
}

func (m *JsonMap) SetTime(key string, value time.Time) {
	m.setValue(key, value)
}

func (m *JsonMap) SetNullTime(key string, value time.Time) {
	m.setValue(key, sql.NullTime{Valid: value.After(datetime.UtcMinTime), Time: value})
}

func (m *JsonMap) SetNullTimeValue(key string, value sql.NullTime) {
	m.setValue(key, value)
}

func (m *JsonMap) GetTimeOrDefault(key string, defaultValue time.Time) time.Time {
	timeVal, err := m.TryGetTime(key)
	if err != nil {
		if errors.Is(err, models.ErrNilValue) {
			return defaultValue
		}
		panic(fmt.Sprintf("GetTime error, key: %s, error: %v", key, err))
	}
	return timeVal
}

//
//
//func GetInt[T IIntValue](m *JsonMap, key string) *MapValue[T] {
//	if v, ok := (*m).innerMap[key]; ok {
//		switch v := v.(type) {
//		case int:
//			return NewMapValue((T)(v))
//		case uint:
//			return NewMapValue((T)(v))
//		case int8:
//			return NewMapValue((T)(v))
//		case uint8:
//			return NewMapValue((T)(v))
//		case int16:
//			return NewMapValue((T)(v))
//		case uint16:
//			return NewMapValue((T)(v))
//		case int32:
//			return NewMapValue((T)(v))
//		case uint32:
//			return NewMapValue((T)(v))
//		case int64:
//			return NewMapValue((T)(v))
//		case uint64:
//			return NewMapValue((T)(v))
//		case float32:
//			return NewMapValue((T)(v))
//		case float64:
//			return NewMapValue((T)(v))
//		}
//	}
//	err := fmt.Errorf("MustGetInt error, not found key: %s", key)
//	return NewErrorMapValue[T](err)
//}
