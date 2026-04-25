package parser

import (
	"mmth-etl/i18n"
	"mmth-etl/types"
	"testing"
)

func TestExtractChallengeRecord_English(t *testing.T) {
	mgr := i18n.NewManager()
	mgr.SetLanguage(i18n.LangEn)
	types.InitI18n(mgr)

	tests := []struct {
		body          string
		expectType    types.ChallengeType
		expectTower   types.TowerType
		expectLevel   string
		expectStatus  types.ChallengeStatus
		desc          string
	}{
		{
			body:         "Challenge 41-60 boss one time：You have failed.,  total：1, Success：0, Err: 0",
			expectType:   types.ChallengeTypeQuest,
			expectLevel:  "41-60",
			expectStatus: types.ChallengeStatusFailed,
			desc:         "EN Quest challenge failed",
		},
		{
			body:         "Challenge 36-13 boss one time：You have triumphed.,  total：1, Success：1, Err: 0",
			expectType:   types.ChallengeTypeQuest,
			expectLevel:  "36-13",
			expectStatus: types.ChallengeStatusSuccess,
			desc:         "EN Quest challenge success",
		},
		{
			body:         "Challenge Tower of Infinity 800 layer one time：You have triumphed.,  total：1, Success：1, Err: 0",
			expectType:   types.ChallengeTypeTower,
			expectTower:  types.TowerInfinity,
			expectLevel:  "800",
			expectStatus: types.ChallengeStatusSuccess,
			desc:         "EN Tower Infinity challenge",
		},
		{
			body:         "Challenge Tower of Crimson 801 layer one time：You have failed.,  total：1, Success：0, Err: 0",
			expectType:   types.ChallengeTypeTower,
			expectTower:  types.TowerCrimson,
			expectLevel:  "801",
			expectStatus: types.ChallengeStatusFailed,
			desc:         "EN Tower Crimson challenge",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			parsed := ParsedLog{Body: tt.body, Character: "test", Timestamp: "2026-04-25T10:00:00+08:00"}
			record := ExtractChallengeRecord(parsed)
			if record == nil {
				t.Fatalf("ExtractChallengeRecord returned nil")
			}

			if record.Type != tt.expectType {
				t.Errorf("Type = %v, want %v", record.Type, tt.expectType)
			}
			if record.Status != tt.expectStatus {
				t.Errorf("Status = %v, want %v", record.Status, tt.expectStatus)
			}
			if tt.expectType == types.ChallengeTypeQuest && record.Level != tt.expectLevel {
				t.Errorf("Level = %v, want %v", record.Level, tt.expectLevel)
			}
			if tt.expectType == types.ChallengeTypeTower && record.TowerType != tt.expectTower {
				t.Errorf("TowerType = %v, want %v", record.TowerType, tt.expectTower)
			}
			if tt.expectType == types.ChallengeTypeTower && record.Level != tt.expectLevel {
				t.Errorf("Level = %v, want %v", record.Level, tt.expectLevel)
			}
		})
	}
}

func TestExtractChallengeRecord_TraditionalChinese(t *testing.T) {
	mgr := i18n.NewManager()
	mgr.SetLanguage(i18n.LangTw)
	types.InitI18n(mgr)

	tests := []struct {
		body          string
		expectType    types.ChallengeType
		expectTower   types.TowerType
		expectLevel   string
		expectStatus  types.ChallengeStatus
		desc          string
	}{
		{
			body:         "挑战 36-13 boss 一次：勝利,  总次数：1, 胜利次数: 1, Err: 0",
			expectType:   types.ChallengeTypeQuest,
			expectLevel:  "36-13",
			expectStatus: types.ChallengeStatusSuccess,
			desc:         "TW Quest challenge success",
		},
		{
			body:         "挑战 業紅之塔 800 层一次: 勝利,  0/10 ,总次数: 1, 胜利次数: 1, Err: 0",
			expectType:   types.ChallengeTypeTower,
			expectTower:  types.TowerCrimson,
			expectLevel:  "800",
			expectStatus: types.ChallengeStatusSuccess,
			desc:         "TW Tower Crimson challenge",
		},
		{
			body:         "挑战 無窮之塔 801 层一次: 敗北,  1/10 ,总次数: 2, 胜利次数: 1, Err: 0",
			expectType:   types.ChallengeTypeTower,
			expectTower:  types.TowerInfinity,
			expectLevel:  "801",
			expectStatus: types.ChallengeStatusFailed,
			desc:         "TW Tower Infinity challenge failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			parsed := ParsedLog{Body: tt.body, Character: "test", Timestamp: "2026-04-25T10:00:00+08:00"}
			record := ExtractChallengeRecord(parsed)
			if record == nil {
				t.Fatalf("ExtractChallengeRecord returned nil")
			}

			if record.Type != tt.expectType {
				t.Errorf("Type = %v, want %v", record.Type, tt.expectType)
			}
			if record.Status != tt.expectStatus {
				t.Errorf("Status = %v, want %v", record.Status, tt.expectStatus)
			}
			if tt.expectType == types.ChallengeTypeQuest && record.Level != tt.expectLevel {
				t.Errorf("Level = %v, want %v", record.Level, tt.expectLevel)
			}
			if tt.expectType == types.ChallengeTypeTower && record.TowerType != tt.expectTower {
				t.Errorf("TowerType = %v, want %v", record.TowerType, tt.expectTower)
			}
			if tt.expectType == types.ChallengeTypeTower && record.Level != tt.expectLevel {
				t.Errorf("Level = %v, want %v", record.Level, tt.expectLevel)
			}
		})
	}
}

func TestExtractChallengeRecord_Japanese(t *testing.T) {
	mgr := i18n.NewManager()
	mgr.SetLanguage(i18n.LangJa)
	types.InitI18n(mgr)

	tests := []struct {
		body          string
		expectType    types.ChallengeType
		expectTower   types.TowerType
		expectLevel   string
		expectStatus  types.ChallengeStatus
		desc          string
	}{
		{
			body:         "36-13 ボスに挑戦 1 回：勝利しました、合計回数：1、勝利回数：1、エラー：0",
			expectType:   types.ChallengeTypeQuest,
			expectLevel:  "36-13",
			expectStatus: types.ChallengeStatusSuccess,
			desc:         "JA Quest challenge success",
		},
		{
			body:         "41-60 ボスに挑戦 1 回：敗北、合計回数：2、勝利回数：1、エラー：0",
			expectType:   types.ChallengeTypeQuest,
			expectLevel:  "41-60",
			expectStatus: types.ChallengeStatusFailed,
			desc:         "JA Quest challenge failed",
		},
		{
			body:         "無窮の塔 800 層に挑戦 1 回：勝利しました、合計回数：1、勝利回数：1、エラー：0",
			expectType:   types.ChallengeTypeTower,
			expectTower:  types.TowerInfinity,
			expectLevel:  "800",
			expectStatus: types.ChallengeStatusSuccess,
			desc:         "JA Tower Infinity challenge success",
		},
		{
			body:         "紅の塔 801 層に挑戦 1 回：敗北、合計回数：2、勝利回数：1、エラー：0",
			expectType:   types.ChallengeTypeTower,
			expectTower:  types.TowerCrimson,
			expectLevel:  "801",
			expectStatus: types.ChallengeStatusFailed,
			desc:         "JA Tower Crimson challenge failed",
		},
		{
			body:         "藍の塔 500 層に挑戦 1 回：勝利しました、合計回数：1、勝利回数：1、エラー：0",
			expectType:   types.ChallengeTypeTower,
			expectTower:  types.TowerAzure,
			expectLevel:  "500",
			expectStatus: types.ChallengeStatusSuccess,
			desc:         "JA Tower Azure challenge",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			parsed := ParsedLog{Body: tt.body, Character: "test", Timestamp: "2026-04-25T10:00:00+08:00"}
			record := ExtractChallengeRecord(parsed)
			if record == nil {
				t.Fatalf("ExtractChallengeRecord returned nil")
			}

			if record.Type != tt.expectType {
				t.Errorf("Type = %v, want %v", record.Type, tt.expectType)
			}
			if record.Status != tt.expectStatus {
				t.Errorf("Status = %v, want %v", record.Status, tt.expectStatus)
			}
			if tt.expectType == types.ChallengeTypeTower && record.TowerType != tt.expectTower {
				t.Errorf("TowerType = %v, want %v", record.TowerType, tt.expectTower)
			}
			if tt.expectType == types.ChallengeTypeTower && record.Level != tt.expectLevel {
				t.Errorf("Level = %v, want %v", record.Level, tt.expectLevel)
			}
		})
	}
}

func TestExtractChallengeRecord_Korean(t *testing.T) {
	mgr := i18n.NewManager()
	mgr.SetLanguage(i18n.LangKo)
	types.InitI18n(mgr)

	tests := []struct {
		body          string
		expectType    types.ChallengeType
		expectTower   types.TowerType
		expectLevel   string
		expectStatus  types.ChallengeStatus
		desc          string
	}{
		{
			body:         "36-13 보스에 도전 1회: 승리, 총 시도 횟수: 1, 승리 횟수: 1, 오류: 0",
			expectType:   types.ChallengeTypeQuest,
			expectLevel:  "36-13",
			expectStatus: types.ChallengeStatusSuccess,
			desc:         "KO Quest challenge success",
		},
		{
			body:         "41-60 보스에 도전 1회: 패배, 총 시도 횟수: 2, 승리 횟수: 1, 오류: 0",
			expectType:   types.ChallengeTypeQuest,
			expectLevel:  "41-60",
			expectStatus: types.ChallengeStatusFailed,
			desc:         "KO Quest challenge failed",
		},
		{
			body:         "무한의 탑 800 층에 도전 1회: 승리, 총 시도 횟수: 1, 승리 횟수: 1, 오류: 0",
			expectType:   types.ChallengeTypeTower,
			expectTower:  types.TowerInfinity,
			expectLevel:  "800",
			expectStatus: types.ChallengeStatusSuccess,
			desc:         "KO Tower Infinity challenge success",
		},
		{
			body:         "홍염의 탑 801 층에 도전 1회: 패배, 총 시도 횟수: 2, 승리 횟수: 1, 오류: 0",
			expectType:   types.ChallengeTypeTower,
			expectTower:  types.TowerCrimson,
			expectLevel:  "801",
			expectStatus: types.ChallengeStatusFailed,
			desc:         "KO Tower Crimson challenge failed",
		},
		{
			body:         "남청의 탑 500 층에 도전 1회: 승리, 총 시도 횟수: 1, 승리 횟수: 1, 오류: 0",
			expectType:   types.ChallengeTypeTower,
			expectTower:  types.TowerAzure,
			expectLevel:  "500",
			expectStatus: types.ChallengeStatusSuccess,
			desc:         "KO Tower Azure challenge",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			parsed := ParsedLog{Body: tt.body, Character: "test", Timestamp: "2026-04-25T10:00:00+08:00"}
			record := ExtractChallengeRecord(parsed)
			if record == nil {
				t.Fatalf("ExtractChallengeRecord returned nil")
			}

			if record.Type != tt.expectType {
				t.Errorf("Type = %v, want %v", record.Type, tt.expectType)
			}
			if record.Status != tt.expectStatus {
				t.Errorf("Status = %v, want %v", record.Status, tt.expectStatus)
			}
			if tt.expectType == types.ChallengeTypeTower && record.TowerType != tt.expectTower {
				t.Errorf("TowerType = %v, want %v", record.TowerType, tt.expectTower)
			}
			if tt.expectType == types.ChallengeTypeTower && record.Level != tt.expectLevel {
				t.Errorf("Level = %v, want %v", record.Level, tt.expectLevel)
			}
		})
	}
}
