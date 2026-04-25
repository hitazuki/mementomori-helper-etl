package parser

import (
	"mmth-etl/i18n"
	"mmth-etl/types"
	"regexp"
	"strings"
)

// sourceTableCache holds the built lookup tables for a specific language.
type sourceTableCache struct {
	lang           i18n.Language
	entries        []i18n.SourceEntry
	lookup         map[string]i18n.SourceEntry
	rewardPatterns []i18n.RewardMissionPattern
}

var currentCache *sourceTableCache

func buildSourceCache(mgr *i18n.Manager, lang i18n.Language) *sourceTableCache {
	entries := mgr.GetSources(lang)
	patterns := mgr.GetRewardMissionPatterns(lang)

	lookup := make(map[string]i18n.SourceEntry, len(entries))
	for _, entry := range entries {
		lookup[entry.Text] = entry
	}

	return &sourceTableCache{
		lang:           lang,
		entries:        entries,
		lookup:         lookup,
		rewardPatterns: patterns,
	}
}

func ensureCache() {
	mgr := types.GetI18nManager()
	currentLang := mgr.CurrentLanguage()

	if currentCache == nil || currentCache.lang != currentLang {
		currentCache = buildSourceCache(mgr, currentLang)
	}
}

// findSource finds a source by exact or prefix match.
// Returns the entry and true if found.
func findSource(source string) (i18n.SourceEntry, bool) {
	ensureCache()

	// Exact match
	if entry, ok := currentCache.lookup[source]; ok {
		return entry, true
	}

	// Prefix match
	for _, entry := range currentCache.entries {
		if strings.HasPrefix(source, entry.Text) {
			return entry, true
		}
	}

	// Reward mission patterns
	for _, pattern := range currentCache.rewardPatterns {
		if strings.HasPrefix(source, pattern.Prefix) {
			return i18n.SourceEntry{
				ID:    pattern.SourceID,
				Alias: pattern.Alias,
				Text:  pattern.Prefix,
			}, true
		}
	}

	return i18n.SourceEntry{}, false
}

// GetSourceID returns the source ID for a given source string.
func GetSourceID(source string) i18n.SourceID {
	if entry, ok := findSource(source); ok {
		return entry.ID
	}
	return 0
}

// GetSourceAlias returns a friendly alias for a source string.
func GetSourceAlias(source string) string {
	if entry, ok := findSource(source); ok {
		return entry.Alias
	}
	return source
}

// MapSourceWithID maps source to alias and returns the source ID.
// Note: Gacha/Open sources should already be cleaned (quantity suffix removed)
// before calling this function (see CleanSourceSuffix).
func MapSourceWithID(source string) (alias string, sourceID i18n.SourceID) {
	// Gacha: check prefix (source already cleaned)
	if types.GachaPrefixRegex().MatchString(source) {
		return source, i18n.SourceIDGacha
	}

	// Open: check prefix (source already cleaned)
	if types.OpenPrefixRegex().MatchString(source) {
		return source, i18n.SourceIDOpen
	}

	// Check mapping table
	if entry, ok := findSource(source); ok {
		return entry.Alias, entry.ID
	}

	return source, 0
}

// InvalidateSourceCache clears the source cache.
func InvalidateSourceCache() {
	currentCache = nil
}

// gachaCountPattern matches gacha count suffix: " <count> times/次/回/회"
var gachaCountPattern = regexp.MustCompile(` \d+ (times|次|回|회)`)

// CleanSourceSuffix removes quantity suffix from source based on log type.
// This is used when storing source context, avoiding redundant prefix matching.
// For LogTypeGacha: cuts at " <count> times/次/回/회" pattern
// For LogTypeOpen: cuts at " x" pattern (e.g., "Open Box x 5" -> "Open Box")
// For other types: returns original body
func CleanSourceSuffix(body string, logType LogType) string {
	switch logType {
	case LogTypeGacha:
		// Find " <count> times/次/回/회" pattern and cut before it
		if idx := gachaCountPattern.FindStringIndex(body); idx != nil {
			return strings.TrimSpace(body[:idx[0]])
		}
		return strings.TrimSpace(body)

	case LogTypeOpen:
		// Cut at " x" pattern
		if idx := strings.Index(body, " x"); idx != -1 {
			return strings.TrimSpace(body[:idx])
		}
		return strings.TrimSpace(body)

	default:
		return body
	}
}
