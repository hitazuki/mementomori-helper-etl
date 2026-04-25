package parser

import (
	"mmth-etl/types"
)

// LogType 日志类型枚举
type LogType int

const (
	LogTypeNone           LogType = iota // 未知类型
	LogTypeDiamond                       // 钻石记录
	LogTypeCave                          // 洞穴记录
	LogTypeChallenge                     // 挑战记录
	LogTypeRuneTicket                    // 饼干记录
	LogTypeUpgradePanacea                // 红水记录
)

// IdentifyLogType 识别日志类型（一次扫描确定类型）
func IdentifyLogType(body string) LogType {
	// 检查物品变动记录（统一格式：Name: ItemName(Quality) × N）
	if types.DiamondRegex().MatchString(body) {
		return LogTypeDiamond
	}
	if types.RuneTicketRegex().MatchString(body) {
		return LogTypeRuneTicket
	}
	if types.UpgradePanaceaRegex().MatchString(body) {
		return LogTypeUpgradePanacea
	}

	// 检查洞穴记录
	if types.CaveEnterRegex().MatchString(body) || types.CaveFinishRegex().MatchString(body) || types.CaveErrorRegex.MatchString(body) {
		return LogTypeCave
	}

	// 检查挑战记录
	// 必须同时匹配：挑战模式（塔/任务）+ 结果关键词（成功/失败）
	// EN/TW: 以 Challenge/挑战 开头
	// JA/KO: 包含 挑戦/도전 关键词
	if types.ChallengeQuestRegex().MatchString(body) || types.ChallengeTowerRegex().MatchString(body) {
		if types.ChallengeSuccessRegex().MatchString(body) || types.ChallengeFailedRegex().MatchString(body) {
			return LogTypeChallenge
		}
	}

	return LogTypeNone
}
