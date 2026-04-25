package main

import (
	"bufio"
	"log"
	"mmth-etl/aggregator"
	"mmth-etl/i18n"
	"mmth-etl/parser"
	"mmth-etl/storage"
	"os"
)

// DynamicLanguageConfig holds configuration for dynamic language detection.
type DynamicLanguageConfig struct {
	Enabled         bool // Enable dynamic language detection
	WindowSize      int  // Number of lines to analyze for language detection
	SwitchThreshold int  // Minimum score difference to trigger language switch
}

// LogProcessor handles log file processing.
type LogProcessor struct {
	inputLogPath string
	checkpoint   string
	i18nMgr      *i18n.Manager
	detector     *i18n.Detector
	dynamicLang  DynamicLanguageConfig
}

// NewLogProcessor creates a new log processor.
func NewLogProcessor(inputLogPath, checkpoint string, i18nMgr *i18n.Manager, dynamicCfg DynamicLanguageConfig) *LogProcessor {
	return &LogProcessor{
		inputLogPath: inputLogPath,
		checkpoint:   checkpoint,
		i18nMgr:      i18nMgr,
		detector:     i18n.NewDetector(),
		dynamicLang:  dynamicCfg,
	}
}

// Process processes the log file and aggregates records.
func (p *LogProcessor) Process(
	diamondAgg *aggregator.ChangeAggregator,
	caveAgg *aggregator.CaveAggregator,
	challengeAgg *aggregator.ChallengeAggregator,
	runeTicketAgg *aggregator.ChangeAggregator,
	upgradePanaceaAgg *aggregator.ChangeAggregator,
	recordsWriter *storage.RecordsWriter,
) string {
	file, err := os.Open(p.inputLogPath)
	if err != nil {
		log.Fatalf("打开日志文件失败: %v", err)
	}
	defer file.Close()

	startPos, err := findStartPosition(file, p.checkpoint)
	if err != nil {
		log.Printf("定位起始位置失败: %v，从头开始处理", err)
		startPos = 0
	}
	if startPos < 0 {
		log.Println("无新记录需要处理")
		return p.checkpoint
	}

	_, err = file.Seek(startPos, 0)
	if err != nil {
		log.Fatalf("文件定位失败: %v", err)
	}

	// Source context per character
	lastSourceByCharacter := make(map[string]string)

	// Dynamic language detection
	windowSize := p.dynamicLang.WindowSize
	if windowSize <= 0 {
		windowSize = 100
	}
	switchThreshold := p.dynamicLang.SwitchThreshold
	if switchThreshold <= 0 {
		switchThreshold = 5
	}
	var scoreAccumulator *i18n.ScoreAccumulator
	if p.dynamicLang.Enabled {
		scoreAccumulator = i18n.NewScoreAccumulator(p.detector, windowSize)
	}

	// Read with buffer
	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 256*1024)
	scanner.Buffer(buf, 1024*1024)

	var lastLogTime string
	lineCount := 0
	newRecordCount := 0
	languageSwitchCount := 0
	checkInterval := max(windowSize/2, 1)

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		parsed := parser.ParseLogLine(line)
		if !parsed.IsValid {
			continue
		}

		// Skip already processed records
		if p.checkpoint != "" && parsed.Timestamp != "" {
			if parsed.Timestamp <= p.checkpoint {
				continue
			}
		}

		// Dynamic language detection
		if p.dynamicLang.Enabled {
			scoreAccumulator.AddLine(parsed.Body)

			if lineCount%checkInterval == 0 {
				p.checkLanguageSwitch(scoreAccumulator, switchThreshold, &languageSwitchCount)
			}
		}

		// Update source context
		if parser.IsValidSource(parsed.Body) {
			lastSourceByCharacter[parsed.Character] = parsed.Body
		} else if parser.ShouldClearSource(parsed.Body) {
			delete(lastSourceByCharacter, parsed.Character)
		}

		// Dispatch by log type
		p.processLogType(parsed, lastSourceByCharacter, diamondAgg, caveAgg, challengeAgg, runeTicketAgg, upgradePanaceaAgg, recordsWriter, &newRecordCount, &lastLogTime)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("读取日志文件时出错: %v", err)
	}

	log.Printf("扫描了 %d 行，新增 %d 条记录", lineCount, newRecordCount)
	if p.dynamicLang.Enabled && languageSwitchCount > 0 {
		log.Printf("动态语言切换 %d 次", languageSwitchCount)
	}

	return lastLogTime
}

// checkLanguageSwitch checks if language should be switched based on accumulated scores.
func (p *LogProcessor) checkLanguageSwitch(scoreAccumulator *i18n.ScoreAccumulator, switchThreshold int, switchCount *int) {
	scores := scoreAccumulator.GetScores()
	currentLang := p.i18nMgr.CurrentLanguage()

	var maxLang i18n.Language
	maxScore := 0
	for lang, score := range scores {
		if score > maxScore {
			maxScore = score
			maxLang = lang
		}
	}

	if maxLang != "" && maxLang != currentLang {
		currentScore := scores[currentLang]
		if maxScore-currentScore >= switchThreshold {
			p.i18nMgr.SetLanguage(maxLang)
			*switchCount++
			log.Printf("语言切换: %s -> %s (得分: %d vs %d)", currentLang, maxLang, maxScore, currentScore)
		}
	}
}

// processLogType processes a single log line based on its type.
func (p *LogProcessor) processLogType(
	parsed parser.ParsedLog,
	lastSourceByCharacter map[string]string,
	diamondAgg *aggregator.ChangeAggregator,
	caveAgg *aggregator.CaveAggregator,
	challengeAgg *aggregator.ChallengeAggregator,
	runeTicketAgg *aggregator.ChangeAggregator,
	upgradePanaceaAgg *aggregator.ChangeAggregator,
	recordsWriter *storage.RecordsWriter,
	newRecordCount *int,
	lastLogTime *string,
) {
	logType := parser.IdentifyLogType(parsed.Body)

	switch logType {
	case parser.LogTypeDiamond:
		source := lastSourceByCharacter[parsed.Character]
		if source == "" {
			source = "none"
		}
		record := parser.ExtractChangeRecord(parsed, source, logType)
		if record != nil {
			diamondAgg.AddRecord(*record)
			if recordsWriter != nil {
				recordsWriter.AppendRecord("diamond", *record)
			}
			*newRecordCount++
			*lastLogTime = record.Timestamp
		}

	case parser.LogTypeCave:
		caveRecord := parser.ExtractCaveRecord(parsed)
		if caveRecord != nil {
			caveAgg.AddRecord(*caveRecord)
		}

	case parser.LogTypeChallenge:
		challengeRecord := parser.ExtractChallengeRecord(parsed)
		if challengeRecord != nil {
			challengeAgg.AddRecord(*challengeRecord)
		}

	case parser.LogTypeRuneTicket:
		source := lastSourceByCharacter[parsed.Character]
		if source == "" {
			source = "none"
		}
		record := parser.ExtractChangeRecord(parsed, source, logType)
		if record != nil {
			runeTicketAgg.AddRecord(*record)
			if recordsWriter != nil {
				recordsWriter.AppendRecord("rune_ticket", *record)
			}
		}

	case parser.LogTypeUpgradePanacea:
		source := lastSourceByCharacter[parsed.Character]
		if source == "" {
			source = "none"
		}
		record := parser.ExtractChangeRecord(parsed, source, logType)
		if record != nil {
			upgradePanaceaAgg.AddRecord(*record)
			if recordsWriter != nil {
				recordsWriter.AppendRecord("upgrade_panacea", *record)
			}
		}
	}
}
