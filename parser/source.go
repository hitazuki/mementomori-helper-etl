package parser

import (
	"mmth-etl/types"
	"strings"
)

// IsValidSource 检查日志主体是否可作为来源上下文
// 过滤掉物品变动日志（Name:）、挑战日志（Challenge）和错误日志（OnError）
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

	return true
}
