package i18n

import (
	"fmt"
	"regexp"
)

// patternDefinition contains language-specific string values for pattern building.
type patternDefinition struct {
	// Name label (e.g., "Name" in English, "名称" in Chinese)
	NameLabel string

	// Item names
	Diamond        string
	RuneTicket     string
	UpgradePanacea string

	// Cave patterns
	CaveEnter  string
	CaveFinish string

	// Challenge patterns
	ChallengeKeyword string
	SuccessKeyword   string
	FailedKeyword    string

	// Tower names (for pattern building)
	TowerInfinity string
	TowerAzure    string
	TowerCrimson  string
	TowerEmerald  string
	TowerAmber    string

	// Source context patterns
	GachaPrefix string // Gacha prefix (e.g., "Gacha " in English)
	OpenPrefix  string // Open prefix (e.g., "Open " in English)
}

// languageDefinitions maps languages to their pattern definitions.
var languageDefinitions = map[Language]patternDefinition{
	LangEn: {
		NameLabel:        "Name",
		Diamond:          "Diamonds",
		RuneTicket:       "Rune Ticket",
		UpgradePanacea:   "Upgrade Panacea",
		CaveEnter:        "Enter Cave of Space-Time",
		CaveFinish:       "Cave of Space-Time Finished",
		ChallengeKeyword: "Challenge",
		SuccessKeyword:   "triumphed",
		FailedKeyword:    "failed",
		TowerInfinity:    "Tower of Infinity",
		TowerAzure:       "Tower of Azure",
		TowerCrimson:     "Tower of Crimson",
		TowerEmerald:     "Tower of Emerald",
		TowerAmber:       "Tower of Amber",
		GachaPrefix:      "Gacha ",
		OpenPrefix:       "Open ",
	},
	LangTw: {
		NameLabel:        "名称",
		Diamond:          "鑽石",
		RuneTicket:       "符石兌換券",
		UpgradePanacea:   "強化秘藥",
		CaveEnter:        "进入 時空洞窟",
		CaveFinish:       "時空洞窟已完成",
		ChallengeKeyword: "挑战",
		SuccessKeyword:   "勝利",
		FailedKeyword:    "敗北",
		TowerInfinity:    "無窮之塔",
		TowerAzure:       "憂藍之塔",
		TowerCrimson:     "業紅之塔",
		TowerEmerald:     "蒼翠之塔",
		TowerAmber:       "流金之塔",
		GachaPrefix:      "抽卡 ",
		OpenPrefix:       "開啟 ",
	},
	LangJa: {
		NameLabel:        "名前",
		Diamond:          "ダイヤ",
		RuneTicket:       "ルーンチケット",
		UpgradePanacea:   "強化秘薬",
		CaveEnter:        "入る 時空の洞窟",
		CaveFinish:       "時空の洞窟 完了済み",
		ChallengeKeyword: "挑戦",
		SuccessKeyword:   "勝利",
		FailedKeyword:    "敗北",
		TowerInfinity:    "無窮の塔",
		TowerAzure:       "藍の塔",
		TowerCrimson:     "紅の塔",
		TowerEmerald:     "翠の塔",
		TowerAmber:       "黄の塔",
		GachaPrefix:      "ガチャ ",
		OpenPrefix:       "開く ",
	},
	LangKo: {
		NameLabel:        "이름",
		Diamond:          "다이아",
		RuneTicket:       "룬 티켓",
		UpgradePanacea:   "강화의 비약",
		CaveEnter:        "입장 시공의 동굴",
		CaveFinish:       "시공의 동굴 완료됨",
		ChallengeKeyword: "도전",
		SuccessKeyword:   "승리",
		FailedKeyword:    "패배",
		TowerInfinity:    "무한의 탑",
		TowerAzure:       "남청의 탑",
		TowerCrimson:     "홍염의 탑",
		TowerEmerald:     "비취의 탑",
		TowerAmber:       "황철의 탑",
		GachaPrefix:      "가챠 ",
		OpenPrefix:       "열기 ",
	},
}

// buildPatternSet creates a PatternSet for the given language.
func buildPatternSet(lang Language) *PatternSet {
	def := languageDefinitions[lang]

	// Build tower name pattern (union of all tower names)
	towerPattern := fmt.Sprintf("%s|%s|%s|%s|%s",
		regexp.QuoteMeta(def.TowerInfinity),
		regexp.QuoteMeta(def.TowerAzure),
		regexp.QuoteMeta(def.TowerCrimson),
		regexp.QuoteMeta(def.TowerEmerald),
		regexp.QuoteMeta(def.TowerAmber),
	)

	// Challenge patterns differ by language:
	// EN/TW: "Challenge/Tower tower 800 layer" (starts with keyword)
	// JA/KO: "塔名 800 層に挑戦/층에 도전" (starts with tower name)
	var challengeTower *regexp.Regexp
	switch lang {
	case LangJa, LangKo:
		// JA: 無窮の塔 800 層に挑戦 1 回
		// KO: 무한의 탑 800 층에 도전 1회
		challengeTower = regexp.MustCompile(fmt.Sprintf(
			`^(%s) (\d+)`,
			towerPattern,
		))
	default:
		// EN: Challenge Tower of Infinity 800 layer
		// TW: 挑战 無窮之塔 800 层
		challengeTower = regexp.MustCompile(fmt.Sprintf(
			`^%s (%s) (\d+)`,
			regexp.QuoteMeta(def.ChallengeKeyword),
			towerPattern,
		))
	}

	// Quest challenge patterns differ by language:
	// EN/TW: "Challenge 36-13 boss one time"
	// JA: "36-13 ボスに挑戦 1 回"
	// KO: "36-13 보스에 도전 1회"
	var challengeQuest *regexp.Regexp
	switch lang {
	case LangJa:
		challengeQuest = regexp.MustCompile(`^(\d+-\d+) ボスに挑戦`)
	case LangKo:
		challengeQuest = regexp.MustCompile(`^(\d+-\d+) 보스에 도전`)
	default:
		challengeQuest = regexp.MustCompile(fmt.Sprintf(
			`^%s (\d+-\d+) boss`,
			regexp.QuoteMeta(def.ChallengeKeyword),
		))
	}

	return &PatternSet{
		// Item change patterns: Name: ItemName(Quality) × N
		Diamond: regexp.MustCompile(fmt.Sprintf(
			`^%s: %s\(None\) × (-?\d+)`,
			regexp.QuoteMeta(def.NameLabel),
			regexp.QuoteMeta(def.Diamond),
		)),
		RuneTicket: regexp.MustCompile(fmt.Sprintf(
			`^%s: %s\([A-Z]+\) × (-?\d+)`,
			regexp.QuoteMeta(def.NameLabel),
			regexp.QuoteMeta(def.RuneTicket),
		)),
		UpgradePanacea: regexp.MustCompile(fmt.Sprintf(
			`^%s: %s\([A-Z]+\) × (-?\d+)`,
			regexp.QuoteMeta(def.NameLabel),
			regexp.QuoteMeta(def.UpgradePanacea),
		)),

		// Cave patterns
		CaveEnter:  regexp.MustCompile(regexp.QuoteMeta(def.CaveEnter)),
		CaveFinish: regexp.MustCompile(regexp.QuoteMeta(def.CaveFinish)),

		// Challenge patterns
		ChallengeQuest:   challengeQuest,
		ChallengeTower:   challengeTower,
		ChallengeSuccess: regexp.MustCompile(regexp.QuoteMeta(def.SuccessKeyword)),
		ChallengeFailed:  regexp.MustCompile(regexp.QuoteMeta(def.FailedKeyword)),

		// Source context patterns
		GachaPrefix: regexp.MustCompile(`^` + regexp.QuoteMeta(def.GachaPrefix)),
		OpenPrefix:  regexp.MustCompile(`^` + regexp.QuoteMeta(def.OpenPrefix)),
	}
}

// TowerNameToType normalizes a language-specific tower name to a canonical tower type.
// Returns empty string if the name is not recognized.
func (m *Manager) TowerNameToType(name string) string {
	// Check all languages for the tower name
	for _, lang := range []Language{LangEn, LangTw, LangJa, LangKo} {
		def := languageDefinitions[lang]
		switch name {
		case def.TowerInfinity:
			return "Infinity"
		case def.TowerAzure:
			return "Azure"
		case def.TowerCrimson:
			return "Crimson"
		case def.TowerEmerald:
			return "Emerald"
		case def.TowerAmber:
			return "Amber"
		}
	}
	return ""
}
