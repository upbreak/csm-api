package utils

import (
	"database/sql"
	"fmt"
	"strings"
)

func StringWhereConvert(condition string, sqlValue sql.NullString, target string) string {
	if sqlValue.Valid {
		value := strings.TrimSpace(sqlValue.String)
		if value != "" {
			condition += fmt.Sprintf(` AND LOWER(%s) LIKE LOWER('%%%s%%')`, target, value)
		}
	}
	return condition
}

func Int64WhereConvert(condition string, sqlValue sql.NullInt64, target string) string {
	if sqlValue.Valid {
		value := sqlValue.Int64
		if value != 0 {
			condition += fmt.Sprintf(` AND %s = %d`, target, value)
		}
	}
	return condition
}

func TimeWhereConvert(condition string, sqlValue sql.NullString, target string) string {
	if sqlValue.Valid {
		value := strings.TrimSpace(sqlValue.String)
		if value != "" {
			condition += fmt.Sprintf(` AND TO_CHAR(%s, 'YYYY-MM-DD') = '%s'`, target, value)
		}
	}
	return condition
}

func TimeBetweenWhereConvert(condition string, sqlValue1 sql.NullString, sqlValue2 sql.NullString, target string) string {
	if sqlValue1.Valid && sqlValue2.Valid {
		value1 := strings.TrimSpace(sqlValue1.String)
		value2 := strings.TrimSpace(sqlValue2.String)

		if value1 != "" && value2 != "" {
			condition += fmt.Sprintf(` AND TO_CHAR(%s, 'YYYY-MM-DD') BETWEEN '%s' AND '%s'`, target, value1, value2)
		}
	}
	return condition
}
