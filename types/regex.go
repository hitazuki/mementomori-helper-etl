package types

import "regexp"

var (
	// 物品变动日志正则（统一格式：Name: ItemName(Quality) × N）
	DiamondRegex        = regexp.MustCompile(`^Name: Diamonds\(None\) × (-?\d+)`)
	RuneTicketRegex     = regexp.MustCompile(`^Name: Rune Ticket\([A-Z]+\) × (-?\d+)`)
	UpgradePanaceaRegex = regexp.MustCompile(`^Name: Upgrade Panacea\([A-Z]+\) × (-?\d+)`)

	// 洞穴日志正则
	CaveEnterRegex  = regexp.MustCompile(`Enter Cave of Space-Time`)
	CaveFinishRegex = regexp.MustCompile(`Cave of Space-Time Finished`)
	CaveErrorRegex  = regexp.MustCompile(`KeyNotFoundException`)

	// 挑战日志正则
	ChallengeQuestRegex   = regexp.MustCompile(`^Challenge (\d+-\d+) boss`)
	ChallengeTowerRegex   = regexp.MustCompile(`^Challenge Tower of (Infinity|Azure|Crimson|Emerald|Amber) (\d+) layer`)
	ChallengeSuccessRegex = regexp.MustCompile(`triumphed`)
	ChallengeFailedRegex  = regexp.MustCompile(`failed`)
)
