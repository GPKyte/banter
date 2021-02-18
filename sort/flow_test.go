package sort_test

import (
	"testing"

	"github.com/GPKyte/banter/sort"
)

func TestRangeExtendOnce(t *testing.T) {
	var exampleRanges = map[int][]int{
		0:  []int{7, 3, 2, 1, 5},             // These are all included and should work fine
		5:  []int{9, 12, 4, 6},               // These are all included as well
		9:  []int{55, 16, 21, 8, 14, 17, 22}, // This has some included and others not, so will have some drop out, perhaps.
		16: []int{99, 100, 111, 122, 155},    // None of these are expected, what will happen now Mr. Crabs?
	}

	var xob = sort.OrderedBins(sort.ExampleOrderedBinSeries[0])
	var xRangeExtended = xob.ExpandIndexRangesOf(exampleRanges)
	t.Log(xRangeExtended)
}

func TestIncrementStrategy(t *testing.T) {
	var offsets = []int{0, 5, 9, 16}

	with, without := sort.ExampleSyncIncrement(offsets)
	t.Logf("\nWith: %v\nWithout: %v", with, without)
}
