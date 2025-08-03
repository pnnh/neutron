package datastore

import (
	"database/sql"
	"errors"
	"fmt"
	"neutron/services/convert"
	"time"
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

type IConvertTableMap interface {
	ToTableMap() (*TableMap, error)
}

type TableMap struct {
	dataMap map[string]interface{}
}

func NewTableMap() *TableMap {
	return new(TableMap)
}

func ConvertToTableMap(dataMap map[string]interface{}) *TableMap {
	return &TableMap{
		dataMap: dataMap,
	}
}

func (m *TableMap) Set(key string, value interface{}) {
	if m.dataMap == nil {
		m.dataMap[key] = value
	}
}

func (m *TableMap) Keys() []string {
	keys := make([]string, 0, len(m.dataMap))
	for k := range m.dataMap {
		keys = append(keys, k)
	}
	return keys
}

func (m *TableMap) Values() []interface{} {
	values := make([]interface{}, 0, len(m.dataMap))
	for _, v := range m.dataMap {
		values = append(values, v)
	}
	return values
}

func (m *TableMap) Get(key string) (interface{}, bool) {
	if m.dataMap == nil {
		return nil, false
	}
	if v, ok := m.dataMap[key]; ok {
		return v, true

	}
	return nil, false
}

func (m *TableMap) MustGetInt(key string) int {
	if v, ok := m.dataMap[key]; ok {
		return v.(int)
	}
	panic(fmt.Sprintf("MustGetInt error, not found key: %s", key))
}

func (m *TableMap) TryGetInt(key string) (int, error) {
	v, ok := m.dataMap[key]
	if !ok {
		return 0, fmt.Errorf("TryGetInt error, not found key: %s", key)
	}
	intVal, err := convert.ConvertInt(v)
	if err != nil {
		return 0, fmt.Errorf("TryGetInt error, key: %s, value: %v, error: %w", key, v, err)
	}
	return intVal, nil
}

func (m *TableMap) GetInt(key string) int {
	intVal, err := m.TryGetInt(key)
	if err != nil {
		panic(fmt.Sprintf("GetInt error, key: %s, error: %v", key, err))
	}
	return intVal
}

func (m *TableMap) TryGetString(key string) (string, error) {
	v, ok := m.dataMap[key]
	if !ok {
		return "", fmt.Errorf("TryGetInt error, not found key: %s", key)
	}
	strVal, err := convert.ConvertString(v)
	if err != nil {
		return "", fmt.Errorf("TryGetString error, key: %s, value: %v, error: %w", key, v, err)
	}
	return strVal, nil
}

func (m *TableMap) GetString(key string) string {
	strVal, err := m.TryGetString(key)
	if err != nil {
		panic(fmt.Sprintf("GetString error, key: %s, error: %v", key, err))
	}
	return strVal
}

func (m *TableMap) GetNullString(key string) sql.NullString {
	strVal, err := m.TryGetString(key)
	if err != nil {
		if errors.Is(err, convert.ErrNilValue) {
			return sql.NullString{Valid: false}
		}
		panic(fmt.Sprintf("GetString error, key: %s, error: %v", key, err))
	}
	return sql.NullString{String: strVal, Valid: true}
}

func (m *TableMap) TryGetTime(key string) (time.Time, error) {
	v, ok := m.dataMap[key]
	if !ok {
		return time.Time{}, fmt.Errorf("TryGetTime error, not found key: %s", key)
	}
	timeVal, err := convert.ConvertTime(v)
	if err != nil {
		return time.Time{}, fmt.Errorf("TryGetTime error, key: %s, value: %v, error: %w", key, v, err)
	}
	return timeVal, nil
}

func (m *TableMap) GetTime(key string) time.Time {
	timeVal, err := m.TryGetTime(key)
	if err != nil {
		panic(fmt.Sprintf("GetTime error, key: %s, error: %v", key, err))
	}
	return timeVal
}
