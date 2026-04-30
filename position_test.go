package main

import (
	"os"
	"testing"
)

func TestFindStartPositionSkipsInvalidLines(t *testing.T) {
	file, err := os.CreateTemp(t.TempDir(), "app-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	content := "system output without timestamp\n" +
		"[2026-04-30 00:50:15] [a(Lv1)] already processed\n" +
		"wrapped request payload\n" +
		"[2026-04-30 01:50:55] [a(Lv1)] 进入 時空洞窟\n"
	if _, err := file.WriteString(content); err != nil {
		t.Fatal(err)
	}

	pos, err := findStartPosition(file, "2026-04-30T00:50:15+08:00")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := file.Seek(pos, 0); err != nil {
		t.Fatal(err)
	}
	line := make([]byte, len("[2026-04-30 01:50:55]"))
	if _, err := file.Read(line); err != nil {
		t.Fatal(err)
	}
	if got := string(line); got != "[2026-04-30 01:50:55]" {
		t.Fatalf("start line = %q", got)
	}
}
