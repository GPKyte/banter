package sort

import (
	"fmt"
)

// Open question about interfaces: What method could we use to overlap the datastructure of bins to something more under-the-hood clever better for mocks?
// When to and not to create an interface. Perhaps after, but preferably before. I explain that after a concrete implentation is found, one can create an interface to wrap it, but any code written would need updated unless the interface explicity matched the original methods or the whole thing became deprecated with a new call provided for the up to date version. Prototyping is an alternative to In advance interfacing.

// FlowMachine maintains any useful-for-sorting metadata
type FlowMachine struct {
	size int
	obs  OrderedBinSeries
}

// OrderedBinSeries puts a name to the construct used to flow sort indices assumed to be sorted by sequentially significant attributes elsewhere
type OrderedBinSeries []OrderedBins

// OrderedBins get the order from their construction
type OrderedBins []Bin

// Bin is singularly named for a collection of indices into a collection stored elsewhere
type Bin []int

// ExampleOrderedBinSeries illustrates an OBS and is used in tests
var ExampleOrderedBinSeries = OrderedBinSeries{
	[]Bin{ // Input
		[]int{0, 1, 2, 3, 4},
		[]int{5, 6, 7, 8, 9},
		[]int{10, 11, 12, 13, 14},
	},
	[]Bin{ // Alphabetical
		[]int{5, 6, 7, 8, 9},
		[]int{0, 1, 2, 3, 4},
		[]int{10, 11, 12, 13, 14},
	},
	[]Bin{ // Disorder
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
	var reverse = make(map[int]*int) // This helps expand the indexed ranges into further indexed subranges whilst maintaining the appropriate offset

	for relativeIndex, bin := range allThese {
		var relativeIndexUsedByAddress = int(relativeIndex) // Make a copy

		for _, each := range bin {
			reverse[each] = &relativeIndexUsedByAddress
		}
	}

	toThese = make(map[int][]int)
	// Expand the range into subranges
	for _, bin := range ob {
		var sync = initIncSync()

		// We use reference to shared index to avoid many updates
		// values in "shared" bin will have different "shared" indices
		for _, value := range bin {
			var sharedIndex = reverse[value]
			if sharedIndex == nil {
				continue // Skip invalid entries, it's okay to lose stuff for now
			}
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

func (ss *simpleSync) incrementLater(please *int) {
	ss.postponed = append(ss.postponed, please)
}

func (ss *simpleSync) laterIsNow() {
	for _, each := range ss.postponed {
		(*each)++
	}
	ss.postponed = make([]*int, 0)
}

func initIncSync() *simpleSync {
	var pp = make([]*int, 0)
	var ss = simpleSync{postponed: pp}

	return &ss
}

// ExampleSyncIncrement demonstrates the difference between incrementing a shared Key immediately and after waiting
func ExampleSyncIncrement(increasingOffsets []int) (with, without map[int][]int) {
	// For each example take the offsets as inspiration and form of three, regardless of overlap and undercoverage.
	// This is an example, remember?

	// Increment a shared key immediately
	var immediate = make(map[int][]int)
	for i := range increasingOffsets {
		key := int(increasingOffsets[i])

		// Build a sequence of three, notice both the next steps
		// in this inner loop, versus the same loop in the syncronized example further beyond.
		for ii := 1; ii <= 3; ii++ {
			immediate[key] = []int{ii}
			key++
		}

		// The next two steps generate a n:[1] and n+1:[2]
		// Rather than
		for ii := 1; ii <= 2; ii++ {
			immediate[key] = []int{ii}
			key++
		}
	}

	// Increment a shared key after right before moving onto next group
	var sync = initIncSync()
	var syncedLater = make(map[int][]int)
	for i := range increasingOffsets {
		key := int(increasingOffsets[i])

		// Generate a list in the same bin
		for ii := 1; ii <= 3; ii++ {
			syncedLater[key] = append(syncedLater[key], ii)
			sync.incrementLater(&key) // Notice the key value stays the same until before next group is generated
		}
		sync.laterIsNow() // Here is the big differentiator

		for ii := 1; ii <= 2; ii++ {
			syncedLater[key] = append(syncedLater[key], ii)
			sync.incrementLater(&key)
		}
		sync.laterIsNow()
		syncedLater[key] = []int{999} // This is a helpful end of sequence mark, better than the next in sequence int{7}
	}

	return syncedLater, immediate
}

func ExampleFlow() {

}
