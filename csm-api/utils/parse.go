package utils

import (
	"github.com/guregu/null"
	"strconv"
	"time"
)

func ParseNullString(s string) null.String {
	if s == "" {
		return null.NewString("", false)
	}
	return null.NewString(s, true)
}

func ParseNullInt(s string) null.Int {
	if s == "" {
		return null.NewInt(0, false)
	}

	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return null.NewInt(0, false)
	}
	return null.NewInt(i, true)
}

func ParseNullFloat(s string) null.Float {
	if s == "" {
		return null.NewFloat(0, false)
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return null.NewFloat(0, false)
	}
	return null.NewFloat(f, true)
}

func ParseNullDate(s string) null.Time {
	if s == "" {
		return null.NewTime(time.Time{}, false)
	}

	loc, err := time.LoadLocation("Asia/Seoul") // 명시적으로 KST 지정
	if err != nil {
		loc = time.FixedZone("KST", 9*60*60) // Fallback
	}

	t, err := time.ParseInLocation("2006-01-02", s, loc) // ← 여기 핵심
	if err != nil {
		return null.NewTime(time.Time{}, false)
	}

	return null.NewTime(t, true)
}

func ParseNullDateTime(dateStr, timeStr string) null.Time {
	if dateStr == "" || timeStr == "" {
		return null.NewTime(time.Time{}, false)
	}

	combined := dateStr + " " + timeStr
	loc := time.Now().Location()
	t, err := time.ParseInLocation("2006-01-02 15:04:05", combined, loc)
	if err != nil {
		return null.NewTime(time.Time{}, false)
	}

	return null.NewTime(t, true)
}
