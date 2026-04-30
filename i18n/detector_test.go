package i18n

import "testing"

func TestDetectSingleLineUniqueWeightedSignals(t *testing.T) {
	detector := NewDetector()

	tests := []struct {
		name      string
		line      string
		wantLang  Language
		minScore  int
		wantEmpty bool
	}{
		{
			name:     "英文通用物品行",
			line:     "Name: Gold(None) \u00d7 100",
			wantLang: LangEn,
			minScore: 2,
		},
		{
			name:     "英文洞窟行",
			line:     "Enter Cave of Space-Time",
			wantLang: LangEn,
			minScore: 3,
		},
		{
			name:     "繁中钻石物品行",
			line:     "\u540d\u79f0: \u947d\u77f3(None) \u00d7 100",
			wantLang: LangTw,
			minScore: 5,
		},
		{
			name:     "繁中洞窟行",
			line:     "\u8fdb\u5165 \u6642\u7a7a\u6d1e\u7a9f",
			wantLang: LangTw,
			minScore: 3,
		},
		{
			name:      "跨语言共享胜利词不做行级提示",
			line:      "\u52dd\u5229",
			wantEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lang, score := detector.DetectSingleLineUnique(tt.line)
			if tt.wantEmpty {
				if lang != "" || score != 0 {
					t.Fatalf("lang=%s score=%d, want empty", lang, score)
				}
				return
			}
			if lang != tt.wantLang {
				t.Fatalf("lang=%s, want %s", lang, tt.wantLang)
			}
			if score < tt.minScore {
				t.Fatalf("score=%d, want >= %d", score, tt.minScore)
			}
		})
	}
}

func TestScoreAccumulatorUsesWeightedWindow(t *testing.T) {
	detector := NewDetector()
	accumulator := NewScoreAccumulator(detector, 2)

	accumulator.AddLine("Name: Gold(None) \u00d7 100")
	accumulator.AddLine("Enter Cave of Space-Time")
	accumulator.AddLine("\u540d\u79f0: \u947d\u77f3(None) \u00d7 100")

	scores := accumulator.GetScores()
	if scores[LangEn] != 3 {
		t.Fatalf("English score=%d, want 3 after first line slides out", scores[LangEn])
	}
	if scores[LangTw] != 5 {
		t.Fatalf("Traditional Chinese score=%d, want 5", scores[LangTw])
	}
}
