package aggregator

import (
	"encoding/json"
	"mmth-etl/types"
	"mmth-etl/utils"
	"os"
	"sort"
	"time"
)

// ChangeAggregator 物品变动聚合器（钻石/饼干/红水通用）
type ChangeAggregator struct {
	dailyStats   map[string]map[string]*types.DailyStats
	weeklyStats  map[string]map[string]*types.WeeklyStats
	monthlyStats map[string]map[string]*types.MonthlyStats
	totalStats   map[string]*types.TotalStats
	keepRecords  bool
	recordCount  int
}

// NewChangeAggregator 创建聚合器
func NewChangeAggregator(keepRecords bool) *ChangeAggregator {
	return &ChangeAggregator{
		dailyStats:   make(map[string]map[string]*types.DailyStats),
		weeklyStats:  make(map[string]map[string]*types.WeeklyStats),
		monthlyStats: make(map[string]map[string]*types.MonthlyStats),
		totalStats:   make(map[string]*types.TotalStats),
		keepRecords:  keepRecords,
		recordCount:  0,
	}
}

// AddRecord 添加变动记录
func (a *ChangeAggregator) AddRecord(record types.ChangeRecord) {
	character := record.Character

	// 初始化角色统计
	if a.dailyStats[character] == nil {
		a.dailyStats[character] = make(map[string]*types.DailyStats)
		a.weeklyStats[character] = make(map[string]*types.WeeklyStats)
		a.monthlyStats[character] = make(map[string]*types.MonthlyStats)
		a.totalStats[character] = &types.TotalStats{}
		a.totalStats[character].Sources = make(map[string]types.SourceStats)
	}

	// 解析日期
	date := record.Timestamp[:10]
	t, _ := time.Parse("2006-01-02", date)
	week := utils.GetWeekString(t)
	month := utils.GetMonthString(t)

	// 更新统计
	a.updateDaily(character, date, record)
	a.updateWeekly(character, week, record)
	a.updateMonthly(character, month, record)
	a.updateTotal(character, record)

	a.recordCount++
}

func (a *ChangeAggregator) updateDaily(character, date string, record types.ChangeRecord) {
	if a.dailyStats[character][date] == nil {
		a.dailyStats[character][date] = &types.DailyStats{
			Date:        date,
			ChangeStats: types.ChangeStats{Sources: make(map[string]types.SourceStats)},
		}
	}
	ds := a.dailyStats[character][date]

	if record.Amount > 0 {
		ds.Gain += record.Amount
		src := ds.Sources[record.Source]
		src.Gain += record.Amount
		ds.Sources[record.Source] = src
	} else {
		ds.Consume += -record.Amount
		src := ds.Sources[record.Source]
		src.Consume += -record.Amount
		ds.Sources[record.Source] = src
	}
	ds.NetChange = ds.Gain - ds.Consume

	if a.keepRecords {
		ds.Records = append(ds.Records, record)
	}
}

func (a *ChangeAggregator) updateWeekly(character, week string, record types.ChangeRecord) {
	if a.weeklyStats[character][week] == nil {
		a.weeklyStats[character][week] = &types.WeeklyStats{
			Week:        week,
			ChangeStats: types.ChangeStats{Sources: make(map[string]types.SourceStats)},
		}
	}
	ws := a.weeklyStats[character][week]

	if record.Amount > 0 {
		ws.Gain += record.Amount
		src := ws.Sources[record.Source]
		src.Gain += record.Amount
		ws.Sources[record.Source] = src
	} else {
		ws.Consume += -record.Amount
		src := ws.Sources[record.Source]
		src.Consume += -record.Amount
		ws.Sources[record.Source] = src
	}
	ws.NetChange = ws.Gain - ws.Consume
}

func (a *ChangeAggregator) updateMonthly(character, month string, record types.ChangeRecord) {
	if a.monthlyStats[character][month] == nil {
		a.monthlyStats[character][month] = &types.MonthlyStats{
			Month:       month,
			ChangeStats: types.ChangeStats{Sources: make(map[string]types.SourceStats)},
		}
	}
	ms := a.monthlyStats[character][month]

	if record.Amount > 0 {
		ms.Gain += record.Amount
		src := ms.Sources[record.Source]
		src.Gain += record.Amount
		ms.Sources[record.Source] = src
	} else {
		ms.Consume += -record.Amount
		src := ms.Sources[record.Source]
		src.Consume += -record.Amount
		ms.Sources[record.Source] = src
	}
	ms.NetChange = ms.Gain - ms.Consume
}

func (a *ChangeAggregator) updateTotal(character string, record types.ChangeRecord) {
	ts := a.totalStats[character]

	if record.Amount > 0 {
		ts.Gain += record.Amount
		src := ts.Sources[record.Source]
		src.Gain += record.Amount
		ts.Sources[record.Source] = src
	} else {
		ts.Consume += -record.Amount
		src := ts.Sources[record.Source]
		src.Consume += -record.Amount
		ts.Sources[record.Source] = src
	}
	ts.NetChange = ts.Gain - ts.Consume
}

// LoadExistingStats 加载已有数据
func (a *ChangeAggregator) LoadExistingStats(filePath string) {
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
			a.dailyStats[character] = make(map[string]*types.DailyStats)
			a.weeklyStats[character] = make(map[string]*types.WeeklyStats)
			a.monthlyStats[character] = make(map[string]*types.MonthlyStats)
			a.totalStats[character] = &types.TotalStats{}
			a.totalStats[character].Sources = make(map[string]types.SourceStats)
		}

		// 加载每日统计
		if dailyData, ok := charData["daily"].(map[string]any); ok {
			for date, dayData := range dailyData {
				if dayMap, ok := dayData.(map[string]any); ok {
					ds := &types.DailyStats{
						Date:        date,
						ChangeStats: types.ChangeStats{Sources: make(map[string]types.SourceStats)},
					}
					ds.Gain = getInt(dayMap, "gain")
					ds.Consume = getInt(dayMap, "consume")
					ds.NetChange = getInt(dayMap, "net_change")
					if sources, ok := dayMap["sources"].(map[string]any); ok {
						for srcName, srcData := range sources {
							if srcMap, ok := srcData.(map[string]any); ok {
								ds.Sources[srcName] = types.SourceStats{
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
					ws := &types.WeeklyStats{
						Week:        week,
						ChangeStats: types.ChangeStats{Sources: make(map[string]types.SourceStats)},
					}
					ws.Gain = getInt(weekMap, "gain")
					ws.Consume = getInt(weekMap, "consume")
					ws.NetChange = getInt(weekMap, "net_change")
					a.weeklyStats[character][week] = ws
				}
			}
		}

		// 加载每月统计
		if monthlyData, ok := charData["monthly"].(map[string]any); ok {
			for month, monthData := range monthlyData {
				if monthMap, ok := monthData.(map[string]any); ok {
					ms := &types.MonthlyStats{
						Month:       month,
						ChangeStats: types.ChangeStats{Sources: make(map[string]types.SourceStats)},
					}
					ms.Gain = getInt(monthMap, "gain")
					ms.Consume = getInt(monthMap, "consume")
					ms.NetChange = getInt(monthMap, "net_change")
					a.monthlyStats[character][month] = ms
				}
			}
		}

		// 加载总计
		if totalData, ok := charData["total"].(map[string]any); ok {
			ts := &types.TotalStats{
				ChangeStats: types.ChangeStats{Sources: make(map[string]types.SourceStats)},
			}
			ts.Gain = getInt(totalData, "gain")
			ts.Consume = getInt(totalData, "consume")
			ts.NetChange = getInt(totalData, "net_change")
			if sources, ok := totalData["sources"].(map[string]any); ok {
				for srcName, srcData := range sources {
					if srcMap, ok := srcData.(map[string]any); ok {
						ts.Sources[srcName] = types.SourceStats{
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
func (a *ChangeAggregator) ToMap() map[string]map[string]any {
	stats := make(map[string]map[string]any)

	characterKeys := make([]string, 0, len(a.dailyStats))
	for k := range a.dailyStats {
		characterKeys = append(characterKeys, k)
	}
	sort.Strings(characterKeys)

	for _, character := range characterKeys {
		characterData := make(map[string]any)

		// 每日统计
		dailyData := make(map[string]*types.DailyStats)
		dateKeys := make([]string, 0, len(a.dailyStats[character]))
		for k := range a.dailyStats[character] {
			dateKeys = append(dateKeys, k)
		}
		sort.Strings(dateKeys)
		for _, k := range dateKeys {
			dailyData[k] = a.dailyStats[character][k]
		}
		characterData["daily"] = dailyData

		// 每周统计
		weeklyData := make(map[string]*types.WeeklyStats)
		weekKeys := make([]string, 0, len(a.weeklyStats[character]))
		for k := range a.weeklyStats[character] {
			weekKeys = append(weekKeys, k)
		}
		sort.Strings(weekKeys)
		for _, k := range weekKeys {
			weeklyData[k] = a.weeklyStats[character][k]
		}
		characterData["weekly"] = weeklyData

		// 每月统计
		monthlyData := make(map[string]*types.MonthlyStats)
		monthKeys := make([]string, 0, len(a.monthlyStats[character]))
		for k := range a.monthlyStats[character] {
			monthKeys = append(monthKeys, k)
		}
		sort.Strings(monthKeys)
		for _, k := range monthKeys {
			monthlyData[k] = a.monthlyStats[character][k]
		}
		characterData["monthly"] = monthlyData

		// 总计
		characterData["total"] = a.totalStats[character]

		stats[character] = characterData
	}

	return stats
}

// RecordCount 返回已处理的记录数
func (a *ChangeAggregator) RecordCount() int {
	return a.recordCount
}

func getInt(m map[string]any, key string) int {
	switch v := m[key].(type) {
	case int:
		return v
	case float64:
		return int(v)
	}
	return 0
}
