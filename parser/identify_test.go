package parser

import (
	"mmth-etl/i18n"
	"mmth-etl/types"
	"testing"
)

func TestIdentifyLogType_English(t *testing.T) {
	mgr := i18n.NewManager()
	mgr.SetLanguage(i18n.LangEn)
	types.InitI18n(mgr)

	tests := []struct {
		body     string
		expected LogType
		desc     string
	}{
		// Diamond
		{"Name: Diamonds(None) × 100", LogTypeDiamond, "EN Diamond gain"},
		{"Name: Diamonds(None) × -50", LogTypeDiamond, "EN Diamond consume"},

		// RuneTicket
		{"Name: Rune Ticket(SSR) × 10", LogTypeRuneTicket, "EN Rune Ticket"},

		// UpgradePanacea
		{"Name: Upgrade Panacea(SSR) × 5", LogTypeUpgradePanacea, "EN Upgrade Panacea"},

		// Cave
		{"Enter Cave of Space-Time", LogTypeCave, "EN Cave enter"},
		{"Cave of Space-Time Finished", LogTypeCave, "EN Cave finish"},
		{"System.Collections.Generic.KeyNotFoundException: The given key was not present", LogTypeCave, "EN Cave system error"},

		// Challenge
		{"Challenge 36-13 boss one time：You have triumphed.,  total：1, Success：1, Err: 0", LogTypeChallenge, "EN Quest challenge success"},
		{"Challenge 41-60 boss one time：You have failed.,  total：1, Success：0, Err: 0", LogTypeChallenge, "EN Quest challenge failed"},
		{"Challenge Tower of Infinity 800 layer one time：You have triumphed.,  total：1, Success：1, Err: 0", LogTypeChallenge, "EN Tower challenge success"},
		{"Challenge Tower of Crimson 801 layer one time：You have failed.,  total：1, Success：0, Err: 0", LogTypeChallenge, "EN Tower challenge failed"},

		// Gacha
		{"Gacha 10 times", LogTypeGacha, "EN Gacha source"},
		{"Gacha test pool 10 times", LogTypeGacha, "EN Gacha with pool name"},

		// Open
		{"Open Box x 5", LogTypeOpen, "EN Open source"},
		{"Open Gold Sealed Chest x 5", LogTypeOpen, "EN Open with item name"},

		// SystemError
		{"OnError: something went wrong", LogTypeSystemError, "EN OnError clears source"},
		{"System.Exception: error", LogTypeSystemError, "EN System clears source"},

		// NameLabel - item change in different language (should not be source context)
		{"名称: 鑽石(None) × 100", LogTypeNameLabel, "TW Name prefix in EN mode"},
		{"名前: ダイヤ(None) × 100", LogTypeNameLabel, "JA Name prefix in EN mode"},
		{"이름: 다이아(None) × 100", LogTypeNameLabel, "KO Name prefix in EN mode"},

		// None
		{"Some random log message", LogTypeNone, "EN Unknown log"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			result := IdentifyLogType(tt.body)
			if result != tt.expected {
				t.Errorf("IdentifyLogType(%q) = %v, want %v", tt.body, result, tt.expected)
			}
		})
	}
}

func TestIdentifyLogType_TraditionalChinese(t *testing.T) {
	mgr := i18n.NewManager()
	mgr.SetLanguage(i18n.LangTw)
	types.InitI18n(mgr)

	tests := []struct {
		body     string
		expected LogType
		desc     string
	}{
		// Diamond
		{"名称: 鑽石(None) × 100", LogTypeDiamond, "TW Diamond gain"},
		{"名称: 鑽石(None) × -50", LogTypeDiamond, "TW Diamond consume"},

		// RuneTicket
		{"名称: 符石兌換券(SSR) × 10", LogTypeRuneTicket, "TW Rune Ticket"},

		// UpgradePanacea
		{"名称: 強化秘藥(SSR) × 5", LogTypeUpgradePanacea, "TW Upgrade Panacea"},

		// Cave
		{"进入 時空洞窟", LogTypeCave, "TW Cave enter"},
		{"時空洞窟已完成", LogTypeCave, "TW Cave finish"},

		// Challenge
		{"挑战 36-13 boss 一次：勝利,  总次数：1, 胜利次数: 1, Err: 0", LogTypeChallenge, "TW Quest challenge success"},
		{"挑战 41-60 boss 一次：敗北,  总次数：1, 胜利次数: 0, Err: 0", LogTypeChallenge, "TW Quest challenge failed"},
		{"挑战 無窮之塔 800 层一次: 勝利,  0/10 ,总次数: 1, 胜利次数: 1, Err: 0", LogTypeChallenge, "TW Tower challenge success"},
		{"挑战 業紅之塔 801 层一次: 敗北,  1/10 ,总次数: 2, 胜利次数: 1, Err: 0", LogTypeChallenge, "TW Tower challenge failed"},

		// Gacha
		{"抽卡 測試卡池 10 次", LogTypeGacha, "TW Gacha source"},
		{"抽卡 測試卡池 5 次, 消耗 金幣×250000", LogTypeGacha, "TW Gacha with cost"},

		// Open
		{"開啟 上級封印寶箱 x 5", LogTypeOpen, "TW Open source"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			result := IdentifyLogType(tt.body)
			if result != tt.expected {
				t.Errorf("IdentifyLogType(%q) = %v, want %v", tt.body, result, tt.expected)
			}
		})
	}
}

func TestIdentifyLogType_Japanese(t *testing.T) {
	mgr := i18n.NewManager()
	mgr.SetLanguage(i18n.LangJa)
	types.InitI18n(mgr)

	tests := []struct {
		body     string
		expected LogType
		desc     string
	}{
		// Diamond
		{"名前: ダイヤ(None) × 100", LogTypeDiamond, "JA Diamond gain"},

		// RuneTicket
		{"名前: ルーンチケット(SSR) × 10", LogTypeRuneTicket, "JA Rune Ticket"},

		// UpgradePanacea
		{"名前: 強化秘薬(SSR) × 5", LogTypeUpgradePanacea, "JA Upgrade Panacea"},

		// Cave
		{"入る 時空の洞窟", LogTypeCave, "JA Cave enter"},
		{"時空の洞窟 完了済み", LogTypeCave, "JA Cave finish"},

		// Challenge
		{"36-13 ボスに挑戦 1 回：勝利しました、合計回数：1、勝利回数：1、エラー：0", LogTypeChallenge, "JA Quest challenge success"},
		{"41-60 ボスに挑戦 1 回：敗北、合計回数：2、勝利回数：1、エラー：0", LogTypeChallenge, "JA Quest challenge failed"},
		{"無窮の塔 800 層に挑戦 1 回：勝利しました、合計回数：1、勝利回数：1、エラー：0", LogTypeChallenge, "JA Tower challenge success"},
		{"紅の塔 801 層に挑戦 1 回：敗北、合計回数：2、勝利回数：1、エラー：0", LogTypeChallenge, "JA Tower challenge failed"},

		// Gacha
		{"ガチャ テストガチャ 10 回", LogTypeGacha, "JA Gacha source"},

		// Open
		{"開く 上級封印宝箱 x 5", LogTypeOpen, "JA Open source"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			result := IdentifyLogType(tt.body)
			if result != tt.expected {
				t.Errorf("IdentifyLogType(%q) = %v, want %v", tt.body, result, tt.expected)
			}
		})
	}
}

func TestIdentifyLogType_Korean(t *testing.T) {
	mgr := i18n.NewManager()
	mgr.SetLanguage(i18n.LangKo)
	types.InitI18n(mgr)

	tests := []struct {
		body     string
		expected LogType
		desc     string
	}{
		// Diamond
		{"이름: 다이아(None) × 100", LogTypeDiamond, "KO Diamond gain"},

		// RuneTicket
		{"이름: 룬 티켓(SSR) × 10", LogTypeRuneTicket, "KO Rune Ticket"},

		// UpgradePanacea
		{"이름: 강화의 비약(SSR) × 5", LogTypeUpgradePanacea, "KO Upgrade Panacea"},

		// Cave
		{"입장 시공의 동굴", LogTypeCave, "KO Cave enter"},
		{"시공의 동굴 완료됨", LogTypeCave, "KO Cave finish"},

		// Challenge
		{"36-13 보스에 도전 1회: 승리, 총 시도 횟수: 1, 승리 횟수: 1, 오류: 0", LogTypeChallenge, "KO Quest challenge success"},
		{"41-60 보스에 도전 1회: 패배, 총 시도 횟수: 2, 승리 횟수: 1, 오류: 0", LogTypeChallenge, "KO Quest challenge failed"},
		{"무한의 탑 800 층에 도전 1회: 승리, 총 시도 횟수: 1, 승리 횟수: 1, 오류: 0", LogTypeChallenge, "KO Tower challenge success"},
		{"홍염의 탑 801 층에 도전 1회: 패배, 총 시도 횟수: 2, 승리 횟수: 1, 오류: 0", LogTypeChallenge, "KO Tower challenge failed"},

		// Gacha
		{"가챠 테스트가챠 10 회", LogTypeGacha, "KO Gacha source"},

		// Open
		{"열기 상급봉인상자 x 5", LogTypeOpen, "KO Open source"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			result := IdentifyLogType(tt.body)
			if result != tt.expected {
				t.Errorf("IdentifyLogType(%q) = %v, want %v", tt.body, result, tt.expected)
			}
		})
	}
}
