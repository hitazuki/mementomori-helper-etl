package main

import (
	"regexp"
)

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

// ChallengeType 挑战类型
type ChallengeType string

const (
	ChallengeTypeQuest ChallengeType = "quest" // 主线
	ChallengeTypeTower ChallengeType = "tower" // 塔
)

// TowerType 塔类型
type TowerType string

const (
	TowerInfinity TowerType = "Infinity"
	TowerAzure    TowerType = "Azure"
	TowerCrimson  TowerType = "Crimson"
	TowerEmerald  TowerType = "Emerald"
	TowerAmber    TowerType = "Amber"
)

// ChallengeStatus 挑战状态
type ChallengeStatus string

const (
	ChallengeStatusSuccess ChallengeStatus = "success"
	ChallengeStatusFailed  ChallengeStatus = "failed"
)

// ChallengeRecord 挑战记录（解析用，不持久化）
type ChallengeRecord struct {
	Character   string          `json:"-"`
	Timestamp   string          `json:"-"`
	PreciseTime string          `json:"-"`
	Type        ChallengeType   `json:"-"`
	Level       string          `json:"level"`
	TowerType   TowerType       `json:"tower_type,omitempty"`
	Status      ChallengeStatus `json:"-"`
}

// ChallengeLevelStats 单关卡统计
type ChallengeLevelStats struct {
	Level    string `json:"level"`
	Attempts int    `json:"attempts"` // 尝试次数
	Success  bool   `json:"success"`  // 是否成功过
	LastTime string `json:"last_time,omitempty"` // 最后挑战时间
}

// ChallengeStats 角色挑战统计
type ChallengeStats struct {
	Quest  map[string]*ChallengeLevelStats               `json:"quest"`  // level -> stats
	Towers map[TowerType]map[string]*ChallengeLevelStats `json:"towers"` // tower -> level -> stats
}

// ItemRecord 物品记录（饼干/红水）
type ItemRecord struct {
	Character   string `json:"character"`
	Timestamp   string `json:"timestamp"`
	PreciseTime string `json:"-"`
	Amount      int    `json:"amount"` // 正数获取，负数消耗
	Source      string `json:"source,omitempty"`
}

// ItemDailyStats 物品每日统计
type ItemDailyStats struct {
	Date      string                 `json:"date"`
	Gain      int                    `json:"gain"`
	Consume   int                    `json:"consume"`
	NetChange int                    `json:"net_change"`
	Sources   map[string]SourceStats `json:"sources,omitempty"`
}

// ItemWeeklyStats 物品每周统计
type ItemWeeklyStats struct {
	Week      string                 `json:"week"`
	Gain      int                    `json:"gain"`
	Consume   int                    `json:"consume"`
	NetChange int                    `json:"net_change"`
	Sources   map[string]SourceStats `json:"sources,omitempty"`
}

// ItemMonthlyStats 物品每月统计
type ItemMonthlyStats struct {
	Month     string                 `json:"month"`
	Gain      int                    `json:"gain"`
	Consume   int                    `json:"consume"`
	NetChange int                    `json:"net_change"`
	Sources   map[string]SourceStats `json:"sources,omitempty"`
}

// ItemTotalStats 物品总计
type ItemTotalStats struct {
	Gain      int                    `json:"gain"`
	Consume   int                    `json:"consume"`
	NetChange int                    `json:"net_change"`
	Sources   map[string]SourceStats `json:"sources,omitempty"`
}

// ItemStats 单个物品的统计
type ItemStats struct {
	Daily   map[string]*ItemDailyStats   `json:"daily"`
	Weekly  map[string]*ItemWeeklyStats  `json:"weekly"`
	Monthly map[string]*ItemMonthlyStats `json:"monthly"`
	Total   *ItemTotalStats              `json:"total"`
}

// 正则表达式
var (
	diamondGainRegex    = regexp.MustCompile(`Diamonds\(None\) × (\d+)`)
	diamondConsumeRegex = regexp.MustCompile(`Diamonds\(None\) × -(\d+)`)
	caveEnterRegex      = regexp.MustCompile(`Enter Cave of Space-Time`)
	caveFinishRegex     = regexp.MustCompile(`Cave of Space-Time Finished`)
	caveErrorRegex      = regexp.MustCompile(`KeyNotFoundException`)

	// 挑战日志正则
	challengeQuestRegex = regexp.MustCompile(`^Challenge (\d+-\d+) boss`)
	challengeTowerRegex = regexp.MustCompile(`^Challenge Tower of (Infinity|Azure|Crimson|Emerald|Amber) (\d+) layer`)
	challengeSuccessRegex = regexp.MustCompile(`triumphed`)
	challengeFailedRegex  = regexp.MustCompile(`failed`)

	// 物品变动日志正则
	runeTicketRegex     = regexp.MustCompile(`^Name: Rune Ticket\([A-Z]+\) × (-?\d+)`)
	upgradePanaceaRegex = regexp.MustCompile(`^Name: Upgrade Panacea\([A-Z]+\) × (-?\d+)`)
)
