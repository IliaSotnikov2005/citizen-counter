package counter

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	minExpectedUniques = 100
	maxExpectedUniques = 20_000
)

type Result struct {
	Counts map[string]int
}

func Count(filePath string, bufferMB int) (*Result, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}
	estimatedCap := estimateCapacity(info.Size())

	counts := make(map[string]int, estimatedCap)

	bufSize := max(bufferMB*1024*1024, 64*1024)
	reader := bufio.NewReaderSize(file, bufSize)

	var totalLines int64

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("failed to read line: %w", err)
		}

		name := strings.TrimSpace(line)
		if name != "" {
			counts[name]++
		}
		totalLines++
	}

	return &Result{
		Counts: counts,
	}, nil
}

func estimateCapacity(fileSize int64) int {
	estimatedLines := fileSize / 15
	estimatedUnique := int(float64(estimatedLines) * 0.3)

	switch {
	case estimatedUnique < minExpectedUniques:
		return minExpectedUniques
	case estimatedUnique > maxExpectedUniques:
		return maxExpectedUniques
	default:
		return estimatedUnique
	}
}
