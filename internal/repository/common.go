package repository

import (
	"errors"
	"fmt"
	"strings"
)

const (
	CommonFieldID        = "id"
	CommonFieldCreatedAt = "created_at"
	CommonFieldUpdatedAt = "updated_at"
)

var ErrNoRecordFound = errors.New("no rows in result set")

func BuildSaveQuery(tableName string, fields []string, isUpdate bool) string {
	if isUpdate {
		var setFields []string
		for _, field := range fields {
			if field == CommonFieldID || field == CommonFieldCreatedAt {
				continue
			}
			setFields = append(setFields, fmt.Sprintf("%s=:%s", field, field))
		}
		return fmt.Sprintf("UPDATE %s SET %s WHERE %s=:%s", tableName, strings.Join(setFields, ", "), CommonFieldID, CommonFieldID)
	}
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, strings.Join(fields, ", "), ":"+strings.Join(fields, ", :"))
}
