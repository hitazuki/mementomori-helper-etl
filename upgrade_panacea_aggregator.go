package main

// UpgradePanaceaAggregator 红水聚合器
type UpgradePanaceaAggregator = ItemAggregator

// NewUpgradePanaceaAggregator 创建红水聚合器
func NewUpgradePanaceaAggregator() *UpgradePanaceaAggregator {
	return NewItemAggregator()
}

// SaveUpgradePanaceaStats 保存红水统计结果
func SaveUpgradePanaceaStats(stats map[string]map[string]any, filePath string) {
	SaveItemStats(stats, filePath)
}
