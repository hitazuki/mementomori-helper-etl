package i18n

import (
	"fmt"
	"regexp"
)

// Detector provides language auto-detection for log files.
type Detector struct {
	signatures     map[Language][]weightedSignature
	gameLogPattern *regexp.Regexp
}

type weightedSignature struct {
	pattern *regexp.Regexp
	weight  int
}

type languageScore struct {
	lang  Language
	score int
}

const (
	minSingleLineScore  = 2
	minSingleLineMargin = 2
)

// NewDetector creates a new language detector.
func NewDetector() *Detector {
	return &Detector{
		gameLogPattern: regexp.MustCompile(`\[\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\] \[.+?\(Lv\d+\)\]`),
		signatures:     buildDetectorSignatures(),
	}
}

// buildDetectorSignatures generates detection patterns from language definitions.
func buildDetectorSignatures() map[Language][]weightedSignature {
	result := make(map[Language][]weightedSignature)

	for lang := range languageDefinitions {
		def := languageDefinitions[lang]
		patterns := []weightedSignature{
			// 物品行标签在真实日志中大量出现，单独标签也能较稳定地区分语言。
			{regexp.MustCompile(fmt.Sprintf(`^%s:`, regexp.QuoteMeta(def.NameLabel))), 2},
			// 钻石物品名是高置信信号，适合用于行级提示。
			{regexp.MustCompile(fmt.Sprintf(`^%s: %s`, regexp.QuoteMeta(def.NameLabel), regexp.QuoteMeta(def.Diamond))), 3},
			// 洞窟完整短语短且唯一，权重高。
			{regexp.MustCompile(regexp.QuoteMeta(def.CaveEnter)), 3},
			{regexp.MustCompile(regexp.QuoteMeta(def.CaveFinish)), 3},
			// 单个胜负词可能被多语言共享，只作为窗口累计的弱信号。
			{regexp.MustCompile(fmt.Sprintf(`%s|%s`, regexp.QuoteMeta(def.SuccessKeyword), regexp.QuoteMeta(def.FailedKeyword))), 1},
		}

		// 塔名通常是完整挑战行的一部分，区分度高于单个胜负词。
		towerPattern := fmt.Sprintf(`%s|%s|%s|%s|%s`,
			regexp.QuoteMeta(def.TowerInfinity),
			regexp.QuoteMeta(def.TowerAzure),
			regexp.QuoteMeta(def.TowerCrimson),
			regexp.QuoteMeta(def.TowerEmerald),
			regexp.QuoteMeta(def.TowerAmber),
		)
		patterns = append(patterns, weightedSignature{regexp.MustCompile(towerPattern), 2})

		if def.ChallengeKeyword != "" {
			patterns = append(patterns, weightedSignature{regexp.MustCompile(regexp.QuoteMeta(def.ChallengeKeyword)), 2})
		}
		if def.GachaPrefix != "" {
			patterns = append(patterns, weightedSignature{regexp.MustCompile(`^` + regexp.QuoteMeta(def.GachaPrefix)), 2})
		}
		if def.OpenPrefix != "" {
			patterns = append(patterns, weightedSignature{regexp.MustCompile(`^` + regexp.QuoteMeta(def.OpenPrefix)), 2})
		}

		for _, source := range SourceDefinitions[lang] {
			if len([]rune(source.Text)) >= 6 {
				patterns = append(patterns, weightedSignature{regexp.MustCompile(regexp.QuoteMeta(source.Text)), 2})
			}
		}

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

		for lang, score := range d.scoreLine(line) {
			scores[lang] += score
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
	scores := d.scoreLine(line)
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

// DetectSingleLineUnique detects the language for one line only when there is a
// single highest-scoring language. Ambiguous lines are ignored so they do not
// temporarily flip the parser into the wrong language.
func (d *Detector) DetectSingleLineUnique(line string) (Language, int) {
	scores := d.scoreLine(line)
	var maxLang Language
	maxScore := 0
	secondScore := 0
	for lang, score := range scores {
		if score > maxScore {
			secondScore = maxScore
			maxScore = score
			maxLang = lang
		} else if score > secondScore {
			secondScore = score
		}
	}

	if maxScore < minSingleLineScore || maxScore-secondScore < minSingleLineMargin {
		return "", 0
	}
	return maxLang, maxScore
}

func (d *Detector) scoreLine(line string) map[Language]int {
	scores := make(map[Language]int)
	for lang, patterns := range d.signatures {
		for _, sig := range patterns {
			if sig.pattern.MatchString(line) {
				scores[lang] += sig.weight
			}
		}
	}
	return scores
}

// LanguageScores returns the language scores for a batch of lines.
// Used for dynamic language detection during processing.
func (d *Detector) LanguageScores(lines []string) map[Language]int {
	scores := make(map[Language]int)

	for _, line := range lines {
		if !d.gameLogPattern.MatchString(line) {
			continue
		}

		for lang, score := range d.scoreLine(line) {
			scores[lang] += score
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
	scoreWindow []languageScore // stores detected language scores for sliding window
}

// NewScoreAccumulator creates a new score accumulator.
func NewScoreAccumulator(detector *Detector, windowSize int) *ScoreAccumulator {
	return &ScoreAccumulator{
		scores:      make(map[Language]int),
		detector:    detector,
		windowSize:  windowSize,
		scoreWindow: make([]languageScore, 0, windowSize),
	}
}

// AddLine adds a line and updates scores incrementally.
// Returns the detected language for this line (may be empty for non-game logs).
func (a *ScoreAccumulator) AddLine(line string) Language {
	lang, score := a.detector.DetectSingleLineUnique(line)
	if score == 0 {
		return ""
	}

	// Add to window
	a.scoreWindow = append(a.scoreWindow, languageScore{lang: lang, score: score})
	if len(a.scoreWindow) > a.windowSize {
		// Remove oldest score
		oldest := a.scoreWindow[0]
		a.scoreWindow = a.scoreWindow[1:]
		if oldest.lang != "" {
			a.scores[oldest.lang] -= oldest.score
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
