package sort

// Several types to aid flow sort
type FlowMachine struct {
	size int
	obs  OrderedBinSeries
}
type OrderedBinSeries []OrderedBins
type OrderedBins []Bin
type Bin []int

var ExampleOrderedBinSeries = OrderedBinSeries{
	[]Bin{
		[]int{0, 1, 2, 3, 4},
		[]int{5, 6, 7, 8, 9},
		[]int{10, 11, 12, 13, 14},
	},
	[]Bin{
		[]int{5, 6, 7, 8, 9},
		[]int{0, 1, 2, 3, 4},
		[]int{10, 11, 12, 13, 14},
	},
	[]Bin{
		[]int{10, 11, 12, 13, 14},
		[]int{5, 6, 7, 8, 9},
		[]int{0, 1, 2, 3, 4},
	},
}

// flow orders indices
// Given a series of groupings of ordered bins
// Return the ordered indices
func flow(obs OrderedBinSeries) []int {
	// Sanity check our data
	// TODO: check len(obs) type(obs[0])...
	sample := make([]int, 0)

	return sample
}

// Flow helper recurses... probably
func flowHelper(chan struct{ data, place int }, []int) {

}

// ExpandIndexRangesOf allThese intoThese which will have more indices and more ranges but the same contents
func (ob OrderedBins) ExpandIndexRangesOf(allThese map[int][]int) (toThese map[int][]int) {
	// Build a reverse lookup; given [value] return the reference to a copy of the initial offset given by allThese
	var reverse map[int]*int // This helps expand the indexed ranges into further indexed subranges whilst maintaining the appropriate offset

	for relativeIndex, bin := range allThese {
		var relativeIndexUsedByAddress = int(relativeIndex) // Make a copy

		for _, each := range bin {
			reverse[each] = &relativeIndexUsedByAddress
		}
	}

	// Expand the range into subranges
	for _, bin := range ob {
		var sync = initIncSync()

		// We use reference to shared index to avoid many updates
		// values in "shared" bin will have different "shared" indices
		for _, value := range bin {
			var sharedIndex = reverse[value]
			var key = *sharedIndex

			sync.incrementLater(sharedIndex)
			// Note that the shared Index of a range will have increasing values each iteration,
			// Thus, the "expansion" occurs when we change the value of the shared index and reuse
			// it as a key in the finally returned datastructure "toThese"
			toThese[key] = append(toThese[key], value)
		}
		// This effectively increments the relative, i.e. start, index for a range of items by len(range)
		sync.laterIsNow()
	}

	return toThese
}

type simpleSync struct {
	postponed []*int
}

func (ss simpleSync) incrementLater(please *int) {
	ss.postponed = append(ss.postponed, please)
}

func (ss simpleSync) laterIsNow() {
	for _, each := range ss.postponed {
		(*each)++
	}
	ss.postponed = make([]*int, 0)
}

func initIncSync() simpleSync {
	var pp = make([]*int, 0)
	var ss = simpleSync{postponed: pp}

	return ss
}

func ExampleFlow() {

}
