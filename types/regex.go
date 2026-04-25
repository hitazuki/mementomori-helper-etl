package types

import (
	"mmth-etl/i18n"
	"regexp"
)

// globalI18nManager is the global i18n manager for pattern access.
var globalI18nManager *i18n.Manager

// InitI18n initializes the global i18n manager.
// This must be called before using any pattern functions.
func InitI18n(manager *i18n.Manager) {
	globalI18nManager = manager
}

// GetI18nManager returns the global i18n manager.
func GetI18nManager() *i18n.Manager {
	if globalI18nManager == nil {
		// Fallback: create default manager with English
		globalI18nManager = i18n.NewManager()
	}
	return globalI18nManager
}

// Pattern accessors - these replace the old direct regex variables.
// They fetch patterns from the i18n manager for the current language.

// DiamondRegex returns the diamond change pattern for the current language.
func DiamondRegex() *regexp.Regexp {
	return GetI18nManager().CurrentPatterns().Diamond
}

// RuneTicketRegex returns the rune ticket change pattern for the current language.
func RuneTicketRegex() *regexp.Regexp {
	return GetI18nManager().CurrentPatterns().RuneTicket
}

// UpgradePanaceaRegex returns the upgrade panacea change pattern for the current language.
func UpgradePanaceaRegex() *regexp.Regexp {
	return GetI18nManager().CurrentPatterns().UpgradePanacea
}

// CaveEnterRegex returns the cave enter pattern for the current language.
func CaveEnterRegex() *regexp.Regexp {
	return GetI18nManager().CurrentPatterns().CaveEnter
}

// CaveFinishRegex returns the cave finish pattern for the current language.
func CaveFinishRegex() *regexp.Regexp {
	return GetI18nManager().CurrentPatterns().CaveFinish
}

// CaveErrorRegex matches error patterns in cave logs (language-independent).
var CaveErrorRegex = regexp.MustCompile(`KeyNotFoundException`)

// ChallengeQuestRegex returns the quest challenge pattern for the current language.
func ChallengeQuestRegex() *regexp.Regexp {
	return GetI18nManager().CurrentPatterns().ChallengeQuest
}

// ChallengeTowerRegex returns the tower challenge pattern for the current language.
func ChallengeTowerRegex() *regexp.Regexp {
	return GetI18nManager().CurrentPatterns().ChallengeTower
}

// ChallengeSuccessRegex returns the challenge success pattern for the current language.
func ChallengeSuccessRegex() *regexp.Regexp {
	return GetI18nManager().CurrentPatterns().ChallengeSuccess
}

// ChallengeFailedRegex returns the challenge failed pattern for the current language.
func ChallengeFailedRegex() *regexp.Regexp {
	return GetI18nManager().CurrentPatterns().ChallengeFailed
}

// GachaPrefixRegex returns the gacha prefix pattern for the current language.
func GachaPrefixRegex() *regexp.Regexp {
	return GetI18nManager().CurrentPatterns().GachaPrefix
}

// OpenPrefixRegex returns the open prefix pattern for the current language.
func OpenPrefixRegex() *regexp.Regexp {
	return GetI18nManager().CurrentPatterns().OpenPrefix
}

// SystemErrorRegex matches system/error logs that should clear source context.
var SystemErrorRegex = regexp.MustCompile(`^(OnError|System\.)`)
