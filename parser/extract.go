package parser

import (
	"mmth-etl/types"
	"regexp"
	"strconv"
)

// ExtractChangeRecord 提取物品变动记录（钻石/饼干/红水通用）
func ExtractChangeRecord(parsed ParsedLog, source string, logType LogType) *types.ChangeRecord {
	var regex *regexp.Regexp

	switch logType {
	case LogTypeDiamond:
		regex = types.DiamondRegex
	case LogTypeRuneTicket:
		regex = types.RuneTicketRegex
	case LogTypeUpgradePanacea:
		regex = types.UpgradePanaceaRegex
	default:
		return nil
	}

	amount := extractAmount(parsed.Body, regex)
	if amount == 0 {
		return nil
	}

	// 统一应用来源映射
	mappedSource := MapSourceToAlias(source)
	if mappedSource == "" {
		mappedSource = "none"
	}

	return &types.ChangeRecord{
		Character:   parsed.Character,
		Timestamp:   parsed.Timestamp,
		Amount:      amount,
		Source:      mappedSource,
	}
}

// extractAmount 提取物品变动数量（统一方法）
func extractAmount(body string, regex *regexp.Regexp) int {
	matches := regex.FindStringSubmatch(body)
	if len(matches) < 2 {
		return 0
	}
	amount, _ := strconv.Atoi(matches[1])
	return amount
}

// ExtractCaveRecord 提取洞穴记录
func ExtractCaveRecord(parsed ParsedLog) *types.CaveRecord {
	body := parsed.Body
	date := parsed.Timestamp[:10]

	if types.CaveErrorRegex.MatchString(body) {
		return &types.CaveRecord{
			Character:   parsed.Character,
			Timestamp:   parsed.Timestamp,
			Status:      types.CaveStatusError,
			Date:        date,
		}
	}

	if types.CaveFinishRegex.MatchString(body) {
		return &types.CaveRecord{
			Character:   parsed.Character,
			Timestamp:   parsed.Timestamp,
			Status:      types.CaveStatusFinished,
			Date:        date,
		}
	}

	if types.CaveEnterRegex.MatchString(body) {
		return &types.CaveRecord{
			Character:   parsed.Character,
			Timestamp:   parsed.Timestamp,
			Status:      types.CaveStatusStarted,
			Date:        date,
		}
	}

	return nil
}

// ExtractChallengeRecord 提取挑战记录
func ExtractChallengeRecord(parsed ParsedLog) *types.ChallengeRecord {
	body := parsed.Body

	record := &types.ChallengeRecord{
		Character:   parsed.Character,
		Timestamp:   parsed.Timestamp,
	}

	if types.ChallengeSuccessRegex.MatchString(body) {
		record.Status = types.ChallengeStatusSuccess
	} else {
		record.Status = types.ChallengeStatusFailed
	}

	// 匹配主线挑战
	if matches := types.ChallengeQuestRegex.FindStringSubmatch(body); len(matches) > 1 {
		record.Type = types.ChallengeTypeQuest
		record.Level = matches[1]
		return record
	}

	// 匹配塔挑战
	if matches := types.ChallengeTowerRegex.FindStringSubmatch(body); len(matches) > 2 {
		record.Type = types.ChallengeTypeTower
		record.TowerType = types.TowerType(matches[1])
		record.Level = matches[2]
		return record
	}

	return nil
}
