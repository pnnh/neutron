package datastore

import (
	"fmt"
	"neutron/services/strutil"
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
