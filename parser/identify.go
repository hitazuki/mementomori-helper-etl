package parser

import (
	"mmth-etl/i18n"
	"mmth-etl/types"
)

type LogType int

const (
	LogTypeNone LogType = iota
	LogTypeDiamond
	LogTypeCave
	LogTypeChallenge
	LogTypeRuneTicket
	LogTypeUpgradePanacea
	LogTypeGacha
	LogTypeOpen
	LogTypeSystemError
	LogTypeNameLabel
	LogTypeAutoRefreshDiamond
)

// IsNameLabelPrefix checks if body starts with any language's Name: prefix.
func IsNameLabelPrefix(body string) bool {
	for _, prefix := range i18n.GetAllNameLabels() {
		if len(body) >= len(prefix) && body[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}

// IdentifyLogType identifies a log type with a single pass over known patterns.
func IdentifyLogType(body string) LogType {
	isItemChangeLog := IsNameLabelPrefix(body)

	if isItemChangeLog {
		if types.DiamondRegex().MatchString(body) {
			return LogTypeDiamond
		}
		if types.RuneTicketRegex().MatchString(body) {
			return LogTypeRuneTicket
		}
		if types.UpgradePanaceaRegex().MatchString(body) {
			return LogTypeUpgradePanacea
		}
		return LogTypeNameLabel
	}

	if types.AutoRefreshDiamondRegex.MatchString(body) {
		return LogTypeAutoRefreshDiamond
	}

	if types.CaveEnterRegex().MatchString(body) || types.CaveFinishRegex().MatchString(body) || types.CaveErrorRegex.MatchString(body) {
		return LogTypeCave
	}

	if types.ChallengeQuestRegex().MatchString(body) || types.ChallengeTowerRegex().MatchString(body) {
		if types.ChallengeSuccessRegex().MatchString(body) || types.ChallengeFailedRegex().MatchString(body) {
			return LogTypeChallenge
		}
	}

	if types.SystemErrorRegex.MatchString(body) {
		return LogTypeSystemError
	}

	if types.GachaPrefixRegex().MatchString(body) {
		return LogTypeGacha
	}
	if types.OpenPrefixRegex().MatchString(body) {
		return LogTypeOpen
	}

	return LogTypeNone
}
