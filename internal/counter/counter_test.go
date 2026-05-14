package counter

import (
	"os"
	"testing"
)

var commonTests = []struct {
	name    string
	content string
	want    map[string]int
	wantErr bool
}{
	{
		name:    "Basic case",
		content: "Алёна\nМиша\nАлёна\nДима\n",
		want:    map[string]int{"Алёна": 2, "Миша": 1, "Дима": 1},
		wantErr: false,
	},
	{
		name:    "Empty file",
		content: "",
		want:    map[string]int{},
		wantErr: false,
	},
	{
		name:    "Spaces and empty lines",
		content: "  Иван  \n\n  Иван\n",
		want:    map[string]int{"Иван": 2},
		wantErr: false,
	},
	{
		name:    "Different cases",
		content: "Юлия\nюлия\nЮлия \n",
		want:    map[string]int{"Юлия": 2, "юлия": 1},
		wantErr: false,
	},
}

func runCounterTest(t *testing.T, c Counter, tt struct {
	name    string
	content string
	want    map[string]int
	wantErr bool
}) {
	tmpFile, err := os.CreateTemp("", "test_names_*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			_ = err
		}
	}()

	if _, err := tmpFile.WriteString(tt.content); err != nil {
		t.Fatal(err)
	}
	if err := tmpFile.Close(); err != nil {
		_ = err
	}

	result, err := c.Count(tmpFile.Name())
	if (err != nil) != tt.wantErr {
		t.Errorf("%s: Count() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		return
	}

	if err == nil {
		if len(result.Counts) != len(tt.want) {
			t.Errorf("%s: got %d unique names, want %d", tt.name, len(result.Counts), len(tt.want))
		}

		for k, v := range tt.want {
			if result.Counts[k] != v {
				t.Errorf("%s: for name %s: got count %d, want %d", tt.name, k, result.Counts[k], v)
			}
		}
	}
}

func TestSequentialCounter(t *testing.T) {
	c := NewSequentialCounter(1)
	for _, tt := range commonTests {
		t.Run(tt.name, func(t *testing.T) {
			runCounterTest(t, c, tt)
		})
	}
}

func TestParallelCounter(t *testing.T) {
	workerCounts := []int{1, 2, 4}

	for _, wc := range workerCounts {
		c := NewParallelCounter(wc, 1)
		for _, tt := range commonTests {
			t.Run(tt.name, func(t *testing.T) {
				runCounterTest(t, c, tt)
			})
		}
	}
}

func TestEstimateCapacity(t *testing.T) {
	tests := []struct {
		size int64
		want int
	}{
		{size: 10, want: minExpectedUniques},
		{size: 100 * 1024 * 1024, want: maxExpectedUniques},
	}

	for _, tt := range tests {
		got := estimateCapacity(tt.size)
		if got < minExpectedUniques || got > maxExpectedUniques {
			t.Errorf("estimateCapacity(%d) = %d; out of bounds [%d, %d]",
				tt.size, got, minExpectedUniques, maxExpectedUniques)
		}
	}
}
