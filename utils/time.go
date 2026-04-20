package utils

import (
	"fmt"
	"time"
)

// GetWeekString 获取周字符串（年-周）
func GetWeekString(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-W%02d", year, week)
}

// GetMonthString 获取月字符串（年-月）
func GetMonthString(t time.Time) string {
	year := t.Year()
	month := t.Month()
	return fmt.Sprintf("%d-%02d", year, month)
}
