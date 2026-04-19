package main

import (
	"regexp"
)

// DiamondRecord 表示钻石记录
type DiamondRecord struct {
	Character   string `json:"character"` // 角色名
	Timestamp   string `json:"timestamp"` // 显示用时间（秒级精度）
	PreciseTime string `json:"-"`         // 检查点用精确时间（纳秒级精度，不序列化）
	Amount      int    `json:"amount"`
	Type        string `json:"type"`             // "gain" 或 "consume"
	Source      string `json:"source,omitempty"` // 钻石来源（模糊匹配）
}

// SourceStats 按来源的统计
type SourceStats struct {
	Gain    int `json:"gain"`
	Consume int `json:"consume"`
}

// DailyStats 每日统计
type DailyStats struct {
	Date      string                 `json:"date"`
	Gain      int                    `json:"gain"`
	Consume   int                    `json:"consume"`
	NetChange int                    `json:"net_change"`
	Records   []DiamondRecord        `json:"records,omitempty"` // 可选保留详细记录
	Sources   map[string]SourceStats `json:"sources,omitempty"` // 按来源统计
}

// ProcessConfig 处理配置
type ProcessConfig struct {
	KeepRecords bool // 是否保留详细记录（默认 false，节省内存）
}

// WeeklyStats 每周统计
type WeeklyStats struct {
	Week      string                 `json:"week"`
	Gain      int                    `json:"gain"`
	Consume   int                    `json:"consume"`
	NetChange int                    `json:"net_change"`
	Sources   map[string]SourceStats `json:"sources,omitempty"` // 按来源统计
}

// MonthlyStats 每月统计
type MonthlyStats struct {
	Month     string                 `json:"month"`
	Gain      int                    `json:"gain"`
	Consume   int                    `json:"consume"`
	NetChange int                    `json:"net_change"`
	Sources   map[string]SourceStats `json:"sources,omitempty"` // 按来源统计
}

// TotalStats 总统计
type TotalStats struct {
	Gain      int                    `json:"gain"`
	Consume   int                    `json:"consume"`
	NetChange int                    `json:"net_change"`
	Sources   map[string]SourceStats `json:"sources,omitempty"` // 按来源统计
}

// LogProcessor 主处理器结构体
type LogProcessor struct {
	outputJSONPath string // 输出JSON文件路径
	stateFilePath  string // 状态文件路径
	inputLogPath   string // 输入日志文件路径
	checkpoint     string // 上次处理的时间戳检查点
}

// 正则表达式
var (
	diamondGainRegex    = regexp.MustCompile(`Diamonds\(None\) × (\d+)`)
	diamondConsumeRegex = regexp.MustCompile(`Diamonds\(None\) × -(\d+)`)
)
