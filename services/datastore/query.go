package datastore

import (
	"database/sql"
	"fmt"
	"github.com/iancoleman/strcase"
	"neutron/services/strutil"
	"reflect"
	"strings"
)

func IsValidTableName(tableName string) bool {
	return strutil.IsValidName(tableName)
}

func NewGetQuery(tableName string, whereText, orderText, extraText string,
	sqlParams map[string]any) (tabMap *TableMap, getErr error) {
	if !IsValidTableName(tableName) {
		return nil, fmt.Errorf("invalid table name: %s", tableName)
	}
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf(" select * from %s ", tableName))
	if whereText != "" {
		builder.WriteString(" where " + whereText)
	}
	if orderText != "" {
		builder.WriteString(" order by " + orderText)
	}
	if extraText != "" {
		builder.WriteString(" " + extraText)
	}
	pageSqlText := builder.String()

	rows, err := NamedQuery(pageSqlText, sqlParams)
	if err != nil {
		return nil, fmt.Errorf("NewSelectQuery: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			getErr = fmt.Errorf("rows.Close: %w", closeErr)
		}
	}()

	var firstMap map[string]interface{}
	// Iterate over the rows and scan each into a map
	for rows.Next() {
		rowMap := make(map[string]interface{})
		if err := rows.MapScan(rowMap); err != nil {
			return nil, fmt.Errorf("MapScan: %w", err)
		}
		firstMap = rowMap
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	if firstMap == nil {
		return nil, nil
	}
	tableMap := ConvertToTableMap(firstMap)

	return tableMap, nil
}

func ReflectColumns(s interface{}) (map[string]any, error) {
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)

	// 如果是指针，取其元素
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("getStructFields kind is not struct")
	}
	columnMap := make(map[string]any)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		//logrus.Printf("成员名: %s, 类型: %s, 类型名称: %s, 值: %v\n",
		//	field.Name,
		//	field.Type.Kind(),      // 基本类型，如 string、int、struct 等
		//	field.Type.Name(),      // 类型名称，如 int、string、自定义类型名
		//	fieldValue.Interface(), // 字段值
		//)
		colName := strcase.ToSnake(field.Name)
		dbTag := field.Tag.Get("db")
		if dbTag == "-" {
			continue
		} else if dbTag != "" {
			colName = dbTag
		}
		insertTag := field.Tag.Get("insert")
		if insertTag == "skip" {
			continue
		}
		var colValue any
		switch val := fieldValue.Interface().(type) {
		case sql.NullString:
			if val.Valid {
				colValue = val.String
			}
		default:
			colValue = val
		}
		columnMap[colName] = colValue
	}
	return columnMap, nil
}
