package parser

import (
	"mmth-etl/i18n"
	"mmth-etl/types"
	"testing"
)

func TestGetSourceID_English(t *testing.T) {
	mgr := i18n.NewManager()
	mgr.SetLanguage(i18n.LangEn)
	types.InitI18n(mgr)
	InvalidateSourceCache()

	tests := []struct {
		source     string
		expectedID i18n.SourceID
		desc       string
	}{
		{"Fountain of Prayers:", i18n.SourceIDFountainOfPrayers, "Fountain of Prayers"},
		{"Presents Box Claim All", i18n.SourceIDPresentsBox, "Presents Box"},
		{"Monthly Boost Already Claimed", i18n.SourceIDMonthlyBoost, "Monthly Boost"},
		{"Login", i18n.SourceIDLoginBonus, "Login"},
		{"Auto Buy Store Items", i18n.SourceIDAutoBuyStore, "Auto Buy Store Items"},
		{"Current expected value: 65. Today, 0/9 auto-refreshes completed. Refreshing now.", i18n.SourceIDFountainOfPrayers, "Auto Refresh"},
		{"You have no more challenges left.", i18n.SourceIDMissionsClaimed, "Missions Claim All"},
		{"Cave of Space-TimeFinished", i18n.SourceIDMissionsClaimed, "Cave Finished"},
		{"Tower of Infinity:", i18n.SourceIDTowerInfinity, "Tower of Infinity"},
		{"You have triumphed.", i18n.SourceIDTempleIllusions, "Temple of Illusions"},
		{"Get Daily 's 60 Reward", i18n.SourceIDDailyMissionReward, "Get Daily Reward"},
		{"Get Weekly 's 80 Reward", i18n.SourceIDWeeklyMissionReward, "Get Weekly Reward"},
		{"Get Daily 's 100 Reward", i18n.RewardMissionSourceID(i18n.MissionGroupDailyID, 100), "Get Daily 100 Reward"},
		{"Get Guild 's 2000 Reward", i18n.SourceIDGuildMissionReward, "Get Guild Reward"},
		{"Get Guild 's 3000 Reward", i18n.RewardMissionSourceID(i18n.SourceIDGuild, 3000), "Get Guild 3000 Reward"},
		{"Total Logins This Month: 15/30", i18n.SourceIDTotalLogins, "Total Logins"},
		{"A player in your World clears Floor {0}", i18n.SourceIDWorldClears, "World clears"},
		{"Unknown source", 0, "Unknown source"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			result := GetSourceID(tt.source)
			if result != tt.expectedID {
				t.Errorf("GetSourceID(%q) = %d, want %d", tt.source, result, tt.expectedID)
			}
		})
	}
}

func TestGetSourceID_TraditionalChinese(t *testing.T) {
	mgr := i18n.NewManager()
	mgr.SetLanguage(i18n.LangTw)
	types.InitI18n(mgr)
	InvalidateSourceCache()

	tests := []struct {
		source     string
		expectedID i18n.SourceID
		desc       string
	}{
		{"祈願之泉:", i18n.SourceIDFountainOfPrayers, "祈願之泉"},
		{"禮物箱", i18n.SourceIDPresentsBox, "Presents Box TW"},
		{"每月強化組合包", i18n.SourceIDMonthlyBoost, "Monthly Boost TW"},
		{"簽到獎勵:", i18n.SourceIDLoginBonus, "Login Bonus TW"},
		{"自动购买商城物品", i18n.SourceIDAutoBuyStore, "Auto Buy TW"},
		{"当前期望值：65.5，今日已自动刷新 1/9 次，现在刷新。", i18n.SourceIDFountainOfPrayers, "Auto Refresh ZH"},
		{"剩餘挑戰次數不足", i18n.SourceIDMissionsClaimed, "Missions Claim All TW"},
		{"時空洞窟已完成", i18n.SourceIDMissionsClaimed, "Cave Finished TW"},
		{"無窮之塔:", i18n.SourceIDTowerInfinity, "Tower of Infinity TW"},
		{"勝利", i18n.SourceIDTempleIllusions, "Temple of Illusions TW"},
		{"领取 Daily 的 60 奖励", i18n.SourceIDDailyMissionReward, "领取 Daily 奖励"},
		{"领取 Weekly 的 80 奖励", i18n.SourceIDWeeklyMissionReward, "领取 Weekly 奖励"},
		{"领取 Weekly 的 120 奖励", i18n.RewardMissionSourceID(i18n.MissionGroupWeeklyID, 120), "领取 Weekly 120 奖励"},
		{"领取 Guild 的 2000 奖励", i18n.SourceIDGuildMissionReward, "领取 Guild 奖励"},
		{"领取 Guild 的 3000 奖励", i18n.RewardMissionSourceID(i18n.SourceIDGuild, 3000), "领取 Guild 3000 奖励"},
		{"本月累計簽到天數：15/30", i18n.SourceIDTotalLogins, "Total Logins TW"},
		{"本世界首次有玩家通關", i18n.SourceIDWorldClears, "World clears TW"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			result := GetSourceID(tt.source)
			if result != tt.expectedID {
				t.Errorf("GetSourceID(%q) = %d, want %d", tt.source, result, tt.expectedID)
			}
		})
	}
}

func TestGetSourceID_Japanese(t *testing.T) {
	mgr := i18n.NewManager()
	mgr.SetLanguage(i18n.LangJa)
	types.InitI18n(mgr)
	InvalidateSourceCache()

	tests := []struct {
		source     string
		expectedID i18n.SourceID
		desc       string
	}{
		{"祈りの泉:", i18n.SourceIDFountainOfPrayers, "Fountain of Prayers JA"},
		{"プレゼントボックス", i18n.SourceIDPresentsBox, "Presents Box JA"},
		{"ログイン", i18n.SourceIDLoginBonus, "Login Bonus JA"},
		{"自動購入ストアアイテム", i18n.SourceIDAutoBuyStore, "Auto Buy JA"},
		{"現在の期待値：65.5。本日自動リフレッシュ 1/9 回完了。リフレッシュ中。", i18n.SourceIDFountainOfPrayers, "Auto Refresh JA"},
		{"残り挑戦回数がありません", i18n.SourceIDMissionsClaimed, "Missions Claimed JA"},
		{"時空の洞窟完了", i18n.SourceIDMissionsClaimed, "Cave Finished JA"},
		{"無窮の塔:", i18n.SourceIDTowerInfinity, "Tower of Infinity JA"},
		{"勝利しました", i18n.SourceIDTempleIllusions, "Temple of Illusions JA"},
		{"今月の合計ログイン日数：", i18n.SourceIDTotalLogins, "Total Logins JA"},
		{"ワールド内のプレイヤーが初めて", i18n.SourceIDWorldClears, "World clears JA"},
		{"Daily の 60 報酬", i18n.SourceIDDailyMissionReward, "Daily Mission Reward JA"},
		{"Weekly の 80 報酬", i18n.SourceIDWeeklyMissionReward, "Weekly Mission Reward JA"},
		{"Daily の 60 の報酬を受け取る", i18n.SourceIDDailyMissionReward, "Daily Mission Reward helper JA"},
		{"Main の 100 報酬", i18n.MissionGroupMainID, "Main Mission Reward JA"},
		{"Guild の 2000 の報酬を受け取る", i18n.SourceIDGuildMissionReward, "Guild Mission Reward JA"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			result := GetSourceID(tt.source)
			if result != tt.expectedID {
				t.Errorf("GetSourceID(%q) = %d, want %d", tt.source, result, tt.expectedID)
			}
		})
	}
}

func TestGetSourceID_Korean(t *testing.T) {
	mgr := i18n.NewManager()
	mgr.SetLanguage(i18n.LangKo)
	types.InitI18n(mgr)
	InvalidateSourceCache()

	tests := []struct {
		source     string
		expectedID i18n.SourceID
		desc       string
	}{
		{"기원의 샘:", i18n.SourceIDFountainOfPrayers, "Fountain of Prayers KO"},
		{"선물 상자", i18n.SourceIDPresentsBox, "Presents Box KO"},
		{"월간 부스트", i18n.SourceIDMonthlyBoost, "Monthly Boost KO"},
		{"이번 달 보상 수령:", i18n.SourceIDTotalLogins, "Total Logins KO"},
		{"월드 내 플레이어가 최초로", i18n.SourceIDWorldClears, "World clears KO"},
		{"로그인", i18n.SourceIDLoginBonus, "Login Bonus KO"},
		{"자동으로 상점 아이템 구매", i18n.SourceIDAutoBuyStore, "Auto Buy KO"},
		{"현재 기대값: 65.5. 오늘 1/9회 자동 새로고침 완료. 지금 새로고침합니다.", i18n.SourceIDFountainOfPrayers, "Auto Refresh KO"},
		{"현재 작업의 다이아몬드 예상 값이 20 미만이므로", i18n.SourceIDMissionsClaimed, "Expected Value KO"},
		{"시공의 동굴 완료", i18n.SourceIDMissionsClaimed, "Cave Finished KO"},
		{"무한의 탑:", i18n.SourceIDTowerInfinity, "Tower of Infinity KO"},
		{"승리했습니다.", i18n.SourceIDTempleIllusions, "Temple of Illusions KO"},
		{"일일 의 60 보상", i18n.SourceIDDailyMissionReward, "Daily Mission Reward KO"},
		{"주간 의 80 보상", i18n.SourceIDWeeklyMissionReward, "Weekly Mission Reward KO"},
		{"일일의 60 보상을 수령합니다", i18n.SourceIDDailyMissionReward, "Daily Mission Reward helper KO"},
		{"Weekly의 80 보상을 수령합니다", i18n.SourceIDWeeklyMissionReward, "Weekly Mission Reward helper KO"},
		{"메인 의 100 보상", i18n.MissionGroupMainID, "Main Mission Reward KO"},
		{"Guild의 2000 보상을 수령합니다", i18n.SourceIDGuildMissionReward, "Guild Mission Reward KO"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			result := GetSourceID(tt.source)
			if result != tt.expectedID {
				t.Errorf("GetSourceID(%q) = %d, want %d", tt.source, result, tt.expectedID)
			}
		})
	}
}

func TestMapSourceWithID(t *testing.T) {
	mgr := i18n.NewManager()
	mgr.SetLanguage(i18n.LangEn)
	types.InitI18n(mgr)
	InvalidateSourceCache()

	tests := []struct {
		source     string
		expected   string
		expectedID i18n.SourceID
		desc       string
	}{
		{"Gacha 黒葬武具ガチャ 10 times", "Gacha 黒葬武具ガチャ 10 times", i18n.SourceIDGacha, "Gacha EN"},
		{"Open Gold Sealed Chest x 5", "Open Gold Sealed Chest x 5", i18n.SourceIDOpen, "Open EN"},
		{"Tower of Infinity: 800", "Tower of Infinity", i18n.SourceIDTowerInfinity, "Tower of Infinity EN"},
		{"You have triumphed.", "Temple of Illusions", i18n.SourceIDTempleIllusions, "Temple of Illusions EN"},
		{"Fountain of Prayers:", "Fountain of Prayers", i18n.SourceIDFountainOfPrayers, "Fountain of Prayers"},
		{"Login", "Login Bonus", i18n.SourceIDLoginBonus, "Login"},
		{"Get Daily 's 60 Reward", "id:23214000060", i18n.SourceIDDailyMissionReward, "Daily reward composite ID"},
		{"Get Guild 's 2000 Reward", "id:111002000", i18n.SourceIDGuildMissionReward, "Guild reward composite ID"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			alias, sourceID := MapSourceWithID(tt.source)
			if alias != tt.expected {
				t.Errorf("MapSourceWithID(%q) alias = %q, want %q", tt.source, alias, tt.expected)
			}
			if sourceID != tt.expectedID {
				t.Errorf("MapSourceWithID(%q) sourceID = %d, want %d", tt.source, sourceID, tt.expectedID)
			}
		})
	}
}

func TestMapSourceWithID_Chinese(t *testing.T) {
	mgr := i18n.NewManager()
	mgr.SetLanguage(i18n.LangTw)
	types.InitI18n(mgr)
	InvalidateSourceCache()

	tests := []struct {
		source     string
		expected   string
		expectedID i18n.SourceID
		desc       string
	}{
		{"抽卡 黒葬武具ガチャ 5 次, 消耗 金幣×250000", "抽卡 黒葬武具ガチャ 5 次, 消耗 金幣×250000", i18n.SourceIDGacha, "Gacha TW"},
		{"開啟 上級封印寶箱 x 5", "開啟 上級封印寶箱 x 5", i18n.SourceIDOpen, "Open TW"},
		{"無窮之塔: 800", "Tower of Infinity", i18n.SourceIDTowerInfinity, "Tower of Infinity TW"},
		{"勝利", "Temple of Illusions", i18n.SourceIDTempleIllusions, "Temple of Illusions TW"},
		{"领取 Weekly 的 80 奖励", "id:23215000080", i18n.SourceIDWeeklyMissionReward, "Weekly reward composite ID"},
		{"领取 Guild 的 2000 奖励", "id:111002000", i18n.SourceIDGuildMissionReward, "Guild reward composite ID"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			alias, sourceID := MapSourceWithID(tt.source)
			if alias != tt.expected {
				t.Errorf("MapSourceWithID(%q) alias = %q, want %q", tt.source, alias, tt.expected)
			}
			if sourceID != tt.expectedID {
				t.Errorf("MapSourceWithID(%q) sourceID = %d, want %d", tt.source, sourceID, tt.expectedID)
			}
		})
	}
}

func TestExtractChangeRecord_CompositeSourceUsesIDKey(t *testing.T) {
	mgr := i18n.NewManager()
	mgr.SetLanguage(i18n.LangTw)
	types.InitI18n(mgr)
	InvalidateSourceCache()

	record := ExtractChangeRecord(ParsedLog{
		Character: "test",
		Timestamp: "2026-04-30T12:00:00+08:00",
		Body:      "名称: 鑽石(None) × 300",
	}, "领取 Daily 的 60 奖励", LogTypeDiamond)

	if record == nil {
		t.Fatal("ExtractChangeRecord returned nil")
	}
	if record.Source != "id:23214000060" {
		t.Errorf("record.Source = %q, want %q", record.Source, "id:23214000060")
	}
	if record.SourceID != int(i18n.SourceIDDailyMissionReward) {
		t.Errorf("record.SourceID = %d, want %d", record.SourceID, i18n.SourceIDDailyMissionReward)
	}
}

func TestExtractChangeRecord_AutoRefreshDiamond(t *testing.T) {
	record := ExtractChangeRecord(ParsedLog{
		Character: "test",
		Timestamp: "2026-04-30T12:00:00+08:00",
		Body:      "Current expected value: 65.5. Today, 1/9 auto-refreshes completed. Refreshing now.",
	}, "", LogTypeAutoRefreshDiamond)

	if record == nil {
		t.Fatal("ExtractChangeRecord returned nil")
	}
	if record.Amount != -20 {
		t.Errorf("record.Amount = %d, want -20", record.Amount)
	}
	if record.Source != "Fountain of Prayers" {
		t.Errorf("record.Source = %q, want %q", record.Source, "Fountain of Prayers")
	}
	if record.SourceID != int(i18n.SourceIDFountainOfPrayers) {
		t.Errorf("record.SourceID = %d, want %d", record.SourceID, i18n.SourceIDFountainOfPrayers)
	}
}

func TestCleanSourceSuffix(t *testing.T) {
	tests := []struct {
		body     string
		logType  LogType
		expected string
		desc     string
	}{
		// Gacha
		{"Gacha test pool 10 times", LogTypeGacha, "Gacha test pool", "EN Gacha with count"},
		{"Gacha pool", LogTypeGacha, "Gacha pool", "EN Gacha no count"},
		{"抽卡 測試卡池 10 次", LogTypeGacha, "抽卡 測試卡池", "TW Gacha with count"},
		{"ガチャ テスト 5 回", LogTypeGacha, "ガチャ テスト", "JA Gacha with count"},

		// Open
		{"Open Gold Chest x 5", LogTypeOpen, "Open Gold Chest", "EN Open with count"},
		{"Open Box", LogTypeOpen, "Open Box", "EN Open no count"},
		{"開啟 上級封印寶箱 x 5", LogTypeOpen, "開啟 上級封印寶箱", "TW Open with count"},
		{"開く 宝箱 x 10", LogTypeOpen, "開く 宝箱", "JA Open with count"},

		// Other types - no change
		{"Some other log", LogTypeNone, "Some other log", "None type unchanged"},
		{"Login bonus", LogTypeDiamond, "Login bonus", "Diamond type unchanged"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			result := CleanSourceSuffix(tt.body, tt.logType)
			if result != tt.expected {
				t.Errorf("CleanSourceSuffix(%q, %v) = %q, want %q", tt.body, tt.logType, result, tt.expected)
			}
		})
	}
}
