package main

import (
	"flag"
	"fmt"
	"log"
	"mmth-etl/aggregator"
	"mmth-etl/storage"
	"os"
	"path/filepath"
)

// 类型别名，语义化
type DiamondAggregator = aggregator.ChangeAggregator
type RuneTicketAggregator = aggregator.ChangeAggregator
type UpgradePanaceaAggregator = aggregator.ChangeAggregator

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

	// 配置路径
	diamondJSONPath := filepath.Join(*outputDir, "diamond_stats.json")
	caveJSONPath := filepath.Join(*outputDir, "cave_stats.json")
	challengeJSONPath := filepath.Join(*outputDir, "challenge_stats.json")
	runeTicketJSONPath := filepath.Join(*outputDir, "rune_ticket_stats.json")
	upgradePanaceaJSONPath := filepath.Join(*outputDir, "upgrade_panacea_stats.json")
	stateFilePath := filepath.Join(*outputDir, "mmth_etl_state.json")

	// 加载检查点
	checkpointMgr := storage.NewCheckpointManager(stateFilePath)
	checkpoint := checkpointMgr.Load()

	// 创建处理器
	processor := NewLogProcessor(inputLogPath, checkpoint)

	// 创建聚合器
	diamondAgg := aggregator.NewChangeAggregator(false)
	caveAgg := aggregator.NewCaveAggregator()
	challengeAgg := aggregator.NewChallengeAggregator()
	runeTicketAgg := aggregator.NewChangeAggregator(false)
	upgradePanaceaAgg := aggregator.NewChangeAggregator(false)

	// 加载已有统计（增量处理）
	diamondAgg.LoadExistingStats(diamondJSONPath)
	caveAgg.LoadExistingStats(caveJSONPath)
	challengeAgg.LoadExistingStats(challengeJSONPath)
	runeTicketAgg.LoadExistingStats(runeTicketJSONPath)
	upgradePanaceaAgg.LoadExistingStats(upgradePanaceaJSONPath)

	log.Println("已加载现有统计数据")

	// 流式处理日志
	lastLogTime := processor.Process(diamondAgg, caveAgg, challengeAgg, runeTicketAgg, upgradePanaceaAgg)

	// 检查是否有新记录
	if diamondAgg.RecordCount() == 0 && caveAgg.RecordCount() == 0 && challengeAgg.RecordCount() == 0 && runeTicketAgg.RecordCount() == 0 && upgradePanaceaAgg.RecordCount() == 0 {
		fmt.Println("没有新的记录需要处理")
		return
	}

	fmt.Printf("新增 %d 条钻石记录\n", diamondAgg.RecordCount())
	fmt.Printf("新增 %d 条洞穴记录\n", caveAgg.RecordCount())
	fmt.Printf("新增 %d 条挑战记录\n", challengeAgg.RecordCount())
	fmt.Printf("新增 %d 条饼干记录\n", runeTicketAgg.RecordCount())
	fmt.Printf("新增 %d 条红水记录\n", upgradePanaceaAgg.RecordCount())

	// 输出统计结果
	if diamondAgg.RecordCount() > 0 {
		storage.SaveStats(diamondAgg.ToMap(), diamondJSONPath)
	}

	if caveAgg.RecordCount() > 0 {
		SaveCaveStats(caveAgg.ToMap(), caveJSONPath)
	}

	if challengeAgg.RecordCount() > 0 {
		SaveChallengeStats(challengeAgg.ToMap(), challengeJSONPath)
	}

	if runeTicketAgg.RecordCount() > 0 {
		storage.SaveStats(runeTicketAgg.ToMap(), runeTicketJSONPath)
	}

	if upgradePanaceaAgg.RecordCount() > 0 {
		storage.SaveStats(upgradePanaceaAgg.ToMap(), upgradePanaceaJSONPath)
	}

	// 更新检查点
	if lastLogTime != "" {
		checkpointMgr.Save(lastLogTime)
	}

	fmt.Printf("处理完成\n")
}
