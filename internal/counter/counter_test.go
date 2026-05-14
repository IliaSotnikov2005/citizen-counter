package counter

import (
	"os"
	"testing"
)

func TestCount(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    map[string]int
		wantErr bool
	}{
		{
			name:    "Base case",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "test_names_*.txt")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpFile.Name())

			if _, err := tmpFile.WriteString(tt.content); err != nil {
				t.Fatal(err)
			}
			tmpFile.Close()

			res, err := Count(tmpFile.Name(), 1)
			if (err != nil) != tt.wantErr {
				t.Errorf("Count() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(res.Counts) != len(tt.want) {
				t.Errorf("got %d unique names, want %d", len(res.Counts), len(tt.want))
			}

			for k, v := range tt.want {
				if res.Counts[k] != v {
					t.Errorf("for name %s: got count %d, want %d", k, res.Counts[k], v)
				}
			}
		})
	}
}

func TestEstimateCapacity(t *testing.T) {
	tests := []struct {
		size int64
		want int
	}{
		{size: 89, want: minExpectedUniques},
		{size: 9999999999999999, want: maxExpectedUniques},
	}

	for _, tt := range tests {
		got := estimateCapacity(tt.size)
		expected := tt.want
		if got != expected && (tt.size > 1000) {
			t.Errorf("estimateCapacity(%d) = %d, want %d", tt.size, got, expected)
		}
	}
}
