package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// timestampRegex matches log timestamp format: [YYYY-MM-DD HH:MM:SS]
var timestampRegex = regexp.MustCompile(`^\[(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})\]`)

// extractTimestamp extracts timestamp from a log line.
// Supports both Docker JSON format and plain text format.
// Returns empty string if not found.
func extractTimestamp(line string) string {
	line = strings.TrimSpace(line)
	if line == "" {
		return ""
	}

	// Try JSON format (Docker logs)
	var logEntry struct{ Log string }
	if err := json.Unmarshal([]byte(line), &logEntry); err == nil && logEntry.Log != "" {
		line = logEntry.Log
	}

	// Extract timestamp
	if matches := timestampRegex.FindStringSubmatch(line); len(matches) >= 2 {
		return strings.Replace(matches[1], " ", "T", 1) + "+08:00"
	}
	return ""
}

// extractTimestampFromLine extracts timestamp from a single line.
func extractTimestampFromLine(line string) (string, error) {
	ts := extractTimestamp(line)
	if ts == "" {
		return "", fmt.Errorf("无法提取时间戳")
	}
	return ts, nil
}

// readTimestampAt reads timestamp at the specified file offset.
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

	ts := extractTimestamp(line)
	if ts == "" {
		return "", fmt.Errorf("无法提取时间戳")
	}
	return ts, nil
}

// findStartPosition uses binary search to find the first record with timestamp > lastLogTime.
func findStartPosition(file *os.File, lastLogTime string) (int64, error) {
	if lastLogTime == "" {
		return 0, nil
	}
	return findStartPositionLinear(file, lastLogTime)
}

// findStartPositionLinear finds the first record after the checkpoint by
// scanning from the beginning. Logs can contain system output and wrapped lines
// without timestamps, which makes timestamp binary search unsafe.
func findStartPositionLinear(file *os.File, lastLogTime string) (int64, error) {
	if _, err := file.Seek(0, 0); err != nil {
		return 0, err
	}

	reader := bufio.NewReader(file)
	var offset int64
	var lastValidTimestamp string

	for {
		lineStart := offset
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return 0, err
		}

		if ts := extractTimestamp(line); ts != "" {
			lastValidTimestamp = ts
			if ts > lastLogTime {
				return lineStart, nil
			}
		}

		offset += int64(len(line))
		if err == io.EOF {
			break
		}
	}

	if lastValidTimestamp != "" && lastValidTimestamp <= lastLogTime {
		return -1, nil
	}
	return -1, nil
}

// findLineStart searches backward for the nearest newline, returns next line start position.
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

// findLastLineStart finds the start position of the last complete line in file.
func findLastLineStart(file *os.File, fileSize int64) int64 {
	if fileSize <= 1 {
		return 0
	}

	buf := make([]byte, 1)
	pos := fileSize - 1

	// Skip trailing newlines
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

	// Find previous newline
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

// findFirstValidTimestamp finds the first valid timestamp in file.
// Returns timestamp and file position.
func findFirstValidTimestamp(file *os.File, fileSize int64) (string, int64) {
	_, err := file.Seek(0, 0)
	if err != nil {
		return "", 0
	}

	reader := bufio.NewReader(file)
	var offset int64 = 0

	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		if line != "" {
			timestamp, _ := extractTimestampFromLine(line)
			if timestamp != "" {
				return timestamp, offset
			}
		}

		offset += int64(len(line) + 1)
		if err == io.EOF {
			break
		}
	}

	return "", 0
}
