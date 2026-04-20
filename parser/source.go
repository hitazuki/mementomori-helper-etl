package parser

import "strings"

// MapSourceToAlias 将来源日志映射为友好的来源名称
func MapSourceToAlias(source string) string {
	// 特殊映射
	if source == "You have triumphed." {
		return "temple of illusions"
	}

	// Gacha 日志：提取 " N times" 之前的内容
	if strings.HasPrefix(source, "Gacha ") {
		if idx := strings.Index(source, " times"); idx != -1 {
			beforeTimes := source[:idx]
			if lastSpace := strings.LastIndex(beforeTimes, " "); lastSpace != -1 {
				return strings.TrimSpace(beforeTimes[:lastSpace])
			}
		}
	}

	// Open 日志：提取 " x N" 之前的内容
	if strings.HasPrefix(source, "Open ") {
		if idx := strings.Index(source, " x "); idx != -1 {
			return strings.TrimSpace(source[:idx])
		}
	}

	return source
}

// IsValidSource 检查日志主体是否可作为来源
func IsValidSource(body string) bool {
	if strings.HasPrefix(body, "Name:") {
		return false
	}
	if strings.HasPrefix(body, "Challenge") {
		return false
	}
	if strings.HasPrefix(body, "OnError") {
		return false
	}
	return true
}
