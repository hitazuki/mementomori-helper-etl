package parser

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"
)

// ParsedLog 表示解析后的日志结构
type ParsedLog struct {
	RawLine   string // 原始行内容
	Timestamp string // 时间戳（秒级精度，ISO格式）
	Character string // 角色名
	Body      string // 日志主体
	IsValid   bool   // 格式是否有效
}

// ParseLogLine 解析单行日志，提取通用信息
// 支持两种格式：
//  1. Docker JSON 格式: {"log":"...", "time":"..."}
//  2. 纯文本格式: [时间戳] [角色] 内容
// 统一使用秒级精度时间戳（ISO格式，本地时区）
func ParseLogLine(line string) ParsedLog {
	result := ParsedLog{RawLine: line}

	// 确定日志内容来源
	logContent := line
	var logEntry struct {
		Log string `json:"log"`
	}
	if err := json.Unmarshal([]byte(line), &logEntry); err == nil && logEntry.Log != "" {
		logContent = logEntry.Log
	}

	// 提取通用字段
	timestamp, character, body := extractLogComponents(logContent)
	if timestamp == "" {
		return result
	}

	// 解析时间戳并使用本地时区
	if t, err := time.ParseInLocation("2006-01-02 15:04:05", timestamp, time.Local); err == nil {
		result.Timestamp = t.Format("2006-01-02T15:04:05Z07:00")
	} else {
		// 解析失败时使用原始时间戳
		result.Timestamp = strings.Replace(timestamp, " ", "T", 1)
	}
	result.Character = character
	result.Body = body
	result.IsValid = true

	return result
}

// extractLogComponents 从日志内容中提取时间戳、角色名和日志主体
func extractLogComponents(logContent string) (string, string, string) {
	logContent = strings.TrimSpace(logContent)

	if !strings.HasPrefix(logContent, "[") {
		return "", "", ""
	}

	timeRegex := regexp.MustCompile(`^\[(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})\]`)
	timeMatches := timeRegex.FindStringSubmatch(logContent)
	if len(timeMatches) < 2 {
		return "", "", ""
	}
	timestamp := timeMatches[1]

	firstCloseIdx := strings.Index(logContent, "]")
	if firstCloseIdx == -1 || firstCloseIdx+1 >= len(logContent) {
		return "", "", ""
	}

	remaining := strings.TrimSpace(logContent[firstCloseIdx+1:])

	if !strings.HasPrefix(remaining, "[") {
		return "", "", ""
	}

	secondCloseIdx := strings.Index(remaining, "]")
	if secondCloseIdx == -1 || secondCloseIdx <= 1 {
		return "", "", ""
	}

	nameSection := remaining[1:secondCloseIdx]
	character := extractCharacterName(nameSection)

	if secondCloseIdx+1 >= len(remaining) {
		return "", "", ""
	}
	body := strings.TrimSpace(remaining[secondCloseIdx+1:])

	if character == "" || body == "" {
		return "", "", ""
	}

	return timestamp, character, body
}

// extractCharacterName 从名字部分提取角色名（去掉等级）
func extractCharacterName(nameSection string) string {
	nameSection = strings.TrimSpace(nameSection)

	if idx := strings.Index(nameSection, "(Lv"); idx != -1 {
		return strings.TrimSpace(nameSection[:idx])
	}

	return nameSection
}
