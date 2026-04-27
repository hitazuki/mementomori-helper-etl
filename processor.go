package main

import (
	"bufio"
	"log"
	"mmth-etl/aggregator"
	"mmth-etl/i18n"
	"mmth-etl/parser"
	"mmth-etl/storage"
	"mmth-etl/types"
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

		// Update checkpoint for all valid records
		if parsed.Timestamp != "" {
			lastLogTime = parsed.Timestamp
		}

		// Dynamic language detection
		if p.dynamicLang.Enabled {
			scoreAccumulator.AddLine(parsed.Body)

			if lineCount%checkInterval == 0 {
				p.checkLanguageSwitch(scoreAccumulator, switchThreshold, &languageSwitchCount)
			}
		}

		// Step 1: Identify log type first
		logType := parser.IdentifyLogType(parsed.Body)

		// Step 2: Process based on log type
		switch logType {
		case parser.LogTypeDiamond, parser.LogTypeRuneTicket, parser.LogTypeUpgradePanacea:
			// Item change records need source context
			source := lastSourceByCharacter[parsed.Character]
			if source == "" {
				source = "none"
			}
			p.processItemChange(parsed, logType, source, diamondAgg, runeTicketAgg, upgradePanaceaAgg, recordsWriter, &newRecordCount)

		case parser.LogTypeCave:
			// Cave logs: enter/finish are source context, errors clear source
			if types.CaveErrorRegex.MatchString(parsed.Body) || types.SystemErrorRegex.MatchString(parsed.Body) {
				// Cave error logs clear source context
				delete(lastSourceByCharacter, parsed.Character)
			} else {
				// Cave enter/finish logs are valid source context
				lastSourceByCharacter[parsed.Character] = parsed.Body
			}
			caveRecord := parser.ExtractCaveRecord(parsed)
			if caveRecord != nil {
				caveAgg.AddRecord(*caveRecord)
				newRecordCount++
			}

		case parser.LogTypeChallenge:
			challengeRecord := parser.ExtractChallengeRecord(parsed)
			if challengeRecord != nil {
				challengeAgg.AddRecord(*challengeRecord)
				newRecordCount++
			}

		case parser.LogTypeGacha, parser.LogTypeOpen:
			// Gacha/Open logs are valid source context
			// Clean quantity suffix when storing (e.g., "Gacha test 5 times" -> "Gacha test")
			cleanSource := parser.CleanSourceSuffix(parsed.Body, logType)
			lastSourceByCharacter[parsed.Character] = cleanSource

		case parser.LogTypeSystemError:
			// System/error logs clear source context
			delete(lastSourceByCharacter, parsed.Character)

		case parser.LogTypeNameLabel:
			// Name: prefix logs (item change in different language)
			// Do not store as source context, do not clear existing source

		case parser.LogTypeNone:
			// Unknown type - still valid as source context
			lastSourceByCharacter[parsed.Character] = parsed.Body
		}
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

// processItemChange handles item change records (diamond, rune ticket, upgrade panacea).
func (p *LogProcessor) processItemChange(
	parsed parser.ParsedLog,
	logType parser.LogType,
	source string,
	diamondAgg *aggregator.ChangeAggregator,
	runeTicketAgg *aggregator.ChangeAggregator,
	upgradePanaceaAgg *aggregator.ChangeAggregator,
	recordsWriter *storage.RecordsWriter,
	newRecordCount *int,
) {
	record := parser.ExtractChangeRecord(parsed, source, logType)
	if record == nil {
		return
	}

	switch logType {
	case parser.LogTypeDiamond:
		diamondAgg.AddRecord(*record)
		if recordsWriter != nil {
			recordsWriter.AppendRecord("diamond", *record)
		}
		*newRecordCount++

	case parser.LogTypeRuneTicket:
		runeTicketAgg.AddRecord(*record)
		if recordsWriter != nil {
			recordsWriter.AppendRecord("rune_ticket", *record)
		}
		*newRecordCount++

	case parser.LogTypeUpgradePanacea:
		upgradePanaceaAgg.AddRecord(*record)
		if recordsWriter != nil {
			recordsWriter.AppendRecord("upgrade_panacea", *record)
		}
		*newRecordCount++
	}
}
