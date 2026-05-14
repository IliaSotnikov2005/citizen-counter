package output

import (
	"fmt"
	"sort"
)

type pair struct {
	name  string
	count int
}

func Print(counts map[string]int, topN int, sortByFreq bool) {
	if len(counts) == 0 {
		return
	}

	pairs := make([]pair, 0, len(counts))
	for name, cnt := range counts {
		pairs = append(pairs, pair{name, cnt})
	}

	if sortByFreq {
		sort.Slice(pairs, func(i, j int) bool {
			if pairs[i].count == pairs[j].count {
				return pairs[i].name < pairs[j].name
			}
			return pairs[i].count > pairs[j].count
		})
	} else {
		sort.Slice(pairs, func(i, j int) bool {
			return pairs[i].name < pairs[j].name
		})
	}

	if topN > 0 && topN < len(pairs) {
		pairs = pairs[:topN]
		fmt.Printf("\n=== First %d names\n", topN)
	} else {
		fmt.Println("\n=== All Names")
	}

	for i, p := range pairs {
		fmt.Printf("%d. %-20s %d\n", i+1, p.name, p.count)
	}

	fmt.Printf("\nTotal: %d unique names\n", len(counts))
}
