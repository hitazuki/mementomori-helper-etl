package main

import (
	"mmth-etl/aggregator"
	"mmth-etl/i18n"
	"mmth-etl/types"
	"os"
	"path/filepath"
	"testing"
)

func TestProcessUsesLineLanguageBeforeWindowSwitch(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "app.log")
	content := "[2026-04-30 00:50:15] [test(Lv1)] already processed\n" +
		"[2026-04-30 01:50:55] [test(Lv1)] 进入 時空洞窟\n"
	if err := os.WriteFile(logPath, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	mgr := i18n.NewManager()
	mgr.SetLanguage(i18n.LangEn)
	types.InitI18n(mgr)

	processor := NewLogProcessor(logPath, "2026-04-30T00:50:15+08:00", mgr, DynamicLanguageConfig{
		Enabled:         true,
		WindowSize:      100,
		SwitchThreshold: 5,
	})

	caveAgg := aggregator.NewCaveAggregator()
	processor.Process(
		aggregator.NewChangeAggregator(),
		caveAgg,
		aggregator.NewChallengeAggregator(),
		aggregator.NewChangeAggregator(),
		aggregator.NewChangeAggregator(),
		nil,
	)

	if got := caveAgg.RecordCount(); got != 1 {
		t.Fatalf("cave records = %d, want 1", got)
	}
	if got := mgr.CurrentLanguage(); got != i18n.LangEn {
		t.Fatalf("stable language = %s, want %s", got, i18n.LangEn)
	}
}
