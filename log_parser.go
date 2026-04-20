package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// ParsedLog 表示解析后的日志结构
type ParsedLog struct {
	RawLine      string // 原始行内容
	Timestamp    string // 时间戳（日志内容中的，秒级精度）
	PreciseTime  string // 精确时间（JSON time字段，纳秒级精度，用于检查点）
	Character    string // 角色名
	Body         string // 日志主体
	IsValid      bool   // 格式是否有效
}

// parseLogLine 解析单行日志，提取通用信息
// 统一解析 "[时间] [名字] 日志主体" 格式，供后续处理复用
func parseLogLine(line string) ParsedLog {
	result := ParsedLog{RawLine: line}

	// 解析JSON格式的日志
	var logEntry struct {
		Log  string `json:"log"`
		Time string `json:"time"`
	}
	if err := json.Unmarshal([]byte(line), &logEntry); err != nil {
		return result
	}

	logContent := logEntry.Log
	result.PreciseTime = logEntry.Time // 保存精确时间用于检查点

	// 验证并提取 "[时间] [名字] 日志主体" 格式
	timestamp, character, body := extractLogComponents(logContent)
	if timestamp == "" {
		return result
	}

	result.Timestamp = timestamp
	result.Character = character
	result.Body = body
	result.IsValid = true

	return result
}

// extractLogComponents 从日志内容中提取时间戳、角色名和日志主体
// 返回: (timestamp, character, body)
func extractLogComponents(logContent string) (string, string, string) {
	logContent = strings.TrimSpace(logContent)

	// 检查基本格式：必须以 [ 开始
	if !strings.HasPrefix(logContent, "[") {
		return "", "", ""
	}

	// 提取时间戳：第一个方括号内的内容
	timeRegex := regexp.MustCompile(`^\[(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})\]`)
	timeMatches := timeRegex.FindStringSubmatch(logContent)
	if len(timeMatches) < 2 {
		return "", "", ""
	}
	timestamp := timeMatches[1]

	// 查找第一个 ] 后的内容
	firstCloseIdx := strings.Index(logContent, "]")
	if firstCloseIdx == -1 || firstCloseIdx+1 >= len(logContent) {
		return "", "", ""
	}

	// 提取角色名和日志主体：第二个方括号内的内容
	// 格式: [时间] [名字 (Lv等级)] 日志主体
	remaining := strings.TrimSpace(logContent[firstCloseIdx+1:])

	// 必须以 [ 开始
	if !strings.HasPrefix(remaining, "[") {
		return "", "", ""
	}

	// 找到第二个 ]
	secondCloseIdx := strings.Index(remaining, "]")
	if secondCloseIdx == -1 || secondCloseIdx <= 1 {
		return "", "", ""
	}

	// 提取角色名（去掉等级部分）
	nameSection := remaining[1:secondCloseIdx] // 去掉开头的 [
	character := extractCharacterName(nameSection)

	// 提取日志主体（第二个 ] 后面的内容）
	if secondCloseIdx+1 >= len(remaining) {
		return "", "", ""
	}
	body := strings.TrimSpace(remaining[secondCloseIdx+1:])

	if character == "" || body == "" {
		return "", "", ""
	}

	return timestamp, character, body
}

// extractCharacterName 从名字部分提取角色名（去掉等级）
// 输入: "角色名 (Lv100)" 输出: "角色名"
func extractCharacterName(nameSection string) string {
	nameSection = strings.TrimSpace(nameSection)

	// 去掉等级部分 (Lvxxx)
	if idx := strings.Index(nameSection, "(Lv"); idx != -1 {
		return strings.TrimSpace(nameSection[:idx])
	}

	return nameSection
}

// isValidSource 检查日志主体是否可以作为钻石来源
// 返回 true 如果不是以下类型：
// - 物品变动日志（以"Name:"开头）
// - 挑战记录（以"Challenge"开头）
// - 错误日志（以"OnError"开头）
func isValidSource(body string) bool {
	if strings.HasPrefix(body, "Name:") {
		return false
	}
	if strings.HasPrefix(body, "Challenge") {
		return false
	}
	if strings.HasPrefix(body, "OnError") {
		return false
	}
	return true
}

// LogType 日志类型枚举
type LogType int

const (
	LogTypeNone           LogType = iota // 未知类型
	LogTypeDiamond                        // 钻石记录
	LogTypeCave                           // 洞穴记录
	LogTypeChallenge                      // 挑战记录
	LogTypeRuneTicket                     // 饼干记录
	LogTypeUpgradePanacea                 // 红水记录
)

// identifyLogType 识别日志类型（一次扫描确定类型）
func identifyLogType(body string) LogType {
	// 检查钻石记录（包含 Diamonds 关键字）
	if strings.Contains(body, "Diamonds(None)") {
		return LogTypeDiamond
	}

	// 检查饼干记录
	if runeTicketRegex.MatchString(body) {
		return LogTypeRuneTicket
	}

	// 检查红水记录
	if upgradePanaceaRegex.MatchString(body) {
		return LogTypeUpgradePanacea
	}

	// 检查洞穴记录（包含洞穴关键字）
	if caveEnterRegex.MatchString(body) || caveFinishRegex.MatchString(body) || caveErrorRegex.MatchString(body) {
		return LogTypeCave
	}

	// 检查挑战记录（以 Challenge 开头且包含状态关键字）
	if strings.HasPrefix(body, "Challenge") && (challengeSuccessRegex.MatchString(body) || challengeFailedRegex.MatchString(body)) {
		return LogTypeChallenge
	}

	return LogTypeNone
}

// parseDiamondLine 从已解析的日志中提取钻石记录
func parseDiamondLine(parsed ParsedLog, lastSource string) *DiamondRecord {
	body := parsed.Body

	// 查找获取钻石的记录
	if matches := diamondGainRegex.FindStringSubmatch(body); len(matches) > 1 {
		amount, _ := strconv.Atoi(matches[1])
		return &DiamondRecord{
			Character:   parsed.Character,
			Timestamp:   parsed.Timestamp,
			PreciseTime: parsed.PreciseTime,
			Amount:      amount,
			Type:        "gain",
			Source:      lastSource,
		}
	}

	// 查找消耗钻石的记录
	if matches := diamondConsumeRegex.FindStringSubmatch(body); len(matches) > 1 {
		amount, _ := strconv.Atoi(matches[1])
		return &DiamondRecord{
			Character:   parsed.Character,
			Timestamp:   parsed.Timestamp,
			PreciseTime: parsed.PreciseTime,
			Amount:      amount,
			Type:        "consume",
			Source:      lastSource,
		}
	}

	return nil
}

// parseCaveLine 从已解析的日志中提取洞穴记录
func parseCaveLine(parsed ParsedLog) *CaveRecord {
	body := parsed.Body
	date := parsed.Timestamp[:10]

	// 检查三种状态（优先级：error > finished > started）
	if caveErrorRegex.MatchString(body) {
		return &CaveRecord{
			Character:   parsed.Character,
			Timestamp:   parsed.Timestamp,
			PreciseTime: parsed.PreciseTime,
			Status:      CaveStatusError,
			Date:        date,
		}
	}

	if caveFinishRegex.MatchString(body) {
		return &CaveRecord{
			Character:   parsed.Character,
			Timestamp:   parsed.Timestamp,
			PreciseTime: parsed.PreciseTime,
			Status:      CaveStatusFinished,
			Date:        date,
		}
	}

	if caveEnterRegex.MatchString(body) {
		return &CaveRecord{
			Character:   parsed.Character,
			Timestamp:   parsed.Timestamp,
			PreciseTime: parsed.PreciseTime,
			Status:      CaveStatusStarted,
			Date:        date,
		}
	}

	return nil
}

// parseChallengeLine 从已解析的日志中提取挑战记录
func parseChallengeLine(parsed ParsedLog) *ChallengeRecord {
	body := parsed.Body

	record := &ChallengeRecord{
		Character:   parsed.Character,
		Timestamp:   parsed.Timestamp,
		PreciseTime: parsed.PreciseTime,
	}

	// 检测状态：优先检测成功，否则为失败
	if challengeSuccessRegex.MatchString(body) {
		record.Status = ChallengeStatusSuccess
	} else {
		record.Status = ChallengeStatusFailed
	}

	// 尝试匹配主线挑战: "Challenge 43-6 boss one time：..."
	if matches := challengeQuestRegex.FindStringSubmatch(body); len(matches) > 1 {
		record.Type = ChallengeTypeQuest
		record.Level = matches[1]
		return record
	}

	// 尝试匹配塔挑战: "Challenge Tower of Amber 1303 layer one time：..."
	if matches := challengeTowerRegex.FindStringSubmatch(body); len(matches) > 2 {
		record.Type = ChallengeTypeTower
		record.TowerType = TowerType(matches[1])
		record.Level = matches[2]
		return record
	}

	return nil
}

// parseRuneTicketLine 从已解析的日志中提取饼干记录
func parseRuneTicketLine(parsed ParsedLog, lastSource string) *ItemRecord {
	body := parsed.Body

	matches := runeTicketRegex.FindStringSubmatch(body)
	if len(matches) < 2 {
		return nil
	}

	amount, _ := strconv.Atoi(matches[1])

	return &ItemRecord{
		Character:   parsed.Character,
		Timestamp:   parsed.Timestamp,
		PreciseTime: parsed.PreciseTime,
		Amount:      amount,
		Source:      lastSource,
	}
}

// parseUpgradePanaceaLine 从已解析的日志中提取红水记录
func parseUpgradePanaceaLine(parsed ParsedLog, lastSource string) *ItemRecord {
	body := parsed.Body

	matches := upgradePanaceaRegex.FindStringSubmatch(body)
	if len(matches) < 2 {
		return nil
	}

	amount, _ := strconv.Atoi(matches[1])

	return &ItemRecord{
		Character:   parsed.Character,
		Timestamp:   parsed.Timestamp,
		PreciseTime: parsed.PreciseTime,
		Amount:      amount,
		Source:      lastSource,
	}
}

// mapSourceToAlias 将来源日志映射为友好的来源名称
func mapSourceToAlias(source string) string {
	// 优先级：特殊短语 > 模式匹配

	// 特殊映射
	if source == "You have triumphed." {
		return "temple of illusions"
	}

	// Gacha 日志：提取 " N times" 之前的内容
	// 例如: "Gacha 黒葬武具ガチャ 5 times, used Gold×250000" -> "Gacha 黒葬武具ガチャ"
	if strings.HasPrefix(source, "Gacha ") {
		if idx := strings.Index(source, " times"); idx != -1 {
			// 找到 " times" 后，向前查找 " N" 的起始位置
			beforeTimes := source[:idx]
			// 从末尾向前找空格，提取数字之前的部分
			if lastSpace := strings.LastIndex(beforeTimes, " "); lastSpace != -1 {
				return strings.TrimSpace(beforeTimes[:lastSpace])
			}
		}
	}

	// Open 日志：提取 " x N" 之前的内容
	// 例如: "Open Silver Sealed Chest x 2" -> "Open Silver Sealed Chest"
	if strings.HasPrefix(source, "Open ") {
		if idx := strings.Index(source, " x "); idx != -1 {
			return strings.TrimSpace(source[:idx])
		}
	}

	return source
}

// readTimestampAt 读取指定偏移处的日志时间戳（精确时间，用于二分查找）
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

	// 从JSON中提取精确时间
	var logEntry struct {
		Time string `json:"time"`
	}
	if err := json.Unmarshal([]byte(line), &logEntry); err != nil {
		return "", err
	}

	if logEntry.Time == "" {
		return "", fmt.Errorf("无法提取时间戳")
	}

	return logEntry.Time, nil
}

// findLineStart 向前查找最近的换行符，返回下一行的起始位置
// 确保二分查找定位到行首，避免从行中间开始读取
func findLineStart(file *os.File, pos int64) int64 {
	if pos <= 0 {
		return 0
	}

	// 向前读取最多 4KB（足够容纳多行）
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

	// 从后向前找换行符
	for i := n - 1; i >= 0; i-- {
		if buf[i] == '\n' {
			return start + int64(i) + 1
		}
	}

	// 没找到，继续向前或返回文件开头
	if start > 0 {
		return findLineStart(file, start)
	}
	return 0
}

// findStartPosition 二分查找定位第一个时间戳 > lastLogTime 的字节位置
// 对于上百MB的文件，只需约 20 次读取即可定位（log2(100MB/4KB) ≈ 15）
func findStartPosition(file *os.File, lastLogTime string) (int64, error) {
	if lastLogTime == "" {
		return 0, nil // 从头开始处理
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return 0, err
	}
	fileSize := fileInfo.Size()
	if fileSize == 0 {
		return 0, nil
	}

	// 快速边界检查：读取文件末尾附近，检查是否有新记录
	// 从文件末尾向前找行首，确保读取完整行
	endPos := fileSize - 4096
	if endPos < 0 {
		endPos = 0
	}
	endPos = findLineStart(file, endPos) // 确保从行首开始读取
	lastTimeInFile, err := readTimestampAt(file, endPos)
	if err == nil && lastTimeInFile <= lastLogTime {
		return -1, nil // 无新记录，无需处理
	}

	// 快速边界检查：读取文件开头，如果第一条 > lastLogTime，处理整个文件
	firstTimeInFile, err := readTimestampAt(file, 0)
	if err == nil && firstTimeInFile > lastLogTime {
		return 0, nil // 处理整个文件
	}

	// 二分查找：找到第一个时间戳 > lastLogTime 的位置
	left, right := int64(0), fileSize
	var result int64 = fileSize

	for left < right {
		mid := (left + right) / 2
		// 调整到行首位置，确保读取完整行
		lineStart := findLineStart(file, mid)

		// 避免死循环：如果行首位置没有前进
		if lineStart <= left && mid > left {
			left = mid + 1
			continue
		}

		timestamp, err := readTimestampAt(file, lineStart)
		if err != nil {
			// 读取失败，可能是损坏行，向右调整
			left = mid + 1
			continue
		}

		if timestamp > lastLogTime {
			// 这个位置可能是目标，记录并继续向左找更小的
			result = lineStart
			right = mid
		} else {
			// 时间戳 <= lastLogTime，向右查找
			left = mid + 1
		}
	}

	return result, nil
}

// processStream 流式处理日志（内存友好，适合 GB 级文件）
// 直接将记录添加到聚合器，不缓存所有记录
func (p *LogProcessor) processStream(
	agg *Aggregator,
	caveAgg *CaveAggregator,
	challengeAgg *ChallengeAggregator,
	runeTicketAgg *RuneTicketAggregator,
	upgradePanaceaAgg *UpgradePanaceaAggregator,
) string {
	// 打开日志文件
	file, err := os.Open(p.inputLogPath)
	if err != nil {
		log.Fatalf("打开日志文件失败: %v", err)
	}
	defer file.Close()

	// 二分查找定位起始位置
	startPos, err := findStartPosition(file, p.checkpoint)
	if err != nil {
		log.Printf("定位起始位置失败: %v，从头开始处理", err)
		startPos = 0
	}
	if startPos < 0 {
		log.Println("无新记录需要处理")
		return p.checkpoint
	}

	// 跳转到起始位置
	_, err = file.Seek(startPos, 0)
	if err != nil {
		log.Fatalf("文件定位失败: %v", err)
	}

	// 来源上下文（按角色隔离，处理期间必须保留）
	lastSourceByCharacter := make(map[string]string)

	// 分块读取：256KB 缓冲区
	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 256*1024)
	scanner.Buffer(buf, 1024*1024)

	var lastLogTime string
	lineCount := 0
	newRecordCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		// 解析日志行
		parsed := parseLogLine(line)
		if !parsed.IsValid {
			continue
		}

		// 时间戳二次校验（防止定位偏差或时间乱序）
		// 使用精确时间进行比较，避免秒级时间戳重复导致遗漏
		if p.checkpoint != "" && parsed.PreciseTime != "" && parsed.PreciseTime <= p.checkpoint {
			continue
		}

		// 更新来源上下文（非特殊日志才作为来源）
		if isValidSource(parsed.Body) {
			lastSourceByCharacter[parsed.Character] = parsed.Body
		}

		// 识别日志类型并分发给对应处理函数
		logType := identifyLogType(parsed.Body)
		switch logType {
		case LogTypeDiamond:
			// 提取钻石记录，来源为空时使用 "none"
			source := lastSourceByCharacter[parsed.Character]
			if source == "" {
				source = "none"
			}
			record := parseDiamondLine(parsed, source)
			if record != nil {
				agg.AddRecord(*record)
				newRecordCount++
				// 使用精确时间作为检查点
				if record.PreciseTime != "" {
					lastLogTime = record.PreciseTime
				} else {
					lastLogTime = record.Timestamp
				}
			}

		case LogTypeCave:
			// 提取洞穴记录
			caveRecord := parseCaveLine(parsed)
			if caveRecord != nil {
				caveAgg.AddRecord(*caveRecord)
			}

		case LogTypeChallenge:
			// 提取挑战记录
			challengeRecord := parseChallengeLine(parsed)
			if challengeRecord != nil {
				challengeAgg.AddRecord(*challengeRecord)
			}

		case LogTypeRuneTicket:
			// 提取饼干记录
			rawSource := lastSourceByCharacter[parsed.Character]
			source := mapSourceToAlias(rawSource)
			if source == "" {
				source = "none"
			}
			record := parseRuneTicketLine(parsed, source)
			if record != nil {
				runeTicketAgg.AddRecord(*record)
			}

		case LogTypeUpgradePanacea:
			// 提取红水记录
			rawSource := lastSourceByCharacter[parsed.Character]
			source := mapSourceToAlias(rawSource)
			if source == "" {
				source = "none"
			}
			record := parseUpgradePanaceaLine(parsed, source)
			if record != nil {
				upgradePanaceaAgg.AddRecord(*record)
			}
		}
		// line、parsed 在此作用域结束后可被 GC 回收
	}

	if err := scanner.Err(); err != nil {
		log.Printf("读取日志文件时出错: %v", err)
	}

	log.Printf("扫描了 %d 行，新增 %d 条记录", lineCount, newRecordCount)

	return lastLogTime
}
