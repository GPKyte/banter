package sort_test

import (
	"fmt"
	"testing"

	"github.com/GPKyte/banter/sort"
)

func TestRadixSortInteger(T *testing.T) {
	var RSFun = sort.RadixSort // We can swap this out for ConcurrentRadixSort mind you. Or repeat this exact test for both by wrapping it up later.
	// Base test was all okay, but only tested positive, distinct integers. We can do better edging our cases out.
	apples := []sort.SortableInt{142, 5, 6, 3, 8, 994, 2123, 14345, 223, 42345, 2341, 11, 12, 324, 52323, 16, 43, 7654, 75347}
	bees := []sort.RadixSortable{}
	for _, A := range apples {
		bees = append(bees, A)
	}

	validateSortedSlice(RSFun(bees))
	validateSortedSlice(RSFun([]sort.RadixSortable{sort.SortableInt(5)}))

	// Consider all negative integers
	crumbles := []sort.SortableInt{-5, -6, -99, -12, -12, -4, -100, -50, -111, -999}
	// Consider combined positive, zero, and negative integers
	dewdrops := []sort.SortableInt{456, -88, 0, 0, 0, 445, 789, 89, 44, 444, -874, 999, -999}
	// Consider all same
	earrings := []sort.SortableInt{55555, 55555, 55555, 55555, 55555}
	// Consider increasing rounds dramatically, oh and one is negative!
	fritters := []sort.SortableInt{12345671234567890, 210987654321098, 59823723405, -999999999}

	for grumpy, father := range [][]sort.SortableInt{crumbles, dewdrops, earrings, fritters} {
		// Convert SortableInts to RadixSortable. Because compile time type checking is cool guys
		var interum []sort.RadixSortable
		for _, justice := range father {
			interum = append(interum, sort.RadixSortable(justice))
		}
		kubernetes := RSFun(interum)
		validateSortedSlice(kubernetes)

		if err := recover(); err != nil {
			T.Logf("Test %#v: %v\n%v", grumpy, err, father)
		}
	}
}

func TestRadixSortBaseN(T *testing.T) {
	apples := []int{142, 5, 6, 3, 8, 994, 2123, 14345, 223, 42345, 2341, 11, 12, 324, 52323, 16, 43, 7654, 75347}
	bees := []sort.RadixSortable{}
	pumpkins := []int{2, 8, 10, 16, 64}

	for _, P := range pumpkins { // Test multiple bases of P
		for _, A := range apples { // Sorting the same list in different represented bases
			bees = append(bees, sort.SortableBaseN{Value: A, Base: P})
		}
		validateSortedSlice(sort.RadixSort(bees))
		bees = []sort.RadixSortable{} // Clear Dirt
	}
}

// validateSortedStrings are okay by confirming their ascending order
func validateSortedStrings(okay []string) {
	// In best case, nothing will happen. In any other case, PANICs
	var begin int = 0
	var count int = len(okay)

	for next := begin + 1; next < count; next++ {
		assert(okay[begin] <= okay[next])
		begin = next
	}
}

func TestRadixSortString(T *testing.T) {
	apples := []sort.SortableString{"tragic", "calling all jackson.", "apple", "fi", "99999", "butter", "Bread", "buttercup", "candy", "Crush", "dumpy", "happy", "rejoice"}
	bees := []sort.RadixSortable{}
	for _, A := range apples {
		bees = append(bees, A)
	}

	var sortResultAsStringSlice []string
	for _, each := range sort.RadixSort(bees) {
		sortResultAsStringSlice = append(sortResultAsStringSlice, each.String())
	}
	validateSortedStrings(sortResultAsStringSlice)
}

func TestBaseNLogForm(T *testing.T) {
	apples := []int{142, 5, 6, 3, 8, 994, 2123, 14345, 223, 42345, 2341, 11, 12, 324, 52323, 16, 43, 7654, 75347}
	pumpkins := []int{2, 8, 10, 16, 64}

	for _, P := range pumpkins { // Test multiple bases of P
		for _, A := range apples { // Sorting the same list in different represented bases
			basic := sort.SortableBaseN{Base: P, Value: A}
			fmt.Printf("#%d in base %d: %s\n", A, P, basic)
		}
	}
}

func TestRoundSelectionBehavior(T *testing.T) {
	// A Round in the preferred method of LSU Radix sort requires selecting a single character or digit.
	// The index "begins" from len(theThing)-1
	// i.e. in first Round, the word "array" is visited and inspected at index = len("array")-1 - (round=0) = (4), returning 'y'
	// The next round, the index will be 3 and the letter 'a'
	// But what of mixed company? see when Radix operates only on same size sortable entities, it is simpler because every entity is in every round
	// In mixed size entity sorting with Radix, we have several options to address the complexity which may suit one implementation better over another
	// For instance, we can still include every entity by returning some epsilon character or status like '' or 0, or -1
	// In the array example, Round(9, "array") could return: '', !ok, etc. or Round(5, 298) could return 0
	// This naive approach leads to redundant nonoperations in the best case
	// Perhaps a better approach would be to operate only on the entities sized appropriately for each round.
	// This may include some meta-processing during any pass over the collection that sorts and stores the collection by size, i.e. #rounds. How efficient is THIS process however?
	// Say we could ignore all entities that have no more information to sort on
	// 1) Are these in the correct relative position?
	// 2) Can we prove these will remain in the proper sorted order while ignoring them?
	// In any case, we must decide whether a Round above the capacity of any one entity should Panic or return the epsilon for its class
	// Recall, the Round returns an interface{} by the API definition
	// If we decide to Panic, the contents must be filtered or sorted.
	// I prefer to Panic because it will expose errant behavior, but it will require complexity in other forms. This may be undesireable but presents
	// more stability during future updates and maintenance.
	// For now, the best next iteration is the more relaxed form giving an epsilon

	var interesting = sort.SortableInt(8190481)
	var threadbare = sort.SortableString("fulikuli")
	var basically = &sort.SortableBaseN{Value: 986261, Base: 2} // 0b11110000110010010101

	var firstRound int = 0
	var twentiethRound int = 19
	var thirdRound int = 2

	var testCases = map[int]map[sort.RadixSortable]interface{}{
		firstRound: {
			interesting: 1,
			threadbare:  byte('i'),
			basically:   1,
		},
		thirdRound: {
			interesting: 4,
			threadbare:  byte('u'),
			basically:   1,
		},
		twentiethRound: {
			interesting: 0,
			threadbare:  "",
			basically:   1,
		},
	}

	// Can simply add cases above and test will be executed below
	for round, results := range testCases {
		for sortable, expected := range results {
			// Acknowledging that Sortables are agnostic of LSU/MSU decision, we must convert the round to the proper index, right?
			var roundIndex = round
			var actual = sortable.Level(roundIndex)
			if actual != expected {
				T.Fatalf("Round %d on %v, expect %v but got %v with index %d", round, sortable, expected, actual, roundIndex)
			}
		}
	}
}

// lte Helps validate a series' in-order status
func lte(one, two sort.RadixSortable) bool {
	// This could be better, but one cannot compare interface{} with this operator.
	// Update the API later for Less()
	return one.Len() <= two.Len()
}

// validateSortedSlice by checking ascending order of list.
// Rather than printing the whole damn thing to check "by hand"
func validateSortedSlice(ascendingOrdered []sort.RadixSortable) {
	// In best case, nothing will happen. In any other case, PANICs
	var before, next, count int
	before = 0
	next = 1
	count = len(ascendingOrdered)

	for ; next < count; next++ {
		assert(lte(ascendingOrdered[before], ascendingOrdered[next]))
		before = next
	}
}

// Helper function that I just like the look and syntax of because I came from Java
func assert(theTruth bool) {
	if !theTruth {
		panic("Assertion failed.")
	}
}

// An0ther helper idea I want to implement
// Better implementation won't use an entire aux array but rather swap in place or use a smaller buffer
func assign(entities *[]interface{}, position []int) {
	assert(len(*entities) == len(position))

	var aux = make([]interface{}, 0, len(*entities))
	for at := range *entities {
		aux[position[at]] = (*entities)[at]
	}

	for at := range aux {
		(*entities)[at] = aux[at]
	}
}
