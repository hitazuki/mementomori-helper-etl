package types

// ChangeRecord 物品变动记录（钻石/饼干/红水通用）
type ChangeRecord struct {
	Character string `json:"character"`
	Timestamp string `json:"timestamp"`
	Amount    int    `json:"amount"` // 正数=获取，负数=消耗
	Source    string `json:"source,omitempty"`
}

// SourceStats 按来源的统计
type SourceStats struct {
	Gain    int `json:"gain"`
	Consume int `json:"consume"`
}

// ChangeStats 变动统计（嵌入型）
type ChangeStats struct {
	Gain      int                    `json:"gain"`
	Consume   int                    `json:"consume"`
	NetChange int                    `json:"net_change"`
	Sources   map[string]SourceStats `json:"sources,omitempty"`
}

// DailyStats 每日统计
type DailyStats struct {
	Date    string         `json:"date"`
	Records []ChangeRecord `json:"records,omitempty"`
	ChangeStats
}

// WeeklyStats 每周统计
type WeeklyStats struct {
	Week string `json:"week"`
	ChangeStats
}

// MonthlyStats 每月统计
type MonthlyStats struct {
	Month string `json:"month"`
	ChangeStats
}

// TotalStats 总计
type TotalStats struct {
	ChangeStats
}

// ProcessConfig 处理配置
type ProcessConfig struct {
	KeepRecords bool // 是否保留详细记录（默认 false，节省内存）
}
