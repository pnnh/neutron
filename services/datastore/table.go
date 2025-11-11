package datastore

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"neutron/helpers"
	"neutron/models"
	"neutron/services/convert"
	"neutron/services/datetime"

	"github.com/sirupsen/logrus"
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

type IConvertDataRow interface {
	ToTableMap() (*DataRow, error)
}

type DataRow struct {
	dataMap map[string]interface{}
	Err     error
}

func NewDataRow() *DataRow {
	return &DataRow{
		dataMap: make(map[string]interface{}),
	}
}

func MapToDataRow(dataMap map[string]interface{}) *DataRow {
	return &DataRow{
		dataMap: dataMap,
	}
}

func (m *DataRow) setValue(key string, value interface{}) {
	if m.dataMap != nil {
		m.dataMap[key] = value
	}
}

func (m *DataRow) Keys() []string {
	keys := make([]string, 0, len(m.dataMap))
	for k := range m.dataMap {
		keys = append(keys, k)
	}
	return keys
}

func (m *DataRow) Values() []interface{} {
	keys := m.Keys()
	values := make([]interface{}, 0, len(m.dataMap))
	for _, k := range keys {
		values = append(values, m.dataMap[k])
	}
	return values
}

func (m *DataRow) ShallowCopyMap() map[string]interface{} {
	mapData := make(map[string]interface{})
	for k, v := range m.dataMap {
		mapData[k] = v
	}
	return mapData
}

func (m *DataRow) InnerMap() map[string]interface{} {
	return m.dataMap
}

func (m *DataRow) getValue(key string) (interface{}, bool) {
	if m.dataMap != nil {
		return nil, false
	}
	if v, ok := m.dataMap[key]; ok {
		return v, true

	}
	return nil, false
}

func (m *DataRow) MustGetInt(key string) int {
	if v, ok := m.dataMap[key]; ok {
		return v.(int)
	}
	panic(fmt.Sprintf("MustGetInt error, not found key: %s", key))
}

func (m *DataRow) TryGetInt(key string) (int, error) {
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

func (m *DataRow) GetInt(key string) int {
	intVal, err := m.TryGetInt(key)
	if err != nil {
		panic(fmt.Sprintf("GetInt error, key: %s, error: %v", key, err))
	}
	return intVal
}

func (m *DataRow) SetInt(key string, value int) {
	m.setValue(key, value)
}

func (m *DataRow) SetIntChain(key string, value int) *DataRow {
	m.setValue(key, value)
	return m
}

func (m *DataRow) SetNullInt(key string) {
	m.setValue(key, sql.NullInt64{Valid: false})
}

type IntGetter interface {
	TryGetInt(key string) (int, error)
}

func (m *DataRow) SetIntChainFrom(key string, getter IntGetter) *DataRow {
	if m.Err != nil {
		return m
	}
	value, err := getter.TryGetInt(key)
	if err != nil {
		m.Err = fmt.Errorf("TryGetInt error, key: %s, error: %v", key, err)
		return m
	}
	m.setValue(key, value)
	return m
}

func (m *DataRow) SetIntDefaultChainFrom(key string, getter IntGetter, defaultValue int) *DataRow {
	if m.Err != nil {
		return m
	}
	value, err := getter.TryGetInt(key)
	if err != nil {
		value = defaultValue
	}
	m.setValue(key, value)
	return m
}

func (m *DataRow) TryGetString(key string) (string, error) {
	v, ok := m.dataMap[key]
	if !ok {
		return "", fmt.Errorf("TryGetInt error, %w key: %s", models.ErrNotFound, key)
	}
	strVal, err := convert.ConvertString(v)
	if err != nil {
		return "", fmt.Errorf("TryGetString error, key: %s, value: %v, error: %w", key, v, err)
	}
	return strVal, nil
}

func (m *DataRow) GetString(key string) string {
	strVal, err := m.TryGetString(key)
	if err != nil {
		panic(fmt.Sprintf("GetString error, key: %s, error: %v", key, err))
	}
	return strVal
}

func (m *DataRow) GetStringOrDefault(key string, defaultValue string) string {
	strVal, err := m.TryGetString(key)
	if err != nil {
		if errors.Is(err, models.ErrNilValue) || errors.Is(err, models.ErrNotFound) {
			return defaultValue
		}
		panic(fmt.Sprintf("GetString error, key: %s, error: %v", key, err))
	}
	return strVal
}
func (m *DataRow) GetStringOrEmpty(key string) string {
	return m.GetStringOrDefault(key, "")
}

func (m *DataRow) GetNullString(key string) sql.NullString {
	strVal, err := m.TryGetString(key)
	if err != nil {
		if errors.Is(err, models.ErrNilValue) {
			return sql.NullString{Valid: false}
		}
		panic(fmt.Sprintf("GetString error, key: %s, error: %v", key, err))
	}
	return sql.NullString{String: strVal, Valid: true}
}

func (m *DataRow) SetString(key string, value string) {
	m.setValue(key, value)
}

type StringGetter interface {
	TryGetString(key string) (string, error)
	IsNullError(err error) bool
}

func (m *DataRow) SetStringChainFrom(key string, getter StringGetter) *DataRow {
	value, err := getter.TryGetString(key)
	if m.Err != nil {
		return m
	}
	if err != nil && !getter.IsNullError(err) {
		m.Err = fmt.Errorf("TryGetString error, key: %s, error: %v", key, err)
		return m
	}
	m.setValue(key, value)
	return m
}

func (m *DataRow) SetNullString(key string, value string) {
	strVal := sql.NullString{Valid: len(value) > 0, String: value}
	m.setValue(key, strVal)
}

func (m *DataRow) SetNullStringChain(key string, value string) *DataRow {
	m.SetNullString(key, value)
	return m
}

func (m *DataRow) SetNullStringChainFrom(key string, getter StringGetter) *DataRow {
	value, err := getter.TryGetString(key)
	if m.Err != nil {
		return m
	}
	if err != nil && !getter.IsNullError(err) {
		m.Err = fmt.Errorf("TryGetString error, key: %s, error: %v", key, err)
		return m
	}
	m.SetNullString(key, value)
	return m
}

func (m *DataRow) SetNullStringValue(key string, value sql.NullString) {
	m.setValue(key, value)
}

// SetNullUuidString 如果得到的是非法的UUID string或者空，者设置为sql.NullString{Valid: false}
func (m *DataRow) SetNullUuidString(key string, value string) {
	strVal := sql.NullString{Valid: helpers.IsUuid(value), String: value}
	m.setValue(key, strVal)
}

// SetNullUuidStringChainFrom 如果得到的是非法的UUID string或者空，者设置为sql.NullString{Valid: false}
func (m *DataRow) SetNullUuidStringChainFrom(key string, getter StringGetter) *DataRow {
	value, err := getter.TryGetString(key)
	if m.Err != nil {
		return m
	}
	if err != nil && !getter.IsNullError(err) {
		m.Err = fmt.Errorf("TryGetString error, key: %s, error: %v", key, err)
		return m
	}
	strVal := sql.NullString{Valid: helpers.IsUuid(value), String: value}
	m.setValue(key, strVal)
	return m
}

func (m *DataRow) SetNullUuidStringChain(key string, uuidString string) *DataRow {
	strVal := sql.NullString{Valid: helpers.IsUuid(uuidString), String: uuidString}
	m.setValue(key, strVal)
	return m
}

func (m *DataRow) TryGetTime(key string) (time.Time, error) {
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

func (m *DataRow) GetTime(key string) time.Time {
	timeVal, err := m.TryGetTime(key)
	if err != nil {
		panic(fmt.Sprintf("GetTime error, key: %s, error: %v", key, err))
	}
	return timeVal
}

func (m *DataRow) SetTime(key string, value time.Time) {
	m.setValue(key, value)
}

func (m *DataRow) SetTimeChain(key string, value time.Time) *DataRow {
	m.setValue(key, value)
	return m
}

func (m *DataRow) SetNullTime(key string, value time.Time) {
	m.setValue(key, sql.NullTime{Valid: value != datetime.NullTime, Time: value})
}

func (m *DataRow) SetNullTimeChain(key string, value time.Time) *DataRow {
	m.setValue(key, sql.NullTime{Valid: value != datetime.NullTime, Time: value})
	return m
}

func (m *DataRow) SetNullTimeValue(key string, value sql.NullTime) {
	m.setValue(key, value)
}

// SetNullTimeStringChainFrom tries to get a string from the getter, parse it as RFC3339 time, and set it as sql.NullTime
func (m *DataRow) SetNullTimeStringChainFrom(key string, getter StringGetter) *DataRow {
	value, err := getter.TryGetString(key)
	if m.Err != nil {
		return m
	}
	if err != nil {
		m.Err = fmt.Errorf("TryGetString error, key: %s, error: %v", key, err)
		return m
	}
	timeValue, err := time.Parse(time.RFC3339, value)
	if err != nil {
		logrus.Errorf("Parse time error, key: %s, value: %s, error: %v", key, value, err)
	}
	strVal := sql.NullTime{Valid: err == nil, Time: timeValue}
	m.setValue(key, strVal)
	return m
}

func (m *DataRow) GetTimeOrDefault(key string, defaultValue time.Time) time.Time {
	timeVal, err := m.TryGetTime(key)
	if err != nil {
		if errors.Is(err, models.ErrNilValue) {
			return defaultValue
		}
		panic(fmt.Sprintf("GetTime error, key: %s, error: %v", key, err))
	}
	return timeVal
}
