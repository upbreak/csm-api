package utils

import (
	"regexp"
	"strings"
	"time"
)

// 공백 제거 + 소문자 변환 + 특수문자 제거
func NormalizeForEqual(s string) string {
	re := regexp.MustCompile(`[^가-힣a-zA-Z0-9]`) // 한글, 영문, 숫자 외 제거
	return strings.ToLower(re.ReplaceAllString(strings.TrimSpace(s), ""))
}

// yyyy-mm-dd -> yy-mm-dd
func ConvertYYYYMMDDToYYMMDD(input string) string {
	t, _ := time.Parse("2006-01-02", strings.TrimSpace(input))
	return t.Format("06-01-02")
}

// mm-dd-yy -> yy-mm-dd
func ConvertMMDDYYToYYMMDD(input string) string {
	t, _ := time.Parse("01-02-06", strings.TrimSpace(input))
	return t.Format("06-01-02")
}
