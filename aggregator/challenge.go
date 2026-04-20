package aggregator

import (
	"encoding/json"
	"mmth-etl/types"
	"os"
	"sort"
)

// ChallengeAggregator 挑战聚合器
type ChallengeAggregator struct {
	stats       map[string]*types.ChallengeStats // character -> stats
	recordCount int
}

// NewChallengeAggregator 创建挑战聚合器
func NewChallengeAggregator() *ChallengeAggregator {
	return &ChallengeAggregator{
		stats: make(map[string]*types.ChallengeStats),
	}
}

func (a *ChallengeAggregator) getOrCreateCharacterStats(character string) *types.ChallengeStats {
	if a.stats[character] == nil {
		a.stats[character] = &types.ChallengeStats{
			Quest:  make(map[string]*types.ChallengeLevelStats),
			Towers: make(map[types.TowerType]map[string]*types.ChallengeLevelStats),
		}
	}
	return a.stats[character]
}

func (a *ChallengeAggregator) getOrCreateTowerStats(charStats *types.ChallengeStats, towerType types.TowerType) map[string]*types.ChallengeLevelStats {
	if charStats.Towers[towerType] == nil {
		charStats.Towers[towerType] = make(map[string]*types.ChallengeLevelStats)
	}
	return charStats.Towers[towerType]
}

func (a *ChallengeAggregator) updateLevelStats(levelMap map[string]*types.ChallengeLevelStats, level string, status types.ChallengeStatus, timestamp string) {
	if levelMap[level] == nil {
		levelMap[level] = &types.ChallengeLevelStats{Level: level, Success: false}
	}
	levelMap[level].Attempts++
	levelMap[level].LastTime = timestamp
	if status == types.ChallengeStatusSuccess {
		levelMap[level].Success = true
	}
}

// AddRecord 添加挑战记录
func (a *ChallengeAggregator) AddRecord(record types.ChallengeRecord) {
	charStats := a.getOrCreateCharacterStats(record.Character)

	switch record.Type {
	case types.ChallengeTypeQuest:
		a.updateLevelStats(charStats.Quest, record.Level, record.Status, record.Timestamp)
	case types.ChallengeTypeTower:
		towerStats := a.getOrCreateTowerStats(charStats, record.TowerType)
		a.updateLevelStats(towerStats, record.Level, record.Status, record.Timestamp)
	}
	a.recordCount++
}

// LoadExistingStats 加载现有统计数据
func (a *ChallengeAggregator) LoadExistingStats(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	var existingStats map[string]*types.ChallengeStats
	if err := json.NewDecoder(file).Decode(&existingStats); err != nil {
		return
	}

	for character, stats := range existingStats {
		if a.stats[character] == nil {
			a.stats[character] = &types.ChallengeStats{
				Quest:  make(map[string]*types.ChallengeLevelStats),
				Towers: make(map[types.TowerType]map[string]*types.ChallengeLevelStats),
			}
		}

		// 合并主线统计
		for level, levelStats := range stats.Quest {
			a.stats[character].Quest[level] = levelStats
		}

		// 合并塔统计
		for towerType, towerStats := range stats.Towers {
			if a.stats[character].Towers[towerType] == nil {
				a.stats[character].Towers[towerType] = make(map[string]*types.ChallengeLevelStats)
			}
			for level, levelStats := range towerStats {
				a.stats[character].Towers[towerType][level] = levelStats
			}
		}
	}
}

// ToMap 转换为输出格式
func (a *ChallengeAggregator) ToMap() map[string]*types.ChallengeStats {
	result := make(map[string]*types.ChallengeStats)

	characterKeys := make([]string, 0, len(a.stats))
	for k := range a.stats {
		characterKeys = append(characterKeys, k)
	}
	sort.Strings(characterKeys)

	for _, character := range characterKeys {
		result[character] = a.stats[character]
	}

	return result
}

// RecordCount 返回已处理的记录数
func (a *ChallengeAggregator) RecordCount() int {
	return a.recordCount
}
