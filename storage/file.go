package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// SaveStats 保存统计结果到文件
func SaveStats(stats map[string]map[string]any, filePath string) {
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

	fmt.Printf("统计结果已保存到: %s\n", filePath)
}

// LoadExistingStats 从文件加载现有统计数据
func LoadExistingStats(filePath string) map[string]map[string]any {
	file, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer file.Close()

	var stats map[string]map[string]any
	if err := json.NewDecoder(file).Decode(&stats); err != nil {
		return nil
	}
	return stats
}

// GetInt 从 map 中获取整数
func GetInt(m map[string]any, key string) int {
	switch v := m[key].(type) {
	case int:
		return v
	case float64:
		return int(v)
	}
	return 0
}
