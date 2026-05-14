package counter

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type SequentialCounter struct {
	BufferSizeMB int
}

func NewSequentialCounter(bufferMB int) *SequentialCounter {
	return &SequentialCounter{
		BufferSizeMB: bufferMB,
	}
}

func (c *SequentialCounter) Count(filePath string) (*Result, error) {
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

	bufSize := max(c.BufferSizeMB*1024*1024, 64*1024)
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
