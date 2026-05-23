package aggregator

import (
	"mmth-etl/types"
	"testing"
)

// makeRecord 构造一条洞穴记录的辅助函数
func makeRecord(character, date, timestamp string, status types.CaveStatus) types.CaveRecord {
	return types.CaveRecord{
		Character: character,
		Timestamp: timestamp,
		Status:    status,
		Date:      date,
	}
}

// TestCaveStatusSingleRecord 单条记录：状态直接采用该记录
func TestCaveStatusSingleRecord(t *testing.T) {
	cases := []struct {
		name   string
		status types.CaveStatus
	}{
		{"仅 started", types.CaveStatusStarted},
		{"仅 finished", types.CaveStatusFinished},
		{"仅 error", types.CaveStatusError},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			agg := NewCaveAggregator()
			agg.AddRecord(makeRecord("角色A", "2026-05-23", "2026-05-23T08:00:00+08:00", tc.status))

			got := agg.caveStats["角色A"]["2026-05-23"].Status
			if got != tc.status {
				t.Fatalf("status = %q, want %q", got, tc.status)
			}
		})
	}
}

// TestCaveStatusLatestFinishedWinsOverEarlierError
// 场景：早上 error，晚上 finished → 最终应为 finished
func TestCaveStatusLatestFinishedWinsOverEarlierError(t *testing.T) {
	agg := NewCaveAggregator()
	agg.AddRecord(makeRecord("角色A", "2026-05-23", "2026-05-23T08:00:00+08:00", types.CaveStatusError))
	agg.AddRecord(makeRecord("角色A", "2026-05-23", "2026-05-23T20:00:00+08:00", types.CaveStatusFinished))

	got := agg.caveStats["角色A"]["2026-05-23"].Status
	if got != types.CaveStatusFinished {
		t.Fatalf("status = %q, want %q", got, types.CaveStatusFinished)
	}
}

// TestCaveStatusLatestErrorWinsOverEarlierFinished
// 场景：早上 finished，晚上 error → 最终应为 error
func TestCaveStatusLatestErrorWinsOverEarlierFinished(t *testing.T) {
	agg := NewCaveAggregator()
	agg.AddRecord(makeRecord("角色A", "2026-05-23", "2026-05-23T08:00:00+08:00", types.CaveStatusFinished))
	agg.AddRecord(makeRecord("角色A", "2026-05-23", "2026-05-23T20:00:00+08:00", types.CaveStatusError))

	got := agg.caveStats["角色A"]["2026-05-23"].Status
	if got != types.CaveStatusError {
		t.Fatalf("status = %q, want %q", got, types.CaveStatusError)
	}
}

// TestCaveStatusLatestStartedAfterFinished
// 场景：先 finished，后 started → 最终应为 started（以最新记录为准）
func TestCaveStatusLatestStartedAfterFinished(t *testing.T) {
	agg := NewCaveAggregator()
	agg.AddRecord(makeRecord("角色A", "2026-05-23", "2026-05-23T08:00:00+08:00", types.CaveStatusFinished))
	agg.AddRecord(makeRecord("角色A", "2026-05-23", "2026-05-23T23:00:00+08:00", types.CaveStatusStarted))

	got := agg.caveStats["角色A"]["2026-05-23"].Status
	if got != types.CaveStatusStarted {
		t.Fatalf("status = %q, want %q", got, types.CaveStatusStarted)
	}
}

// TestCaveStatusMultipleRecordsTimestampOrder
// 场景：乱序添加三条记录，应以时间戳最大的那条为准
func TestCaveStatusMultipleRecordsTimestampOrder(t *testing.T) {
	agg := NewCaveAggregator()
	// 故意乱序添加
	agg.AddRecord(makeRecord("角色A", "2026-05-23", "2026-05-23T15:00:00+08:00", types.CaveStatusStarted))
	agg.AddRecord(makeRecord("角色A", "2026-05-23", "2026-05-23T08:00:00+08:00", types.CaveStatusError))
	agg.AddRecord(makeRecord("角色A", "2026-05-23", "2026-05-23T22:00:00+08:00", types.CaveStatusFinished))

	got := agg.caveStats["角色A"]["2026-05-23"].Status
	if got != types.CaveStatusFinished {
		t.Fatalf("status = %q, want %q (最新时间戳 22:00 对应 finished)", got, types.CaveStatusFinished)
	}
}

// TestCaveStatusIsolatedByCharacterAndDate
// 场景：不同角色、不同日期的状态互不干扰
func TestCaveStatusIsolatedByCharacterAndDate(t *testing.T) {
	agg := NewCaveAggregator()

	// 角色A 05-22：error
	agg.AddRecord(makeRecord("角色A", "2026-05-22", "2026-05-22T10:00:00+08:00", types.CaveStatusError))
	// 角色A 05-23：finished
	agg.AddRecord(makeRecord("角色A", "2026-05-23", "2026-05-23T10:00:00+08:00", types.CaveStatusFinished))
	// 角色B 05-23：started
	agg.AddRecord(makeRecord("角色B", "2026-05-23", "2026-05-23T10:00:00+08:00", types.CaveStatusStarted))

	checks := []struct {
		character string
		date      string
		want      types.CaveStatus
	}{
		{"角色A", "2026-05-22", types.CaveStatusError},
		{"角色A", "2026-05-23", types.CaveStatusFinished},
		{"角色B", "2026-05-23", types.CaveStatusStarted},
	}

	for _, c := range checks {
		got := agg.caveStats[c.character][c.date].Status
		if got != c.want {
			t.Errorf("caveStats[%q][%q].Status = %q, want %q", c.character, c.date, got, c.want)
		}
	}
}

// TestCaveRecordCountAccumulates AddRecord 应正确累加记录数
func TestCaveRecordCountAccumulates(t *testing.T) {
	agg := NewCaveAggregator()
	agg.AddRecord(makeRecord("角色A", "2026-05-23", "2026-05-23T08:00:00+08:00", types.CaveStatusStarted))
	agg.AddRecord(makeRecord("角色A", "2026-05-23", "2026-05-23T20:00:00+08:00", types.CaveStatusFinished))
	agg.AddRecord(makeRecord("角色B", "2026-05-23", "2026-05-23T10:00:00+08:00", types.CaveStatusError))

	if got := agg.RecordCount(); got != 3 {
		t.Fatalf("RecordCount() = %d, want 3", got)
	}
}
