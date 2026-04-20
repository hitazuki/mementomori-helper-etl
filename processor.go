package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mmth-etl/aggregator"
	"mmth-etl/parser"
	"mmth-etl/storage"
	"mmth-etl/types"
	"os"
	"regexp"
	"strings"
)

// LogProcessor 日志处理器
type LogProcessor struct {
	inputLogPath string
	checkpoint   string
}

// NewLogProcessor 创建日志处理器
func NewLogProcessor(inputLogPath, checkpoint string) *LogProcessor {
	return &LogProcessor{
		inputLogPath: inputLogPath,
		checkpoint:   checkpoint,
	}
}

// Process 流式处理日志
func (p *LogProcessor) Process(
	diamondAgg *aggregator.ChangeAggregator,
	caveAgg *aggregator.CaveAggregator,
	challengeAgg *aggregator.ChallengeAggregator,
	runeTicketAgg *aggregator.ChangeAggregator,
	upgradePanaceaAgg *aggregator.ChangeAggregator,
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

	// 来源上下文（按角色隔离）
	lastSourceByCharacter := make(map[string]string)

	// 分块读取
	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 256*1024)
	scanner.Buffer(buf, 1024*1024)

	var lastLogTime string
	lineCount := 0
	newRecordCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		parsed := parser.ParseLogLine(line)
		if !parsed.IsValid {
			continue
		}

		// 时间戳二次校验：跳过已处理的记录
			// 使用 <= 跳过同秒记录，确保不重复计入
			if p.checkpoint != "" && parsed.Timestamp != "" {
				if parsed.Timestamp <= p.checkpoint {
					continue
				}
			}

		// 更新来源上下文
		if parser.IsValidSource(parsed.Body) {
			lastSourceByCharacter[parsed.Character] = parsed.Body
		}

		// 识别日志类型并分发
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
				newRecordCount++
				lastLogTime = record.Timestamp
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
			}

		case parser.LogTypeUpgradePanacea:
			source := lastSourceByCharacter[parsed.Character]
			if source == "" {
				source = "none"
			}
			record := parser.ExtractChangeRecord(parsed, source, logType)
			if record != nil {
				upgradePanaceaAgg.AddRecord(*record)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("读取日志文件时出错: %v", err)
	}

	log.Printf("扫描了 %d 行，新增 %d 条记录", lineCount, newRecordCount)

	return lastLogTime
}

// findStartPosition 二分查找定位第一个时间戳 > lastLogTime 的字节位置
func findStartPosition(file *os.File, lastLogTime string) (int64, error) {
	if lastLogTime == "" {
		return 0, nil
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return 0, err
	}
	fileSize := fileInfo.Size()
	if fileSize == 0 {
		return 0, nil
	}

	// 快速边界检查：读取文件末尾的时间戳
	// 找到最后一个完整行的起始位置
	lastLineStart := findLastLineStart(file, fileSize)
	lastTimeInFile, err := readTimestampAt(file, lastLineStart)

	// 文件最后时间 <= checkpoint 表示全部已处理
	if err == nil && lastTimeInFile <= lastLogTime {
		return -1, nil
	}

	firstTimeInFile, err := readTimestampAt(file, 0)
	if err == nil && firstTimeInFile > lastLogTime {
		return 0, nil
	}

	// 二分查找：找第一个 > checkpoint 的记录
	left, right := int64(0), fileSize
	var result int64 = fileSize

	for left < right {
		mid := (left + right) / 2
		lineStart := findLineStart(file, mid)

		if lineStart <= left && mid > left {
			left = mid + 1
			continue
		}

		timestamp, err := readTimestampAt(file, lineStart)
		if err != nil {
			left = mid + 1
			continue
		}

		if timestamp > lastLogTime {
			result = lineStart
			right = mid
		} else {
			left = mid + 1
		}
	}

	return result, nil
}

// findLineStart 向前查找最近的换行符，返回下一行的起始位置
func findLineStart(file *os.File, pos int64) int64 {
	if pos <= 0 {
		return 0
	}

	const bufSize = 4096
	start := pos - bufSize
	if start < 0 {
		start = 0
	}

	_, err := file.Seek(start, 0)
	if err != nil {
		return 0
	}

	buf := make([]byte, int(pos-start))
	n, _ := file.Read(buf)

	for i := n - 1; i >= 0; i-- {
		if buf[i] == '\n' {
			return start + int64(i) + 1
		}
	}

	if start > 0 {
		return findLineStart(file, start)
	}
	return 0
}

// findLastLineStart 找到文件中最后一个完整行的起始位置
func findLastLineStart(file *os.File, fileSize int64) int64 {
	if fileSize <= 1 {
		return 0
	}

	// 从文件末尾向前查找，跳过可能的末尾换行符
	buf := make([]byte, 1)
	pos := fileSize - 1

	// 跳过末尾的换行符
	for pos >= 0 {
		file.Seek(pos, 0)
		file.Read(buf)
		if buf[0] != '\n' && buf[0] != '\r' {
			break
		}
		pos--
	}

	if pos < 0 {
		return 0
	}

	// 从当前位置向前找上一个换行符
	for pos >= 0 {
		file.Seek(pos, 0)
		file.Read(buf)
		if buf[0] == '\n' {
			return pos + 1
		}
		pos--
	}

	return 0
}

// readTimestampAt 读取指定偏移处的日志时间戳
// 支持 Docker JSON 格式和纯文本格式
func readTimestampAt(file *os.File, offset int64) (string, error) {
	_, err := file.Seek(offset, 0)
	if err != nil {
		return "", err
	}

	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}

	line = strings.TrimSpace(line)
	if line == "" {
		return "", fmt.Errorf("空行")
	}

	// 尝试解析 JSON 格式（Docker 日志）
	var logEntry struct {
		Log string `json:"log"`
	}
	if err := json.Unmarshal([]byte(line), &logEntry); err == nil && logEntry.Log != "" {
		// 从日志内容中提取时间戳（与 parser 保持一致）
		timeRegex := regexp.MustCompile(`^\[(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})\]`)
		matches := timeRegex.FindStringSubmatch(logEntry.Log)
		if len(matches) >= 2 {
			return strings.Replace(matches[1], " ", "T", 1) + "+08:00", nil
		}
		return "", fmt.Errorf("无法从 JSON 日志内容提取时间戳")
	}

	// 尝试解析纯文本格式：[YYYY-MM-DD HH:MM:SS]
	timeRegex := regexp.MustCompile(`^\[(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})\]`)
	matches := timeRegex.FindStringSubmatch(line)
	if len(matches) >= 2 {
		// 转换为 ISO 格式（东8区）
		localTime := matches[1]
		return strings.Replace(localTime, " ", "T", 1) + "+08:00", nil
	}

	return "", fmt.Errorf("无法提取时间戳")
}

// SaveCaveStats 保存洞穴统计结果到文件
func SaveCaveStats(stats map[string]map[string]*types.CaveDailyStats, filePath string) {
	storage.SaveStats(convertCaveStats(stats), filePath)
}

// SaveChallengeStats 保存挑战统计结果到文件
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
