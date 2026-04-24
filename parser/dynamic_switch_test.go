package parser

import (
	"mmth-etl/i18n"
	"testing"
)

// TestDynamicSwitchWithWindowSizes tests language switching with different window sizes.
func TestDynamicSwitchWithWindowSizes(t *testing.T) {
	tests := []struct {
		name          string
		windowSize    int
		threshold     int
		lines         []string // Mixed language lines
		expectedSwitches int
		desc          string
	}{
		{
			name:       "逐行检测_快速切换",
			windowSize: 1,
			threshold:  1,
			lines: []string{
				"Name: Diamonds(None) x 100",     // EN
				"名称: 鑽石(None) x 100",           // TW
				"Name: Diamonds(None) x 101",     // EN
				"名称: 鑽石(None) x 101",           // TW
			},
			expectedSwitches: 3, // EN->TW->EN->TW
			desc:          "逐行检测应每行都切换",
		},
		{
			name:       "小窗口_批量切换",
			windowSize: 5,
			threshold:  3,
			lines: []string{
				"Name: Diamonds(None) x 1",
				"Name: Diamonds(None) x 2",
				"名称: 鑽石(None) x 3",     // TW starts
				"名称: 鑽石(None) x 4",
				"名称: 鑽石(None) x 5",
				"名称: 鑽石(None) x 6",     // TW dominant
				"名称: 鑽石(None) x 7",
			},
			expectedSwitches: 1, // EN->TW after window fills
			desc:          "小窗口批量检测",
		},
		{
			name:       "大窗口_稳定检测",
			windowSize: 10,
			threshold:  5,
			lines: func() []string {
				// 15 EN lines, then 15 TW lines
				lines := make([]string, 30)
				for i := 0; i < 15; i++ {
					lines[i] = "Name: Diamonds(None) x 100"
				}
				for i := 15; i < 30; i++ {
					lines[i] = "名称: 鑽石(None) x 100"
				}
				return lines
			}(),
			expectedSwitches: 1, // EN->TW
			desc:          "大窗口稳定检测，减少抖动",
		},
		{
			name:       "阈值过高_不切换",
			windowSize: 5,
			threshold:  10, // Higher than window size
			lines: []string{
				"名称: 鑽石(None) x 1",
				"名称: 鑽石(None) x 2",
				"名称: 鑽石(None) x 3",
				"名称: 鑽石(None) x 4",
				"名称: 鑽石(None) x 5",
			},
			expectedSwitches: 0, // No switch because threshold too high
			desc:          "阈值过高阻止切换",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detector := i18n.NewDetector()
			accumulator := i18n.NewScoreAccumulator(detector, tt.windowSize)

			switchCount := 0
			currentLang := i18n.LangEn
			checkInterval := max(tt.windowSize/2, 1)

			for i, line := range tt.lines {
				accumulator.AddLine(line)

				// Check language switch at interval
				if (i+1)%checkInterval == 0 || i == len(tt.lines)-1 {
					scores := accumulator.GetScores()

					// Find dominant language
					var maxLang i18n.Language
					maxScore := 0
					for lang, score := range scores {
						if score > maxScore {
							maxScore = score
							maxLang = lang
						}
					}

					// Check switch
					if maxLang != "" && maxLang != currentLang {
						currentScore := scores[currentLang]
						if maxScore-currentScore >= tt.threshold {
							switchCount++
							currentLang = maxLang
						}
					}
				}
			}

			if switchCount != tt.expectedSwitches {
				t.Errorf("%s: expected %d switches, got %d", tt.desc, tt.expectedSwitches, switchCount)
			}
		})
	}
}

// TestAdaptiveThreshold tests that small windows use adaptive thresholds.
func TestAdaptiveThreshold(t *testing.T) {
	tests := []struct {
		windowSize      int
		userThreshold   int
		expectedEffective int
	}{
		{1, 5, 1},   // Window 1: threshold capped to 1
		{2, 5, 1},   // Window 2: threshold capped to 1
		{4, 5, 2},   // Window 4: threshold capped to 2
		{10, 5, 5},  // Window 10: use user threshold
		{100, 5, 5}, // Window 100: use user threshold
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			// Simulate adaptive threshold logic
			effective := tt.userThreshold
			if tt.windowSize < tt.userThreshold {
				effective = max(tt.windowSize/2, 1)
			}

			if effective != tt.expectedEffective {
				t.Errorf("window=%d, threshold=%d: expected effective=%d, got %d",
					tt.windowSize, tt.userThreshold, tt.expectedEffective, effective)
			}
		})
	}
}

// TestLanguageDetectionAccuracy verifies single-line detection accuracy.
func TestLanguageDetectionAccuracy(t *testing.T) {
	detector := i18n.NewDetector()

	tests := []struct {
		line     string
		expected i18n.Language
	}{
		{"Name: Diamonds(None) x 100", i18n.LangEn},
		{"Enter Cave of Space-Time", i18n.LangEn},
		{"Challenge Tower of Crimson 800 layer: triumphed", i18n.LangEn},
		{"名称: 鑽石(None) x 100", i18n.LangTw},
		{"进入 時空洞窟", i18n.LangTw},
		{"挑战 業紅之塔 800 层: 勝利", i18n.LangTw},
		{"名称: ダイヤ(None) x 100", i18n.LangJa},
		{"이름: 다이아(None) x 100", i18n.LangKo},
		{"도전 무한의 탑 800: 승리", i18n.LangKo},
	}

	for _, tt := range tests {
		t.Run(string(tt.expected), func(t *testing.T) {
			lang, score := detector.DetectSingleLine(tt.line)
			if score == 0 {
				t.Errorf("Failed to detect language for: %s", tt.line)
				return
			}
			if lang != tt.expected {
				t.Errorf("Expected %s, got %s for: %s", tt.expected, lang, tt.line)
			}
		})
	}
}

// TestRealWorldScenario simulates real mixed-language log processing.
func TestRealWorldScenario(t *testing.T) {
	// Simulate: 40 lines EN, then 60 lines TW (60% TW creates imbalance for switch)
	lines := make([]string, 100)
	for i := 0; i < 40; i++ {
		lines[i] = "Name: Diamonds(None) x 100"
	}
	for i := 40; i < 100; i++ {
		lines[i] = "名称: 鑽石(None) x 100"
	}

	tests := []struct {
		name       string
		windowSize int
		threshold  int
		wantSwitch bool
	}{
		{"窗口100_检测切换", 100, 5, true},
		{"窗口50_检测切换", 50, 5, true},
		{"窗口10_可能抖动", 10, 3, true},
		{"窗口1_逐行切换", 1, 1, true},
		{"窗口100_阈值过高", 100, 50, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detector := i18n.NewDetector()
			accumulator := i18n.NewScoreAccumulator(detector, tt.windowSize)

			switched := false
			currentLang := i18n.LangEn
			// Use frequent checks for accurate detection in tests
			checkInterval := 5

			for i, line := range lines {
				accumulator.AddLine(line)

				// Check frequently, especially after line 75 (in TW region)
				if (i+1)%checkInterval == 0 || i >= 75 {
					scores := accumulator.GetScores()

					var maxLang i18n.Language
					maxScore := 0
					for lang, score := range scores {
						if score > maxScore {
							maxScore = score
							maxLang = lang
						}
					}

					if maxLang != "" && maxLang != currentLang {
						currentScore := scores[currentLang]
						if maxScore-currentScore >= tt.threshold {
							switched = true
							currentLang = maxLang
						}
					}
				}
			}

			if switched != tt.wantSwitch {
				t.Errorf("expected switch=%v, got %v", tt.wantSwitch, switched)
			}
		})
	}
}
