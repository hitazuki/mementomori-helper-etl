package main

// RuneTicketAggregator 饼干聚合器
type RuneTicketAggregator = ItemAggregator

// NewRuneTicketAggregator 创建饼干聚合器
func NewRuneTicketAggregator() *RuneTicketAggregator {
	return NewItemAggregator()
}

// SaveRuneTicketStats 保存饼干统计结果
func SaveRuneTicketStats(stats map[string]map[string]any, filePath string) {
	SaveItemStats(stats, filePath)
}
