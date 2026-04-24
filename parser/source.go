package parser

import (
	"mmth-etl/types"
	"strings"
)

// IsValidSource 检查日志主体是否可作为来源上下文
// 过滤掉物品变动日志（Name:）、挑战日志（Challenge）、错误日志（OnError）和系统异常日志
func IsValidSource(body string) bool {
	mgr := types.GetI18nManager()

	// 检查是否以 Name: 开头（当前语言）
	if mgr.IsNameLabelLine(body) {
		return false
	}

	// 检查是否以 Challenge 开头（当前语言）
	if mgr.IsChallengeLine(body) {
		return false
	}

	// 检查 OnError（语言无关）
	if strings.HasPrefix(body, "OnError") {
		return false
	}

	// 检查系统异常日志（System. 开头）
	// 这类日志可能是洞窟记录，但不应作为物品变动的来源上下文
	if strings.HasPrefix(body, "System.") {
		return false
	}

	return true
}

// ShouldClearSource 检查是否需要清空来源上下文
// 系统/错误日志需要清空，Name/Challenge 日志保留当前来源
func ShouldClearSource(body string) bool {
	// OnError 日志需要清空来源
	if strings.HasPrefix(body, "OnError") {
		return true
	}

	// 系统异常日志需要清空来源
	if strings.HasPrefix(body, "System.") {
		return true
	}

	return false
}
