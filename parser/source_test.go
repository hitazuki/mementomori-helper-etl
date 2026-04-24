package parser

import (
	"mmth-etl/i18n"
	"mmth-etl/types"
	"testing"
)

func TestIsValidSource_English(t *testing.T) {
	mgr := i18n.NewManager()
	mgr.SetLanguage(i18n.LangEn)
	types.InitI18n(mgr)

	tests := []struct {
		body     string
		expected bool
		desc     string
	}{
		{"Name: Diamonds(None) × 100", false, "EN Name prefix"},
		{"Challenge Tower of Infinity 800 layer", false, "EN Challenge prefix"},
		{"Gacha 10 times", true, "EN Gacha is valid source"},
		{"Open Box x 5", true, "EN Open is valid source"},
		{"OnError: something went wrong", false, "OnError prefix"},
		{"Enter Cave of Space-Time", true, "Cave enter is valid source"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			result := IsValidSource(tt.body)
			if result != tt.expected {
				t.Errorf("IsValidSource(%q) = %v, want %v", tt.body, result, tt.expected)
			}
		})
	}
}

func TestIsValidSource_TraditionalChinese(t *testing.T) {
	mgr := i18n.NewManager()
	mgr.SetLanguage(i18n.LangTw)
	types.InitI18n(mgr)

	tests := []struct {
		body     string
		expected bool
		desc     string
	}{
		{"名称: 鑽石(None) × 100", false, "TW Name prefix"},
		{"挑战 業紅之塔 800 层一次: 勝利", false, "TW Challenge prefix"},
		{"進入 時空洞窟", true, "TW Cave enter is valid source"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			result := IsValidSource(tt.body)
			if result != tt.expected {
				t.Errorf("IsValidSource(%q) = %v, want %v", tt.body, result, tt.expected)
			}
		})
	}
}

func TestIsValidSource_Japanese(t *testing.T) {
	mgr := i18n.NewManager()
	mgr.SetLanguage(i18n.LangJa)
	types.InitI18n(mgr)

	tests := []struct {
		body     string
		expected bool
		desc     string
	}{
		{"名称: ダイヤ(None) × 100", false, "JA Name prefix"},
		{"挑戦 無窮の塔 800", false, "JA Challenge prefix"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			result := IsValidSource(tt.body)
			if result != tt.expected {
				t.Errorf("IsValidSource(%q) = %v, want %v", tt.body, result, tt.expected)
			}
		})
	}
}

func TestIsValidSource_Korean(t *testing.T) {
	mgr := i18n.NewManager()
	mgr.SetLanguage(i18n.LangKo)
	types.InitI18n(mgr)

	tests := []struct {
		body     string
		expected bool
		desc     string
	}{
		{"이름: 다이아(None) × 100", false, "KO Name prefix"},
		{"도전 무한의 탑 800", false, "KO Challenge prefix"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			result := IsValidSource(tt.body)
			if result != tt.expected {
				t.Errorf("IsValidSource(%q) = %v, want %v", tt.body, result, tt.expected)
			}
		})
	}
}

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
		{"Login", i18n.SourceIDLogin, "Login"},
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
		{"登錄", i18n.SourceIDLogin, "Login TW"},
		{"自動購買商城物品", i18n.SourceIDAutoBuyStore, "Auto Buy TW"},
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
		{"ログイン", i18n.SourceIDLogin, "Login JA"},
		{"自動購入ストアアイテム", i18n.SourceIDAutoBuyStore, "Auto Buy JA"},
		{"残り挑戦回数がありません", i18n.SourceIDMissionsClaimed, "Missions Claimed JA"},
		{"時空の洞窟完了", i18n.SourceIDMissionsClaimed, "Cave Finished JA"},
		{"無窮の塔:", i18n.SourceIDTowerInfinity, "Tower of Infinity JA"},
		{"勝利しました", i18n.SourceIDTempleIllusions, "Temple of Illusions JA"},
		{"今月の合計ログイン日数：", i18n.SourceIDTotalLogins, "Total Logins JA"},
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
		{"Login", i18n.SourceIDLogin, "Login"},
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
