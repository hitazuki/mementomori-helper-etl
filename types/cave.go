package types

// CaveStatus 洞穴状态
type CaveStatus string

const (
	CaveStatusStarted  CaveStatus = "started"  // 已执行
	CaveStatusFinished CaveStatus = "finished" // 已完成
	CaveStatusError    CaveStatus = "error"    // 异常
)

// CaveRecord 洞穴记录
type CaveRecord struct {
	Character   string     `json:"character"`
	Timestamp   string     `json:"timestamp"`
	PreciseTime string     `json:"-"`
	Status      CaveStatus `json:"status"`
	Date        string     `json:"-"` // 内部使用，不输出到 JSON
}

// CaveDailyStats 每日洞穴统计
type CaveDailyStats struct {
	Date    string       `json:"date"`
	Records []CaveRecord `json:"records,omitempty"`
	Status  CaveStatus   `json:"status"`
}
