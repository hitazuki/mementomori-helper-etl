package main

import (
	"encoding/json"
	"log"
	"os"
	"sort"
)

// CaveAggregator 洞穴聚合器
type CaveAggregator struct {
	caveStats   map[string]map[string]*CaveDailyStats // character -> date -> stats
	recordCount int
}

// NewCaveAggregator 创建洞穴聚合器
func NewCaveAggregator() *CaveAggregator {
	return &CaveAggregator{
		caveStats: make(map[string]map[string]*CaveDailyStats),
	}
}

// AddRecord 添加洞穴记录
func (a *CaveAggregator) AddRecord(record CaveRecord) {
	character := record.Character
	date := record.Date

	if a.caveStats[character] == nil {
		a.caveStats[character] = make(map[string]*CaveDailyStats)
	}

	if a.caveStats[character][date] == nil {
		a.caveStats[character][date] = &CaveDailyStats{
			Date:    date,
			Records: []CaveRecord{},
		}
	}

	a.caveStats[character][date].Records = append(
		a.caveStats[character][date].Records, record)

	// 更新每日状态
	a.updateDailyStatus(character, date)
	a.recordCount++
}

// updateDailyStatus 更新每日汇总状态
func (a *CaveAggregator) updateDailyStatus(character, date string) {
	daily := a.caveStats[character][date]
	hasStarted := false
	hasFinished := false
	hasError := false

	for _, r := range daily.Records {
		switch r.Status {
		case CaveStatusError:
			hasError = true
		case CaveStatusFinished:
			hasFinished = true
		case CaveStatusStarted:
			hasStarted = true
		}
	}

	// 状态优先级：异常 > 已完成 > 未完成
	if hasError {
		daily.Status = CaveStatusError
	} else if hasFinished {
		daily.Status = CaveStatusFinished
	} else if hasStarted {
		daily.Status = CaveStatusStarted
	}
}

// LoadExistingStats 加载现有统计数据（用于增量合并）
func (a *CaveAggregator) LoadExistingStats(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	var existingStats map[string]map[string]*CaveDailyStats
	if err := json.NewDecoder(file).Decode(&existingStats); err != nil {
		return
	}

	for character, dates := range existingStats {
		if a.caveStats[character] == nil {
			a.caveStats[character] = make(map[string]*CaveDailyStats)
		}
		for date, stats := range dates {
			a.caveStats[character][date] = stats
		}
	}
}

// ToMap 转换为输出格式
func (a *CaveAggregator) ToMap() map[string]map[string]*CaveDailyStats {
	result := make(map[string]map[string]*CaveDailyStats)

	// 按角色名排序
	characterKeys := make([]string, 0, len(a.caveStats))
	for k := range a.caveStats {
		characterKeys = append(characterKeys, k)
	}
	sort.Strings(characterKeys)

	for _, character := range characterKeys {
		dateMap := make(map[string]*CaveDailyStats)

		// 按日期排序
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

// SaveCaveStats 保存洞穴统计结果到文件
func SaveCaveStats(stats map[string]map[string]*CaveDailyStats, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("创建洞穴统计文件失败: %v", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(stats); err != nil {
		log.Printf("写入洞穴统计文件失败: %v", err)
		return
	}

	log.Printf("洞穴统计结果已保存到: %s", filePath)
}
