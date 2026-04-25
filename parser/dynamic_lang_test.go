package parser

import (
	"mmth-etl/i18n"
	"mmth-etl/types"
	"testing"
)

// TestDynamicLanguageSwitch verifies that source mapping correctly
// switches languages and invalidates cache.
func TestDynamicLanguageSwitch(t *testing.T) {
	mgr := i18n.NewManager()
	types.InitI18n(mgr)

	// Test 1: English language - should match English sources
	t.Run("English_Login", func(t *testing.T) {
		mgr.SetLanguage(i18n.LangEn)
		InvalidateSourceCache()

		source := "Login"
		id := GetSourceID(source)
		if id != i18n.SourceIDLoginBonus {
			t.Errorf("English: GetSourceID(%q) = %d, want %d", source, id, i18n.SourceIDLoginBonus)
		}

		alias := GetSourceAlias(source)
		if alias != "Login Bonus" {
			t.Errorf("English: GetSourceAlias(%q) = %q, want %q", source, alias, "Login Bonus")
		}
	})

	// Test 2: Switch to Chinese - should NOT match English "Login"
	t.Run("Chinese_NoMatch_English_Login", func(t *testing.T) {
		mgr.SetLanguage(i18n.LangTw)
		InvalidateSourceCache()

		// English "Login" should NOT match in Chinese mode
		source := "Login"
		id := GetSourceID(source)
		if id != 0 {
			t.Errorf("Chinese mode: GetSourceID(%q) = %d, want 0 (English should not match)", source, id)
		}
	})

	// Test 3: Chinese should match Chinese sources
	t.Run("Chinese_Login", func(t *testing.T) {
		mgr.SetLanguage(i18n.LangTw)
		InvalidateSourceCache()

		source := "簽到獎勵:"
		id := GetSourceID(source)
		if id != i18n.SourceIDLoginBonus {
			t.Errorf("Chinese: GetSourceID(%q) = %d, want %d", source, id, i18n.SourceIDLoginBonus)
		}
	})

	// Test 4: IdentifyLogType should correctly identify Gacha/Open in different languages
	t.Run("IdentifyLogType_Respects_Language", func(t *testing.T) {
		// English mode - English "Gacha " should be identified
		mgr.SetLanguage(i18n.LangEn)
		if IdentifyLogType("Gacha test 10 times") != LogTypeGacha {
			t.Error("English Gacha should be identified in English mode")
		}

		// Chinese mode - Chinese "抽卡 " should be identified
		mgr.SetLanguage(i18n.LangTw)
		if IdentifyLogType("抽卡 測試卡池 10 次") != LogTypeGacha {
			t.Error("Chinese 抽卡 should be identified in Chinese mode")
		}

		// Chinese mode - English "Gacha " should NOT be identified as Gacha
		if IdentifyLogType("Gacha test 10 times") == LogTypeGacha {
			t.Error("English Gacha should NOT be identified as Gacha in Chinese mode")
		}

		// System error should clear source context
		if IdentifyLogType("OnError: something went wrong") != LogTypeSystemError {
			t.Error("OnError should be identified as LogTypeSystemError")
		}
	})

	// Test 5: Cache auto-invalidates on language change (without explicit InvalidateSourceCache)
	t.Run("Cache_Auto_Invalidates", func(t *testing.T) {
		mgr.SetLanguage(i18n.LangEn)
		InvalidateSourceCache()

		// Warm up cache with English
		_ = GetSourceID("Login")

		// Switch to Chinese without explicit invalidate
		mgr.SetLanguage(i18n.LangTw)

		// This should automatically detect language change and rebuild cache
		id := GetSourceID("簽到獎勵:")
		if id != i18n.SourceIDLoginBonus {
			t.Errorf("Auto-invalidate failed: GetSourceID(\"登錄\") = %d, want %d", id, i18n.SourceIDLoginBonus)
		}
	})
}

// TestSourceMatchingWithLogPatterns tests source matching against actual log patterns.
func TestSourceMatchingWithLogPatterns(t *testing.T) {
	mgr := i18n.NewManager()
	types.InitI18n(mgr)

	tests := []struct {
		lang     i18n.Language
		source   string
		expected i18n.SourceID
		desc     string
	}{
		// English patterns
		{i18n.LangEn, "Login", i18n.SourceIDLoginBonus, "EN Login"},
		{i18n.LangEn, "Auto Buy Store Items", i18n.SourceIDAutoBuyStore, "EN Auto Buy"},
		{i18n.LangEn, "Fountain of Prayers: 100 diamonds", i18n.SourceIDFountainOfPrayers, "EN Fountain"},
		{i18n.LangEn, "Get Daily 's 60 Reward", i18n.MissionGroupDailyID, "EN Daily Reward"},

		// Chinese patterns
		{i18n.LangTw, "簽到獎勵:", i18n.SourceIDLoginBonus, "TW Login"},
		{i18n.LangTw, "自动购买商城物品", i18n.SourceIDAutoBuyStore, "TW Auto Buy"},
		{i18n.LangTw, "祈願之泉: 100 diamonds", i18n.SourceIDFountainOfPrayers, "TW Fountain"},
		{i18n.LangTw, "领取 Daily 的 60 奖励", i18n.MissionGroupDailyID, "TW Daily Reward"},

		// Japanese patterns
		{i18n.LangJa, "ログイン", i18n.SourceIDLoginBonus, "JA Login"},
		{i18n.LangJa, "祈りの泉: 100 diamonds", i18n.SourceIDFountainOfPrayers, "JA Fountain"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mgr.SetLanguage(tt.lang)
			InvalidateSourceCache()

			id := GetSourceID(tt.source)
			if id != tt.expected {
				t.Errorf("GetSourceID(%q) = %d, want %d", tt.source, id, tt.expected)
			}
		})
	}
}
