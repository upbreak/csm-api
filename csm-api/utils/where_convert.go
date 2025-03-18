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

func OrTimeBetweenWhereConvert(condition string, sqlValue1 sql.NullString, sqlValue2 sql.NullString, target string) string {
	if sqlValue1.Valid && sqlValue2.Valid {
		value1 := strings.TrimSpace(sqlValue1.String)
		value2 := strings.TrimSpace(sqlValue2.String)

		if value1 != "" && value2 != "" {
			condition += fmt.Sprintf(` OR TO_CHAR(%s, 'YYYY-MM-DD') BETWEEN '%s' AND '%s'`, target, value1, value2)
		}
	}
	return condition
}

func RetrySearchTextConvert(retry string, columns []string) string {
	where := ""
	if retry == "" {
		return ""
	}
	keyArr := strings.Split(retry, "~")

	if len(keyArr) > 0 {
		for _, key := range keyArr {
			arr := strings.Split(key, ":")
			if arr[0] == "ALL" {
				values := strings.Split(arr[1], "|")
				for _, value := range values {
					var temp []string
					for _, column := range columns {
						s := fmt.Sprintf(`LOWER(%s) LIKE LOWER('%%%s%%')`, column, strings.TrimSpace(value))
						temp = append(temp, s)
					}
					where += fmt.Sprintf(`AND (%s)`, strings.Join(temp, " OR "))
				}
			} else {
				for _, column := range columns {
					values := strings.Split(arr[1], "|")
					if strings.Contains(column, arr[0]) {
						for _, value := range values {
							where += fmt.Sprintf(`AND LOWER(%s) LIKE LOWER('%%%s%%')`, column, strings.TrimSpace(value))
						}
						break
					}
				}
			}
		}
	}

	return where
}
