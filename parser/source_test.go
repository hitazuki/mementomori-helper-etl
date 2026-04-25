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
		{"You have no more challenges left.", i18n.SourceIDMissionsClaimed, "Missions Claim All"},
		{"Cave of Space-TimeFinished", i18n.SourceIDMissionsClaimed, "Cave Finished"},
		{"Tower of Infinity:", i18n.SourceIDTowerInfinity, "Tower of Infinity"},
		{"You have triumphed.", i18n.SourceIDTempleIllusions, "Temple of Illusions"},
		{"Get Daily 's 60 Reward", i18n.MissionGroupDailyID, "Get Daily Reward"},
		{"Get Weekly 's 80 Reward", i18n.MissionGroupWeeklyID, "Get Weekly Reward"},
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
		{"剩餘挑戰次數不足", i18n.SourceIDMissionsClaimed, "Missions Claim All TW"},
		{"時空洞窟已完成", i18n.SourceIDMissionsClaimed, "Cave Finished TW"},
		{"無窮之塔:", i18n.SourceIDTowerInfinity, "Tower of Infinity TW"},
		{"勝利", i18n.SourceIDTempleIllusions, "Temple of Illusions TW"},
		{"领取 Daily 的 60 奖励", i18n.MissionGroupDailyID, "领取 Daily 奖励"},
		{"领取 Weekly 的 80 奖励", i18n.MissionGroupWeeklyID, "领取 Weekly 奖励"},
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
		{"残り挑戦回数がありません", i18n.SourceIDMissionsClaimed, "Missions Claimed JA"},
		{"時空の洞窟完了", i18n.SourceIDMissionsClaimed, "Cave Finished JA"},
		{"無窮の塔:", i18n.SourceIDTowerInfinity, "Tower of Infinity JA"},
		{"勝利しました", i18n.SourceIDTempleIllusions, "Temple of Illusions JA"},
		{"今月の合計ログイン日数：", i18n.SourceIDTotalLogins, "Total Logins JA"},
		{"ワールド内のプレイヤーが初めて", i18n.SourceIDWorldClears, "World clears JA"},
		{"Daily の 60 報酬", i18n.MissionGroupDailyID, "Daily Mission Reward JA"},
		{"Weekly の 80 報酬", i18n.MissionGroupWeeklyID, "Weekly Mission Reward JA"},
		{"Main の 100 報酬", i18n.MissionGroupMainID, "Main Mission Reward JA"},
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
		{"현재 작업의 다이아몬드 예상 값이 20 미만이므로", i18n.SourceIDMissionsClaimed, "Expected Value KO"},
		{"시공의 동굴 완료", i18n.SourceIDMissionsClaimed, "Cave Finished KO"},
		{"무한의 탑:", i18n.SourceIDTowerInfinity, "Tower of Infinity KO"},
		{"승리했습니다.", i18n.SourceIDTempleIllusions, "Temple of Illusions KO"},
		{"일일 의 60 보상", i18n.MissionGroupDailyID, "Daily Mission Reward KO"},
		{"주간 의 80 보상", i18n.MissionGroupWeeklyID, "Weekly Mission Reward KO"},
		{"메인 의 100 보상", i18n.MissionGroupMainID, "Main Mission Reward KO"},
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
		expectedID i18n.SourceID
		desc       string
	}{
		{"Gacha 黒葬武具ガチャ 10 times", i18n.SourceIDGacha, "Gacha EN"},
		{"Open Gold Sealed Chest x 5", i18n.SourceIDOpen, "Open EN"},
		{"Tower of Infinity: 800", i18n.SourceIDTowerInfinity, "Tower of Infinity EN"},
		{"You have triumphed.", i18n.SourceIDTempleIllusions, "Temple of Illusions EN"},
		{"Fountain of Prayers:", i18n.SourceIDFountainOfPrayers, "Fountain of Prayers"},
		{"Login", i18n.SourceIDLoginBonus, "Login"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			_, sourceID := MapSourceWithID(tt.source)
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
		expectedID i18n.SourceID
		desc       string
	}{
		{"抽卡 黒葬武具ガチャ 5 次, 消耗 金幣×250000", i18n.SourceIDGacha, "Gacha TW"},
		{"開啟 上級封印寶箱 x 5", i18n.SourceIDOpen, "Open TW"},
		{"無窮之塔: 800", i18n.SourceIDTowerInfinity, "Tower of Infinity TW"},
		{"勝利", i18n.SourceIDTempleIllusions, "Temple of Illusions TW"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			_, sourceID := MapSourceWithID(tt.source)
			if sourceID != tt.expectedID {
				t.Errorf("MapSourceWithID(%q) sourceID = %d, want %d", tt.source, sourceID, tt.expectedID)
			}
		})
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
