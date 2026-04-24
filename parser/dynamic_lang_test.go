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
		if id != i18n.SourceIDLogin {
			t.Errorf("English: GetSourceID(%q) = %d, want %d", source, id, i18n.SourceIDLogin)
		}

		alias := GetSourceAlias(source)
		if alias != "Login" {
			t.Errorf("English: GetSourceAlias(%q) = %q, want %q", source, alias, "Login")
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

		source := "登錄"
		id := GetSourceID(source)
		if id != i18n.SourceIDLogin {
			t.Errorf("Chinese: GetSourceID(%q) = %d, want %d", source, id, i18n.SourceIDLogin)
		}
	})

	// Test 4: IsValidSource should respect current language
	t.Run("IsValidSource_Respects_Language", func(t *testing.T) {
		// English mode - English "Name:" should be invalid (filtered)
		mgr.SetLanguage(i18n.LangEn)
		if IsValidSource("Name: Diamonds(None) x 100") {
			t.Error("English Name: should be invalid in English mode")
		}

		// Chinese mode - English "Name:" should be valid (not filtered, only 名称: is)
		mgr.SetLanguage(i18n.LangTw)
		if !IsValidSource("Name: Diamonds(None) x 100") {
			t.Error("English Name: should be valid in Chinese mode")
		}

		// Chinese mode - Chinese "名称:" should be invalid
		if IsValidSource("名称: 鑽石(None) x 100") {
			t.Error("Chinese 名称: should be invalid in Chinese mode")
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
		id := GetSourceID("登錄")
		if id != i18n.SourceIDLogin {
			t.Errorf("Auto-invalidate failed: GetSourceID(\"登錄\") = %d, want %d", id, i18n.SourceIDLogin)
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
		{i18n.LangEn, "Login", i18n.SourceIDLogin, "EN Login"},
		{i18n.LangEn, "Auto Buy Store Items", i18n.SourceIDAutoBuyStore, "EN Auto Buy"},
		{i18n.LangEn, "Fountain of Prayers: 100 diamonds", i18n.SourceIDFountainOfPrayers, "EN Fountain"},
		{i18n.LangEn, "Get Daily 's 60 Reward", i18n.MissionGroupDailyID, "EN Daily Reward"},

		// Chinese patterns
		{i18n.LangTw, "登錄", i18n.SourceIDLogin, "TW Login"},
		{i18n.LangTw, "自動購買商城物品", i18n.SourceIDAutoBuyStore, "TW Auto Buy"},
		{i18n.LangTw, "祈願之泉: 100 diamonds", i18n.SourceIDFountainOfPrayers, "TW Fountain"},
		{i18n.LangTw, "领取 Daily 的 60 奖励", i18n.MissionGroupDailyID, "TW Daily Reward"},

		// Japanese patterns
		{i18n.LangJa, "ログイン", i18n.SourceIDLogin, "JA Login"},
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
