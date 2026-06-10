package aggregator

import (
	"encoding/json"
	"mmth-etl/types"
	"os"
	"sort"
)

// CaveAggregator 洞穴聚合器
type CaveAggregator struct {
	caveStats   map[string]map[string]*types.CaveDailyStats // character -> date -> stats
	recordCount int
}

// NewCaveAggregator 创建洞穴聚合器
func NewCaveAggregator() *CaveAggregator {
	return &CaveAggregator{
		caveStats: make(map[string]map[string]*types.CaveDailyStats),
	}
}

// AddRecord 添加洞穴记录
func (a *CaveAggregator) AddRecord(record types.CaveRecord) {
	character := record.Character
	date := record.Date

	if a.caveStats[character] == nil {
		a.caveStats[character] = make(map[string]*types.CaveDailyStats)
	}

	daily := a.caveStats[character][date]
	if daily == nil {
		daily = &types.CaveDailyStats{
			Date:    date,
			Records: []types.CaveRecord{},
		}
		a.caveStats[character][date] = daily
	}

	// 统一逻辑：如果新状态是 Finished 或 Error，当前状态必须是 Started 才能生效，否则忽略
	if record.Status != types.CaveStatusStarted && daily.Status != types.CaveStatusStarted {
		return
	}

	daily.Records = append(daily.Records, record)

	a.updateDailyStatus(character, date)
	a.recordCount++
}

func (a *CaveAggregator) updateDailyStatus(character, date string) {
	daily := a.caveStats[character][date]

	// 根据当天最新一条记录的状态决定当日状态
	var latestTimestamp string
	var latestStatus types.CaveStatus

	for _, r := range daily.Records {
		if r.Timestamp >= latestTimestamp {
			latestTimestamp = r.Timestamp
			latestStatus = r.Status
		}
	}

	if latestTimestamp != "" {
		daily.Status = latestStatus
	}
}

// LoadExistingStats 加载现有统计数据
func (a *CaveAggregator) LoadExistingStats(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	var existingStats map[string]map[string]*types.CaveDailyStats
	if err := json.NewDecoder(file).Decode(&existingStats); err != nil {
		return
	}

	for character, dates := range existingStats {
		if a.caveStats[character] == nil {
			a.caveStats[character] = make(map[string]*types.CaveDailyStats)
		}
		for date, stats := range dates {
			a.caveStats[character][date] = stats
		}
	}
}

// ToMap 转换为输出格式
func (a *CaveAggregator) ToMap() map[string]map[string]*types.CaveDailyStats {
	result := make(map[string]map[string]*types.CaveDailyStats)

	characterKeys := make([]string, 0, len(a.caveStats))
	for k := range a.caveStats {
		characterKeys = append(characterKeys, k)
	}
	sort.Strings(characterKeys)

	for _, character := range characterKeys {
		dateMap := make(map[string]*types.CaveDailyStats)

		dateKeys := make([]string, 0, len(a.caveStats[character]))
		for k := range a.caveStats[character] {
			dateKeys = append(dateKeys, k)
		}
		sort.Strings(dateKeys)

		for _, date := range dateKeys {
			dateMap[date] = a.caveStats[character][date]
		}

		result[character] = dateMap
	}

	return result
}

// RecordCount 返回已处理的记录数
func (a *CaveAggregator) RecordCount() int {
	return a.recordCount
}
