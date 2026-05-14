package counter

const (
	minExpectedUniques = 100
	maxExpectedUniques = 20_000
)

type Result struct {
	Counts map[string]int
}

type Counter interface {
	Count(path string) (*Result, error)
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
