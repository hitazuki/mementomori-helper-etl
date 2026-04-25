package main

import (
	"mmth-etl/storage"
	"mmth-etl/types"
)

// SaveCaveStats saves cave statistics to file.
func SaveCaveStats(stats map[string]map[string]*types.CaveDailyStats, filePath string) {
	storage.SaveStats(convertCaveStats(stats), filePath)
}

// SaveChallengeStats saves challenge statistics to file.
func SaveChallengeStats(stats map[string]*types.ChallengeStats, filePath string) {
	storage.SaveStats(convertChallengeStats(stats), filePath)
}

func convertCaveStats(stats map[string]map[string]*types.CaveDailyStats) map[string]map[string]any {
	result := make(map[string]map[string]any)
	for k, v := range stats {
		result[k] = make(map[string]any)
		for k2, v2 := range v {
			result[k][k2] = v2
		}
	}
	return result
}

func convertChallengeStats(stats map[string]*types.ChallengeStats) map[string]map[string]any {
	result := make(map[string]map[string]any)
	for k, v := range stats {
		result[k] = map[string]any{
			"quest":  v.Quest,
			"towers": v.Towers,
		}
	}
	return result
}
