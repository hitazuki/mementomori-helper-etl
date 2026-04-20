package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// ItemAggregator 物品统计聚合器（支持增量更新）
type ItemAggregator struct {
	dailyStats   map[string]map[string]*ItemDailyStats   // character -> date -> stats
	weeklyStats  map[string]map[string]*ItemWeeklyStats  // character -> week -> stats
	monthlyStats map[string]map[string]*ItemMonthlyStats // character -> month -> stats
	totalStats   map[string]*ItemTotalStats              // character -> total
	recordCount  int
}

// NewItemAggregator 创建物品聚合器
func NewItemAggregator() *ItemAggregator {
	return &ItemAggregator{
		dailyStats:   make(map[string]map[string]*ItemDailyStats),
		weeklyStats:  make(map[string]map[string]*ItemWeeklyStats),
		monthlyStats: make(map[string]map[string]*ItemMonthlyStats),
		totalStats:   make(map[string]*ItemTotalStats),
		recordCount:  0,
	}
}

// AddRecord 增量添加记录
func (a *ItemAggregator) AddRecord(record ItemRecord) {
	character := record.Character

	// 初始化角色统计
	if a.dailyStats[character] == nil {
		a.dailyStats[character] = make(map[string]*ItemDailyStats)
		a.weeklyStats[character] = make(map[string]*ItemWeeklyStats)
		a.monthlyStats[character] = make(map[string]*ItemMonthlyStats)
		a.totalStats[character] = &ItemTotalStats{Sources: make(map[string]SourceStats)}
	}

	// 解析日期
	date := record.Timestamp[:10]
	t, _ := time.Parse("2006-01-02", date)
	week := getWeekString(t)
	month := getMonthString(t)

	// 更新每日统计
	if a.dailyStats[character][date] == nil {
		a.dailyStats[character][date] = &ItemDailyStats{
			Date:    date,
			Sources: make(map[string]SourceStats),
		}
	}
	a.updateDailyStats(a.dailyStats[character][date], record)

	// 更新每周统计
	if a.weeklyStats[character][week] == nil {
		a.weeklyStats[character][week] = &ItemWeeklyStats{
			Week:    week,
			Sources: make(map[string]SourceStats),
		}
	}
	a.updateWeeklyStats(a.weeklyStats[character][week], record)

	// 更新每月统计
	if a.monthlyStats[character][month] == nil {
		a.monthlyStats[character][month] = &ItemMonthlyStats{
			Month:   month,
			Sources: make(map[string]SourceStats),
		}
	}
	a.updateMonthlyStats(a.monthlyStats[character][month], record)

	// 更新总计
	a.updateTotalStats(a.totalStats[character], record)

	a.recordCount++
}

// updateDailyStats 更新每日统计
func (a *ItemAggregator) updateDailyStats(ds *ItemDailyStats, record ItemRecord) {
	if record.Amount > 0 {
		ds.Gain += record.Amount
		src := ds.Sources[record.Source]
		src.Gain += record.Amount
		ds.Sources[record.Source] = src
	} else {
		amount := -record.Amount
		ds.Consume += amount
		src := ds.Sources[record.Source]
		src.Consume += amount
		ds.Sources[record.Source] = src
	}
	ds.NetChange = ds.Gain - ds.Consume
}

// updateWeeklyStats 更新每周统计
func (a *ItemAggregator) updateWeeklyStats(ws *ItemWeeklyStats, record ItemRecord) {
	if record.Amount > 0 {
		ws.Gain += record.Amount
		src := ws.Sources[record.Source]
		src.Gain += record.Amount
		ws.Sources[record.Source] = src
	} else {
		amount := -record.Amount
		ws.Consume += amount
		src := ws.Sources[record.Source]
		src.Consume += amount
		ws.Sources[record.Source] = src
	}
	ws.NetChange = ws.Gain - ws.Consume
}

// updateMonthlyStats 更新每月统计
func (a *ItemAggregator) updateMonthlyStats(ms *ItemMonthlyStats, record ItemRecord) {
	if record.Amount > 0 {
		ms.Gain += record.Amount
		src := ms.Sources[record.Source]
		src.Gain += record.Amount
		ms.Sources[record.Source] = src
	} else {
		amount := -record.Amount
		ms.Consume += amount
		src := ms.Sources[record.Source]
		src.Consume += amount
		ms.Sources[record.Source] = src
	}
	ms.NetChange = ms.Gain - ms.Consume
}

// updateTotalStats 更新总计
func (a *ItemAggregator) updateTotalStats(ts *ItemTotalStats, record ItemRecord) {
	if record.Amount > 0 {
		ts.Gain += record.Amount
		src := ts.Sources[record.Source]
		src.Gain += record.Amount
		ts.Sources[record.Source] = src
	} else {
		amount := -record.Amount
		ts.Consume += amount
		src := ts.Sources[record.Source]
		src.Consume += amount
		ts.Sources[record.Source] = src
	}
	ts.NetChange = ts.Gain - ts.Consume
}

// LoadExistingStats 加载现有统计数据
func (a *ItemAggregator) LoadExistingStats(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	var existingStats map[string]map[string]any
	if err := json.NewDecoder(file).Decode(&existingStats); err != nil {
		return
	}

	for character, charData := range existingStats {
		if a.dailyStats[character] == nil {
			a.dailyStats[character] = make(map[string]*ItemDailyStats)
			a.weeklyStats[character] = make(map[string]*ItemWeeklyStats)
			a.monthlyStats[character] = make(map[string]*ItemMonthlyStats)
			a.totalStats[character] = &ItemTotalStats{Sources: make(map[string]SourceStats)}
		}

		// 加载每日统计
		if dailyData, ok := charData["daily"].(map[string]any); ok {
			for date, dayData := range dailyData {
				if dayMap, ok := dayData.(map[string]any); ok {
					ds := &ItemDailyStats{
						Date:      date,
						Gain:      getInt(dayMap, "gain"),
						Consume:   getInt(dayMap, "consume"),
						NetChange: getInt(dayMap, "net_change"),
						Sources:   make(map[string]SourceStats),
					}
					if sources, ok := dayMap["sources"].(map[string]any); ok {
						for srcName, srcData := range sources {
							if srcMap, ok := srcData.(map[string]any); ok {
								ds.Sources[srcName] = SourceStats{
									Gain:    getInt(srcMap, "gain"),
									Consume: getInt(srcMap, "consume"),
								}
							}
						}
					}
					a.dailyStats[character][date] = ds
				}
			}
		}

		// 加载每周统计
		if weeklyData, ok := charData["weekly"].(map[string]any); ok {
			for week, weekData := range weeklyData {
				if weekMap, ok := weekData.(map[string]any); ok {
					ws := &ItemWeeklyStats{
						Week:      week,
						Gain:      getInt(weekMap, "gain"),
						Consume:   getInt(weekMap, "consume"),
						NetChange: getInt(weekMap, "net_change"),
						Sources:   make(map[string]SourceStats),
					}
					if sources, ok := weekMap["sources"].(map[string]any); ok {
						for srcName, srcData := range sources {
							if srcMap, ok := srcData.(map[string]any); ok {
								ws.Sources[srcName] = SourceStats{
									Gain:    getInt(srcMap, "gain"),
									Consume: getInt(srcMap, "consume"),
								}
							}
						}
					}
					a.weeklyStats[character][week] = ws
				}
			}
		}

		// 加载每月统计
		if monthlyData, ok := charData["monthly"].(map[string]any); ok {
			for month, monthData := range monthlyData {
				if monthMap, ok := monthData.(map[string]any); ok {
					ms := &ItemMonthlyStats{
						Month:     month,
						Gain:      getInt(monthMap, "gain"),
						Consume:   getInt(monthMap, "consume"),
						NetChange: getInt(monthMap, "net_change"),
						Sources:   make(map[string]SourceStats),
					}
					if sources, ok := monthMap["sources"].(map[string]any); ok {
						for srcName, srcData := range sources {
							if srcMap, ok := srcData.(map[string]any); ok {
								ms.Sources[srcName] = SourceStats{
									Gain:    getInt(srcMap, "gain"),
									Consume: getInt(srcMap, "consume"),
								}
							}
						}
					}
					a.monthlyStats[character][month] = ms
				}
			}
		}

		// 加载总计
		if totalData, ok := charData["total"].(map[string]any); ok {
			ts := &ItemTotalStats{
				Gain:      getInt(totalData, "gain"),
				Consume:   getInt(totalData, "consume"),
				NetChange: getInt(totalData, "net_change"),
				Sources:   make(map[string]SourceStats),
			}
			if sources, ok := totalData["sources"].(map[string]any); ok {
				for srcName, srcData := range sources {
					if srcMap, ok := srcData.(map[string]any); ok {
						ts.Sources[srcName] = SourceStats{
							Gain:    getInt(srcMap, "gain"),
							Consume: getInt(srcMap, "consume"),
						}
					}
				}
			}
			a.totalStats[character] = ts
		}
	}
}

// ToMap 转换为输出格式
func (a *ItemAggregator) ToMap() map[string]map[string]any {
	stats := make(map[string]map[string]any)

	characterKeys := make([]string, 0, len(a.dailyStats))
	for k := range a.dailyStats {
		characterKeys = append(characterKeys, k)
	}
	sort.Strings(characterKeys)

	for _, character := range characterKeys {
		characterData := make(map[string]any)

		// 保存每日统计
		dailyData := make(map[string]*ItemDailyStats)
		var dateKeys []string
		for k := range a.dailyStats[character] {
			dateKeys = append(dateKeys, k)
		}
		sort.Strings(dateKeys)
		for _, k := range dateKeys {
			dailyData[k] = a.dailyStats[character][k]
		}
		characterData["daily"] = dailyData

		// 保存每周统计
		weeklyData := make(map[string]*ItemWeeklyStats)
		var weekKeys []string
		for k := range a.weeklyStats[character] {
			weekKeys = append(weekKeys, k)
		}
		sort.Strings(weekKeys)
		for _, k := range weekKeys {
			weeklyData[k] = a.weeklyStats[character][k]
		}
		characterData["weekly"] = weeklyData

		// 保存每月统计
		monthlyData := make(map[string]*ItemMonthlyStats)
		var monthKeys []string
		for k := range a.monthlyStats[character] {
			monthKeys = append(monthKeys, k)
		}
		sort.Strings(monthKeys)
		for _, k := range monthKeys {
			monthlyData[k] = a.monthlyStats[character][k]
		}
		characterData["monthly"] = monthlyData

		// 保存总统计
		characterData["total"] = a.totalStats[character]

		stats[character] = characterData
	}

	return stats
}

// RecordCount 返回已处理的记录数
func (a *ItemAggregator) RecordCount() int {
	return a.recordCount
}

// SaveItemStats 保存物品统计结果到文件
func SaveItemStats(stats map[string]map[string]any, filePath string) {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("创建统计文件失败: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(stats); err != nil {
		log.Fatalf("写入统计文件失败: %v", err)
	}

	fmt.Printf("物品统计结果已保存到: %s\n", filePath)
}
