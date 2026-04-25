// Package i18n provides internationalization support for log parsing.
package i18n

import (
	"regexp"
)

// Language represents a supported language for log parsing.
type Language string

const (
	// LangEn is English.
	LangEn Language = "en"
	// LangTw is Traditional Chinese.
	LangTw Language = "tw"
	// LangJa is Japanese.
	LangJa Language = "ja"
	// LangKo is Korean.
	LangKo Language = "ko"
)

// PatternSet contains all regex patterns for a specific language.
type PatternSet struct {
	// Item change patterns
	Diamond        *regexp.Regexp
	RuneTicket     *regexp.Regexp
	UpgradePanacea *regexp.Regexp

	// Cave patterns
	CaveEnter  *regexp.Regexp
	CaveFinish *regexp.Regexp

	// Challenge patterns
	ChallengeQuest   *regexp.Regexp
	ChallengeTower   *regexp.Regexp
	ChallengeSuccess *regexp.Regexp
	ChallengeFailed  *regexp.Regexp

	// Source context patterns
	GachaPrefix *regexp.Regexp // Gacha prefix for source
	OpenPrefix  *regexp.Regexp // Open prefix for source
}

// Manager manages multi-language patterns for log parsing.
type Manager struct {
	currentLang Language
	patterns    map[Language]*PatternSet
}

// NewManager creates a new i18n manager with patterns for all supported languages.
func NewManager() *Manager {
	m := &Manager{
		patterns: make(map[Language]*PatternSet),
	}

	// Pre-compile patterns for all languages
	for _, lang := range []Language{LangEn, LangTw, LangJa, LangKo} {
		m.patterns[lang] = buildPatternSet(lang)
	}

	// Default to English
	m.currentLang = LangEn

	return m
}

// SetLanguage sets the current language for pattern matching.
func (m *Manager) SetLanguage(lang Language) {
	m.currentLang = lang
}

// CurrentLanguage returns the current language.
func (m *Manager) CurrentLanguage() Language {
	return m.currentLang
}

// CurrentPatterns returns the pattern set for the current language.
func (m *Manager) CurrentPatterns() *PatternSet {
	return m.patterns[m.currentLang]
}

// GetPatterns returns the pattern set for a specific language.
func (m *Manager) GetPatterns(lang Language) *PatternSet {
	if patterns, ok := m.patterns[lang]; ok {
		return patterns
	}
	return m.patterns[LangEn]
}

// GetVictoryPhrases returns all victory phrases in all languages.
// Used for source mapping (e.g., "You have triumphed." in English).
func (m *Manager) GetVictoryPhrases() map[string]string {
	result := make(map[string]string)
	for _, lang := range []Language{LangEn, LangTw, LangJa, LangKo} {
		def := languageDefinitions[lang]
		result[def.SuccessKeyword] = def.SuccessKeyword
	}
	return result
}
