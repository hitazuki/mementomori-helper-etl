package parser

import (
	"encoding/json"
	"regexp"
	"strings"
)

// ParsedLog 表示解析后的日志结构
type ParsedLog struct {
	RawLine     string // 原始行内容
	Timestamp   string // 时间戳（日志内容中的，秒级精度）
	PreciseTime string // 精确时间（JSON time字段，纳秒级精度，用于检查点）
	Character   string // 角色名
	Body        string // 日志主体
	IsValid     bool   // 格式是否有效
}

// ParseLogLine 解析单行日志，提取通用信息
func ParseLogLine(line string) ParsedLog {
	result := ParsedLog{RawLine: line}

	var logEntry struct {
		Log  string `json:"log"`
		Time string `json:"time"`
	}
	if err := json.Unmarshal([]byte(line), &logEntry); err != nil {
		return result
	}

	logContent := logEntry.Log
	result.PreciseTime = logEntry.Time

	timestamp, character, body := extractLogComponents(logContent)
	if timestamp == "" {
		return result
	}

	result.Timestamp = timestamp
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
