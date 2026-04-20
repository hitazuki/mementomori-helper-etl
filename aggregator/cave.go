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

	if a.caveStats[character][date] == nil {
		a.caveStats[character][date] = &types.CaveDailyStats{
			Date:    date,
			Records: []types.CaveRecord{},
		}
	}

	a.caveStats[character][date].Records = append(
		a.caveStats[character][date].Records, record)

	a.updateDailyStatus(character, date)
	a.recordCount++
}

func (a *CaveAggregator) updateDailyStatus(character, date string) {
	daily := a.caveStats[character][date]
	hasStarted := false
	hasFinished := false
	hasError := false

	for _, r := range daily.Records {
		switch r.Status {
		case types.CaveStatusError:
			hasError = true
		case types.CaveStatusFinished:
			hasFinished = true
		case types.CaveStatusStarted:
			hasStarted = true
		}
	}

	// 状态优先级：异常 > 已完成 > 未完成
	if hasError {
		daily.Status = types.CaveStatusError
	} else if hasFinished {
		daily.Status = types.CaveStatusFinished
	} else if hasStarted {
		daily.Status = types.CaveStatusStarted
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
