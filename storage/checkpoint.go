package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// CheckpointManager 检查点管理器
type CheckpointManager struct {
	stateFilePath string
	checkpoint    string
}

// NewCheckpointManager 创建检查点管理器
func NewCheckpointManager(stateFilePath string) *CheckpointManager {
	return &CheckpointManager{
		stateFilePath: stateFilePath,
	}
}

// Load 加载检查点
func (m *CheckpointManager) Load() string {
	if _, err := os.Stat(m.stateFilePath); os.IsNotExist(err) {
		m.checkpoint = ""
		return ""
	}

	data, err := os.ReadFile(m.stateFilePath)
	if err != nil {
		log.Fatalf("读取状态文件失败: %v", err)
	}

	var state struct {
		Checkpoint string `json:"checkpoint"`
	}
	if err := json.Unmarshal(data, &state); err != nil {
		log.Fatalf("解析状态文件失败: %v", err)
	}

	m.checkpoint = state.Checkpoint
	if m.checkpoint != "" {
		fmt.Printf("从 %s 开始处理日志\n", m.checkpoint)
	} else {
		fmt.Println("从头开始处理日志")
	}

	return m.checkpoint
}

// Save 保存检查点
func (m *CheckpointManager) Save(checkpoint string) {
	m.checkpoint = checkpoint

	state := struct {
		Checkpoint string `json:"checkpoint"`
	}{
		Checkpoint: checkpoint,
	}

	data, err := json.Marshal(state)
	if err != nil {
		log.Fatalf("序列化状态失败: %v", err)
	}

	if err := os.WriteFile(m.stateFilePath, data, 0644); err != nil {
		log.Fatalf("写入状态文件失败: %v", err)
	}
}

// Get 获取当前检查点
func (m *CheckpointManager) Get() string {
	return m.checkpoint
}
