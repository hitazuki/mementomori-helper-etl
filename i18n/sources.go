package i18n

// SourceID represents a unique identifier for a diamond source.
// - 0: Unknown/unmatched source
// - 1-99999: Game TextResource IDs
// - 100000+: Helper custom IDs
type SourceID int

// SourceEntry defines a single source mapping for a specific language.
type SourceEntry struct {
	ID    SourceID
	Alias string
	Text  string // The source text in the specific language
}

// Game built-in source IDs from TextResource
const (
	SourceIDFountainOfPrayers SourceID = 140   // Fountain of Prayers
	SourceIDOpen              SourceID = 67    // Open (開啟)
	SourceIDLoginBonus        SourceID = 719   // Login Bonus (签到奖励)
	SourceIDPresentsBox       SourceID = 21308 // Presents Box
	SourceIDMonthlyBoost      SourceID = 21332 // Monthly Boost
	SourceIDTotalLogins       SourceID = 3331  // Total Logins
	SourceIDWorldClears       SourceID = 23277 // World Player Clears
)

// Mission group IDs from TextResource
const (
	MissionGroupDailyID  SourceID = 23214
	MissionGroupWeeklyID SourceID = 23215
	MissionGroupMainID   SourceID = 23213
)

// Helper custom source IDs
const (
	SourceIDAutoBuyStore    SourceID = 100002
	SourceIDExpectedValue   SourceID = 100003
	SourceIDMissionsClaimed SourceID = 100004
	SourceIDGacha           SourceID = 100005
	SourceIDTowerInfinity   SourceID = 100007
	SourceIDTempleIllusions SourceID = 100008
)

// SourceDefinitions maps languages to their source entries.
// This is used to build language-specific source tables.
var SourceDefinitions = map[Language][]SourceEntry{
	LangEn: {
		{SourceIDFountainOfPrayers, "Fountain of Prayers", "Fountain of Prayers:"},
		{SourceIDPresentsBox, "Presents Box", "Presents Box Claim All"},
		{SourceIDMonthlyBoost, "Monthly Boost", "Monthly Boost Already Claimed"},
		{SourceIDTotalLogins, "Total Logins This Month", "Total Logins This Month:"},
		{SourceIDWorldClears, "World Player Clears", "A player in your World clears"},
		{SourceIDWorldClears, "World Player Clears", "A player in your World reaches"},
		{SourceIDLoginBonus, "Login Bonus", "Login"},
		{SourceIDAutoBuyStore, "Auto Buy Store Items", "Auto Buy Store Items"},
		{SourceIDExpectedValue, "Expected Value Below 20", "The expected diamond value"},
		{SourceIDMissionsClaimed, "Missions Claim All", "You have no more challenges left."},
		{SourceIDMissionsClaimed, "Missions Claim All", "Cave of Space-TimeFinished"},
		{SourceIDMissionsClaimed, "Missions Claim All", "Nothing to receive"},
		{SourceIDTowerInfinity, "Tower of Infinity", "Tower of Infinity:"},
		{SourceIDTempleIllusions, "Temple of Illusions", "You have triumphed."},
	},
	LangTw: {
		{SourceIDFountainOfPrayers, "Fountain of Prayers", "祈願之泉:"},
		{SourceIDPresentsBox, "Presents Box", "禮物箱"},
		{SourceIDMonthlyBoost, "Monthly Boost", "每月強化組合包"},
		{SourceIDTotalLogins, "Total Logins This Month", "本月累計簽到天數："},
		{SourceIDWorldClears, "World Player Clears", "本世界首次有玩家"},
		{SourceIDLoginBonus, "Login Bonus", "簽到獎勵:"},
		{SourceIDAutoBuyStore, "Auto Buy Store Items", "自動購買商城物品"},
		{SourceIDExpectedValue, "Expected Value Below 20", "當前任務的鑽石數量期望值"},
		{SourceIDMissionsClaimed, "Missions Claim All", "剩餘挑戰次數不足"},
		{SourceIDMissionsClaimed, "Missions Claim All", "没有可以领取的"},
		{SourceIDMissionsClaimed, "Missions Claim All", "時空洞窟已完成"},
		{SourceIDTowerInfinity, "Tower of Infinity", "無窮之塔:"},
		{SourceIDTempleIllusions, "Temple of Illusions", "勝利"},
	},
	LangJa: {
		{SourceIDFountainOfPrayers, "Fountain of Prayers", "祈りの泉:"},
		{SourceIDPresentsBox, "Presents Box", "プレゼントボックス"},
		{SourceIDMonthlyBoost, "Monthly Boost", "月間ブースト"},
		{SourceIDTotalLogins, "Total Logins This Month", "今月の合計ログイン日数："},
		{SourceIDWorldClears, "World Player Clears", "ワールド内のプレイヤーが初めて"},
		{SourceIDLoginBonus, "Login Bonus", "ログイン"},
		{SourceIDAutoBuyStore, "Auto Buy Store Items", "自動購入ストアアイテム"},
		{SourceIDExpectedValue, "Expected Value Below 20", "現在のタスクのダイヤの期待値"},
		{SourceIDMissionsClaimed, "Missions Claim All", "残り挑戦回数がありません"},
		{SourceIDMissionsClaimed, "Missions Claim All", "受け取れるものはありません"},
		{SourceIDMissionsClaimed, "Missions Claim All", "時空の洞窟完了"},
		{SourceIDTowerInfinity, "Tower of Infinity", "無窮の塔:"},
		{SourceIDTempleIllusions, "Temple of Illusions", "勝利しました"},
	},
	LangKo: {
		{SourceIDFountainOfPrayers, "Fountain of Prayers", "기원의 샘:"},
		{SourceIDPresentsBox, "Presents Box", "선물 상자"},
		{SourceIDMonthlyBoost, "Monthly Boost", "월간 부스트"},
		{SourceIDTotalLogins, "Total Logins This Month", "이번 달 보상 수령:"},
		{SourceIDWorldClears, "World Player Clears", "월드 내 플레이어가 최초로"},
		{SourceIDLoginBonus, "Login Bonus", "로그인"},
		{SourceIDAutoBuyStore, "Auto Buy Store Items", "자동으로 상점 아이템 구매"},
		{SourceIDExpectedValue, "Expected Value Below 20", "현재 작업의 다이아몬드 예상 값"},
		{SourceIDMissionsClaimed, "Missions Claim All", "시공의 동굴 완료"},
		{SourceIDMissionsClaimed, "Missions Claim All", "수령 가능한 것이 없습니다"},
		{SourceIDTowerInfinity, "Tower of Infinity", "무한의 탑:"},
		{SourceIDTempleIllusions, "Temple of Illusions", "승리했습니다."},
	},
}

// RewardMissionPattern defines a prefix pattern for reward mission matching.
type RewardMissionPattern struct {
	Prefix   string
	SourceID SourceID
	Alias    string
}

// rewardMissionDefinitions maps languages to their reward mission patterns.
var rewardMissionDefinitions = map[Language][]RewardMissionPattern{
	LangEn: {
		{"Get Daily ", MissionGroupDailyID, "Daily Mission Reward"},
		{"Get Weekly ", MissionGroupWeeklyID, "Weekly Mission Reward"},
		{"Get Main ", MissionGroupMainID, "Main Mission Reward"},
	},
	LangTw: {
		{"领取 Daily ", MissionGroupDailyID, "Daily Mission Reward"},
		{"领取 Weekly ", MissionGroupWeeklyID, "Weekly Mission Reward"},
		{"领取 Main ", MissionGroupMainID, "Main Mission Reward"},
	},
	LangJa: {
		{"Daily の ", MissionGroupDailyID, "Daily Mission Reward"},
		{"Weekly の ", MissionGroupWeeklyID, "Weekly Mission Reward"},
		{"Main の ", MissionGroupMainID, "Main Mission Reward"},
	},
	LangKo: {
		{"일일 의 ", MissionGroupDailyID, "Daily Mission Reward"},
		{"주간 의 ", MissionGroupWeeklyID, "Weekly Mission Reward"},
		{"메인 의 ", MissionGroupMainID, "Main Mission Reward"},
	},
}

// GetSources returns the source entries for the given language.
func (m *Manager) GetSources(lang Language) []SourceEntry {
	if sources, ok := SourceDefinitions[lang]; ok {
		return sources
	}
	return SourceDefinitions[LangEn]
}

// CurrentSources returns the source entries for the current language.
func (m *Manager) CurrentSources() []SourceEntry {
	return m.GetSources(m.currentLang)
}

// GetRewardMissionPatterns returns the reward mission patterns for the given language.
func (m *Manager) GetRewardMissionPatterns(lang Language) []RewardMissionPattern {
	if patterns, ok := rewardMissionDefinitions[lang]; ok {
		return patterns
	}
	return rewardMissionDefinitions[LangEn]
}

// CurrentRewardMissionPatterns returns the patterns for the current language.
func (m *Manager) CurrentRewardMissionPatterns() []RewardMissionPattern {
	return m.GetRewardMissionPatterns(m.currentLang)
}
