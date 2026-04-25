package i18n

import (
	"fmt"
	"regexp"
)

// Detector provides language auto-detection for log files.
type Detector struct {
	signatures     map[Language][]*regexp.Regexp
	gameLogPattern *regexp.Regexp
}

// NewDetector creates a new language detector.
func NewDetector() *Detector {
	return &Detector{
		gameLogPattern: regexp.MustCompile(`\[\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\] \[.+?\(Lv\d+\)\]`),
		signatures:     buildDetectorSignatures(),
	}
}

// buildDetectorSignatures generates detection patterns from language definitions.
func buildDetectorSignatures() map[Language][]*regexp.Regexp {
	result := make(map[Language][]*regexp.Regexp)

	for lang := range languageDefinitions {
		def := languageDefinitions[lang]
		patterns := []*regexp.Regexp{
			// Name label + item name (most common, appears in every item change log)
			regexp.MustCompile(fmt.Sprintf(`%s: %s`, regexp.QuoteMeta(def.NameLabel), regexp.QuoteMeta(def.Diamond))),
			// Cave patterns
			regexp.MustCompile(regexp.QuoteMeta(def.CaveEnter)),
			regexp.MustCompile(regexp.QuoteMeta(def.CaveFinish)),
			// Challenge result keywords (appears in all challenge logs)
			regexp.MustCompile(fmt.Sprintf(`%s|%s`, regexp.QuoteMeta(def.SuccessKeyword), regexp.QuoteMeta(def.FailedKeyword))),
		}

		// Tower names pattern (language-specific, appears in challenge logs)
		// Note: JA/KO challenge logs start with tower name, not "Challenge" keyword
		towerPattern := fmt.Sprintf(`%s|%s|%s|%s|%s`,
			regexp.QuoteMeta(def.TowerInfinity),
			regexp.QuoteMeta(def.TowerAzure),
			regexp.QuoteMeta(def.TowerCrimson),
			regexp.QuoteMeta(def.TowerEmerald),
			regexp.QuoteMeta(def.TowerAmber),
		)
		patterns = append(patterns, regexp.MustCompile(towerPattern))

		result[lang] = patterns
	}

	return result
}

// Detect analyzes a sample of log lines to determine the language.
// Returns the detected language, or LangEn as fallback.
func (d *Detector) Detect(lines []string) Language {
	scores := make(map[Language]int)

	for _, line := range lines {
		// Only consider game log lines (skip system messages)
		if !d.gameLogPattern.MatchString(line) {
			continue
		}

		for lang, patterns := range d.signatures {
			for _, p := range patterns {
				if p.MatchString(line) {
					scores[lang]++
				}
			}
		}
	}

	// Find language with highest score
	var maxLang Language = LangEn
	maxScore := 0
	for lang, score := range scores {
		if score > maxScore {
			maxScore = score
			maxLang = lang
		}
	}

	return maxLang
}

// DetectFromSample detects language from a sample of lines.
// sampleSize is the maximum number of lines to analyze.
func (d *Detector) DetectFromSample(lines []string, sampleSize int) Language {
	if len(lines) > sampleSize {
		lines = lines[:sampleSize]
	}
	return d.Detect(lines)
}

// DetectSingleLine detects language from a single line.
// Returns the detected language and a confidence score.
// Note: The caller should ensure only valid game log content is passed
// (e.g., parsed.Body after IsValid check), so we don't need gameLogPattern filtering.
func (d *Detector) DetectSingleLine(line string) (Language, int) {
	scores := make(map[Language]int)
	for lang, patterns := range d.signatures {
		for _, p := range patterns {
			if p.MatchString(line) {
				scores[lang]++
			}
		}
	}

	var maxLang Language
	maxScore := 0
	for lang, score := range scores {
		if score > maxScore {
			maxScore = score
			maxLang = lang
		}
	}

	return maxLang, maxScore
}

// LanguageScores returns the language scores for a batch of lines.
// Used for dynamic language detection during processing.
func (d *Detector) LanguageScores(lines []string) map[Language]int {
	scores := make(map[Language]int)

	for _, line := range lines {
		if !d.gameLogPattern.MatchString(line) {
			continue
		}

		for lang, patterns := range d.signatures {
			for _, p := range patterns {
				if p.MatchString(line) {
					scores[lang]++
				}
			}
		}
	}

	return scores
}

// ScoreAccumulator accumulates language scores incrementally.
// Used for efficient dynamic language detection without re-scanning.
type ScoreAccumulator struct {
	scores      map[Language]int
	detector    *Detector
	windowSize  int
	scoreWindow []Language // stores detected language for each line (for sliding window)
}

// NewScoreAccumulator creates a new score accumulator.
func NewScoreAccumulator(detector *Detector, windowSize int) *ScoreAccumulator {
	return &ScoreAccumulator{
		scores:      make(map[Language]int),
		detector:    detector,
		windowSize:  windowSize,
		scoreWindow: make([]Language, 0, windowSize),
	}
}

// AddLine adds a line and updates scores incrementally.
// Returns the detected language for this line (may be empty for non-game logs).
func (a *ScoreAccumulator) AddLine(line string) Language {
	lang, score := a.detector.DetectSingleLine(line)
	if score == 0 {
		return ""
	}

	// Add to window
	a.scoreWindow = append(a.scoreWindow, lang)
	if len(a.scoreWindow) > a.windowSize {
		// Remove oldest score
		oldestLang := a.scoreWindow[0]
		a.scoreWindow = a.scoreWindow[1:]
		if oldestLang != "" {
			a.scores[oldestLang]--
		}
	}

	// Add new score
	if lang != "" {
		a.scores[lang] += score
	}

	return lang
}

// GetScores returns current accumulated scores.
func (a *ScoreAccumulator) GetScores() map[Language]int {
	result := make(map[Language]int)
	for lang, score := range a.scores {
		result[lang] = score
	}
	return result
}

// Reset clears all accumulated scores.
func (a *ScoreAccumulator) Reset() {
	a.scores = make(map[Language]int)
	a.scoreWindow = a.scoreWindow[:0]
}
