package entity

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"
)

/**
 * @author 작성자: 김진우
 * @created 작성일: 2025-02-14
 * @modified 최종 수정일:
 * @modifiedBy 최종 수정자:
 * @modified description
 * -
 */

// func: (일반 타입 → SQLNulls 타입 변환)
// @param
// - value: 일반 타입
// @ex
// - nullInt64 := ToSQLNulls(42).(sql.NullInt64)
func ToSQLNulls[T any](value T) any {
	switch v := any(value).(type) {
	case string:
		return sql.NullString{String: v, Valid: v != ""}
	case int:
		return sql.NullInt64{Int64: int64(v), Valid: v != 0}
	case int64:
		return sql.NullInt64{Int64: v, Valid: v != 0}
	case float64:
		return sql.NullFloat64{Float64: v, Valid: v != 0}
	case bool:
		return sql.NullBool{Bool: v, Valid: true} // false도 유효한 값
	case time.Time:
		return sql.NullTime{Time: v, Valid: !v.IsZero()}
	default:
		return nil
	}
}

// func: (SQLNulls 타입 → 일반 타입 변환)
// @Generic
// - T any: 원하는 일반 타입
// @param
// - nullValue: SQLNulls 타입
// @ex
// - regularStr := ToRegular[string](nullValue)
// - regularInt := ToRegular[int64](nullValue)
// - regularFloat := ToRegular[float64](nullValue)
// - regularBool := ToRegular[bool](nullValue)
// - regularTime := ToRegular[time.Time](nullValue)
func ToRegular[T any](nullValue any) T {
	switch v := nullValue.(type) {
	case sql.NullString:
		if v.Valid {
			return any(v.String).(T)
		}
		return any("").(T)
	case sql.NullInt64:
		if v.Valid {
			return any(v.Int64).(T)
		}
		return any(int64(0)).(T)
	case sql.NullFloat64:
		if v.Valid {
			return any(v.Float64).(T)
		}
		return any(float64(0)).(T)
	case sql.NullBool:
		if v.Valid {
			return any(v.Bool).(T)
		}
		return any(false).(T)
	case sql.NullTime:
		if v.Valid {
			return any(v.Time).(T)
		}
		return any(time.Time{}).(T)
	default:
		var zero T
		return zero
	}
}

// func: (단일 일반 타입 구조체 → SQLNulls 타입 구조체 변환)
// @param
// 첫 번째 매개변수(regular): 일반 타입 (포인터 X)
// 두 번째 매개변수(sqlNulls): SQLNulls 타입 구조체 (포인터 O)
func ConvertToSQLNulls(regular any, sqlNulls any) error {
	regularVal := reflect.ValueOf(regular)
	sqlNullsVal := reflect.ValueOf(sqlNulls)

	// sqlNulls가 포인터인지 확인하고 역참조
	if sqlNullsVal.Kind() != reflect.Ptr {
		return fmt.Errorf("sqlNulls must be a pointer to a struct")
	}
	sqlNullsVal = sqlNullsVal.Elem()

	// regular가 포인터가 아니라 값 타입이면 그대로 사용
	if regularVal.Kind() != reflect.Struct {
		return fmt.Errorf("regular must be a struct, got %s", regularVal.Kind())
	}

	regularType := regularVal.Type()

	for i := 0; i < regularVal.NumField(); i++ {
		fieldName := regularType.Field(i).Name
		regularField := regularVal.Field(i)
		sqlNullsField := sqlNullsVal.FieldByName(fieldName)

		if !sqlNullsField.IsValid() {
			return fmt.Errorf("field %q exists in regular but not in SQLNulls", fieldName)
		}

		// 값을 직접 수정
		switch regularField.Kind() {
		case reflect.Int64:
			sqlNullsField.Set(reflect.ValueOf(sql.NullInt64{
				Int64: regularField.Int(),
				Valid: regularField.Int() != 0,
			}))
		case reflect.String:
			sqlNullsField.Set(reflect.ValueOf(sql.NullString{
				String: regularField.String(),
				Valid:  regularField.String() != "",
			}))
		case reflect.Struct:
			if regularField.Type() == reflect.TypeOf(time.Time{}) {
				sqlNullsField.Set(reflect.ValueOf(sql.NullTime{
					Time:  regularField.Interface().(time.Time),
					Valid: !regularField.Interface().(time.Time).IsZero(),
				}))
			} else {
				return fmt.Errorf("unsupported struct field type for field %q", fieldName)
			}
		default:
			return fmt.Errorf("unsupported field type %s for field %q", regularField.Type(), fieldName)
		}
	}

	return nil
}

// func: (단일 SQLNulls 타입 구조체 → 일반 타입 구조체 변환)
// @param
// 첫 번째 매개변수(sqlNulls): SQLNulls 타입 구조체 (포인터 X)
// 두 번째 매개변수(regular): 일반 타입 구조체 (포인터 O)
func ConvertToRegular(sqlNulls any, regular any) error {
	sqlNullsVal := reflect.ValueOf(sqlNulls)
	regularVal := reflect.ValueOf(regular).Elem()

	// sqlNulls가 포인터일 경우 역참조
	if sqlNullsVal.Kind() == reflect.Ptr {
		sqlNullsVal = sqlNullsVal.Elem()
	}

	// sqlNulls가 구조체인지 확인
	if sqlNullsVal.Kind() != reflect.Struct {
		return fmt.Errorf("sqlNulls must be a struct, got %s", sqlNullsVal.Kind())
	}

	sqlNullsType := sqlNullsVal.Type()

	for i := 0; i < sqlNullsVal.NumField(); i++ {
		fieldName := sqlNullsType.Field(i).Name
		sqlNullsField := sqlNullsVal.Field(i)
		regularField := regularVal.FieldByName(fieldName)

		if !regularField.IsValid() {
			return fmt.Errorf("field %q exists in SQLNulls but not in regular", fieldName)
		}

		switch sqlNullsField.Interface().(type) {
		case sql.NullInt64:
			regularField.SetInt(sqlNullsField.Interface().(sql.NullInt64).Int64)
		case sql.NullString:
			regularField.SetString(sqlNullsField.Interface().(sql.NullString).String)
		case sql.NullTime:
			regularField.Set(reflect.ValueOf(sqlNullsField.Interface().(sql.NullTime).Time))
		default:
			return fmt.Errorf("unsupported field type %s for field %q", sqlNullsField.Type(), fieldName)
		}
	}
	return nil
}

// func: (일반 타입 구조체 슬라이스 → SQLNulls 타입 구조체 슬라이스 변환)
// @param
// 첫 번째 매개변수(regularSlice): 일반 타입 슬라이스 (포인터 X)
// 두 번째 매개변수(sqlNullsSlice): SQLNulls 타입 슬라이스 (포인터 O)
func ConvertSliceToSQLNulls(regularSlice any, sqlNullsSlice any) error {
	regularVal := reflect.ValueOf(regularSlice)
	sqlNullsVal := reflect.ValueOf(sqlNullsSlice)

	if regularVal.Kind() != reflect.Slice {
		return fmt.Errorf("regularSlice must be a slice (got %s)", regularVal.Kind())
	}

	if sqlNullsVal.Kind() == reflect.Ptr {
		sqlNullsVal = sqlNullsVal.Elem()
	}

	if sqlNullsVal.Kind() != reflect.Slice {
		return fmt.Errorf("sqlNullsSlice must be a pointer to a slice (got %s)", sqlNullsVal.Kind())
	}

	if sqlNullsVal.IsNil() || sqlNullsVal.Len() != regularVal.Len() {
		sqlNullsVal.Set(reflect.MakeSlice(sqlNullsVal.Type(), regularVal.Len(), regularVal.Len()))
	}

	for i := 0; i < regularVal.Len(); i++ {
		regularItem := regularVal.Index(i)

		// regularItem이 포인터면 값을 가져와서 전달
		if regularItem.Kind() == reflect.Ptr {
			regularItem = regularItem.Elem()
		}

		sqlNullsItem := reflect.New(sqlNullsVal.Type().Elem().Elem())

		err := ConvertToSQLNulls(regularItem.Interface(), sqlNullsItem.Interface()) // 이제 값이 전달됨
		if err != nil {
			return fmt.Errorf("error converting item at index %d: %w", i, err)
		}

		sqlNullsVal.Index(i).Set(sqlNullsItem)
	}

	return nil
}

// func: (SQLNulls 타입 구조체 슬라이스 → 일반 타입 구조체 슬라이스 변환)
// @param
// 첫 번째 매개변수(sqlNullsSlice): SQLNulls 타입 슬라이스 (포인터 X)
// 두 번째 매개변수(regularSlice): 일반 타입 슬라이스 (포인터 O)
func ConvertSliceToRegular(sqlNullsSlice any, regularSlice any) error {
	sqlNullsVal := reflect.ValueOf(sqlNullsSlice)
	regularVal := reflect.ValueOf(regularSlice)

	// sqlNullsSlice가 포인터라면 실제 값을 가져오기
	if sqlNullsVal.Kind() == reflect.Ptr {
		sqlNullsVal = sqlNullsVal.Elem()
	}

	if sqlNullsVal.Kind() != reflect.Slice {
		return fmt.Errorf("sqlNullsSlice must be a slice (got %s)", sqlNullsVal.Kind())
	}

	if regularVal.Kind() != reflect.Ptr {
		return fmt.Errorf("regularSlice must be a pointer (got %s)", regularVal.Kind())
	}

	regularSliceElem := regularVal.Elem()

	if regularSliceElem.Kind() != reflect.Slice {
		return fmt.Errorf("regularSlice must be a pointer to a slice (got %s)", regularSliceElem.Kind())
	}

	if regularSliceElem.IsNil() || regularSliceElem.Len() != sqlNullsVal.Len() {
		regularSliceElem.Set(reflect.MakeSlice(regularSliceElem.Type(), sqlNullsVal.Len(), sqlNullsVal.Len()))
	}

	for i := 0; i < sqlNullsVal.Len(); i++ {
		sqlNullsItem := sqlNullsVal.Index(i)
		regularItem := reflect.New(regularSliceElem.Type().Elem().Elem())

		err := ConvertToRegular(sqlNullsItem.Interface(), regularItem.Interface())
		if err != nil {
			return fmt.Errorf("error converting item at index %d: %w", i, err)
		}

		regularSliceElem.Index(i).Set(regularItem)
	}

	return nil
}
