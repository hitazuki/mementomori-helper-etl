package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// 命令行参数
	outputDir := flag.String("output", "./data", "输出目录路径")
	flag.Parse()

	// 检查日志文件参数
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("用法: mmth_etl [-output <输出目录>] <日志文件路径>")
		os.Exit(1)
	}
	inputLogPath := args[0]

	// 配置
	outputJSONPath := filepath.Join(*outputDir, "diamond_stats.json")
	caveJSONPath := filepath.Join(*outputDir, "cave_stats.json")
	challengeJSONPath := filepath.Join(*outputDir, "challenge_stats.json")
	runeTicketJSONPath := filepath.Join(*outputDir, "rune_ticket_stats.json")
	upgradePanaceaJSONPath := filepath.Join(*outputDir, "upgrade_panacea_stats.json")
	stateFilePath := filepath.Join(*outputDir, "mmth_etl_state.json")

	processor := &LogProcessor{
		outputJSONPath: outputJSONPath,
		stateFilePath:  stateFilePath,
		inputLogPath:   inputLogPath,
	}

	// 加载上次记录的时间戳
	processor.loadCheckpoint()

	// 创建聚合器（不保留详细记录以节省内存）
	// 如果需要保留 records，改为 NewAggregator(true)
	agg := NewAggregator(false)

	// 创建洞穴聚合器
	caveAgg := NewCaveAggregator()

	// 创建挑战聚合器
	challengeAgg := NewChallengeAggregator()

	// 创建物品聚合器
	runeTicketAgg := NewRuneTicketAggregator()
	upgradePanaceaAgg := NewUpgradePanaceaAggregator()

	// 加载已有统计（增量处理）
	if existingStats := loadExistingStats(outputJSONPath); existingStats != nil {
		agg.LoadExistingStats(existingStats)
		fmt.Printf("已加载现有统计数据\n")
	}

	// 加载已有洞穴统计（增量处理）
	caveAgg.LoadExistingStats(caveJSONPath)

	// 加载已有挑战统计（增量处理）
	challengeAgg.LoadExistingStats(challengeJSONPath)

	// 加载已有物品统计（增量处理）
	runeTicketAgg.LoadExistingStats(runeTicketJSONPath)
	upgradePanaceaAgg.LoadExistingStats(upgradePanaceaJSONPath)

	// 流式处理日志（内存友好）
	lastLogTime := processor.processStream(agg, caveAgg, challengeAgg, runeTicketAgg, upgradePanaceaAgg)

	// 检查是否有新记录
	if agg.RecordCount() == 0 && caveAgg.RecordCount() == 0 && challengeAgg.RecordCount() == 0 && runeTicketAgg.RecordCount() == 0 && upgradePanaceaAgg.RecordCount() == 0 {
		fmt.Println("没有新的记录需要处理")
		return
	}

	fmt.Printf("新增 %d 条钻石记录\n", agg.RecordCount())
	fmt.Printf("新增 %d 条洞穴记录\n", caveAgg.RecordCount())
	fmt.Printf("新增 %d 条挑战记录\n", challengeAgg.RecordCount())
	fmt.Printf("新增 %d 条饼干记录\n", runeTicketAgg.RecordCount())
	fmt.Printf("新增 %d 条红水记录\n", upgradePanaceaAgg.RecordCount())

	// 输出统计结果
	if agg.RecordCount() > 0 {
		stats := agg.ToMap()
		SaveStats(stats, outputJSONPath)
	}

	// 输出洞穴统计结果
	if caveAgg.RecordCount() > 0 {
		caveStats := caveAgg.ToMap()
		SaveCaveStats(caveStats, caveJSONPath)
	}

	// 输出挑战统计结果
	if challengeAgg.RecordCount() > 0 {
		challengeStats := challengeAgg.ToMap()
		SaveChallengeStats(challengeStats, challengeJSONPath)
	}

	// 输出饼干统计结果
	if runeTicketAgg.RecordCount() > 0 {
		runeTicketStats := runeTicketAgg.ToMap()
		SaveRuneTicketStats(runeTicketStats, runeTicketJSONPath)
	}

	// 输出红水统计结果
	if upgradePanaceaAgg.RecordCount() > 0 {
		upgradePanaceaStats := upgradePanaceaAgg.ToMap()
		SaveUpgradePanaceaStats(upgradePanaceaStats, upgradePanaceaJSONPath)
	}

	// 更新状态
	if lastLogTime != "" {
		processor.checkpoint = lastLogTime
		processor.saveCheckpoint()
	}

	fmt.Printf("处理完成\n")
}
