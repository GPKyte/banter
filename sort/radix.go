package sort // import "github.com/GPKyte/banter/sort"

import (
	"fmt"
	"strings"
)

type stringer interface {
	String() string
}

// Radix avoids comparison-sort constraints by placing Values given their indexed Value.
// The data types successful here are Base N and have implied incremental powers like any alphanumeric system.
// Conversions are possible if they imply order through the indexing, such as roman alphabet inside Unicode representation.
// Base can be binary, decimal, 26, 16, 8, etc. in theory at least.

// RadixSortable is the minimum method set interface for a list to be sortable by my Radix implementation.
// We need some interfaces and methods to consume them and this fills that need.
type RadixSortable interface {
	Level(int) interface{} // aka index, digit, power, level. Return should be in the symbol language and have implied order elsewhere
	Len() int
	stringer
}

// SortableInt implements A Level method which returns the corresponding arithmetic Base10 digit with conventional placement from small digit to large digit
// Someone please improve this wording
// Example SortableInt{89192}.Level(0) = 89192 & 10^0
type SortableInt int

// SortableBaseN is a challenge type. The idea is to apply Radix to new problems. "New"
// QUESTION: Does sorting work the same in every Base?
type SortableBaseN struct {
	Base       int
	Value      int
	cachedRepr []int
}

// SortableString wraps the string primitive to implement the interface for RadixSortable
type SortableString string

// Level is implemented on SortableInt per the interface RadixSortable
// Note how Level 5 of SortableInt{432} is 0
// 432 % 10**5 / 10**4 = 00432 / 10000 = 0.0... = {0}
func (S SortableInt) Level(this int) interface{} {
	return S.LSD(this)
}

// LSD for SortableInt fulfills interface RadixSortableLSD
func (S SortableInt) LSD(offset int) interface{} {
	// 0 offset returns the first digit (%10)
	// 1 offset returns second digit S/10 %10
	// 200 offset returns 201st digit S/10e200 %10
	return (int(S) / powerOf(10, offset)) % 10
}

// Len returns the number of Base 10 digits in this alias of int
func (S SortableInt) Len() int {
	return psuedoLog10(int(S))
}

func (S SortableInt) String() string {
	return fmt.Sprint(int(S))
}

// Level is implemented for SortableBaseN per the interface RadixSortable
func (S SortableBaseN) Level(this int) interface{} {
	return S.LSD(this)
}

// Len implemented on SortableBaseN for RadixSortable interface{}
func (S SortableBaseN) Len() int {
	return len(S.LogForm())
}

// LSD on SortableBaseN fulfills interface for RadixSortableLSD and provides safe LSD indexing
func (S SortableBaseN) LSD(level int) interface{} {
	const epsilon = -1

	if level >= S.Len() {
		return epsilon
	}
	return S.LogForm()[S.Len()-1-level]
}

// String is implemented on SortableBaseN for printing in the right format
func (S SortableBaseN) String() string {
	var formatted []int = S.LogForm()
	var stringNumbers []string = make([]string, 0, len(formatted))

	for _, v := range formatted {
		stringNumbers = append(stringNumbers, fmt.Sprint(v))
	}
	return strings.Join(stringNumbers, ":")
}

// LogForm helps represent S whether it be binary, hexadecimal, octet, IPAddress, or even decimal, the default representation
// LogForm cannot yet be applied to negative numbers, probably.
func (S SortableBaseN) LogForm() (format []int) {
	// PART ONE: Iterate up and track exponent which make Base^exponent greater than running remainder from S.Value
	// PART TWO: Now for the duration it takes to reduce the running remainder to zero ~exactly~ or less than 1,
	// Decrement the exponent
	// PART THREE: For the duration the (coeff := Base-1) * S.Base ^ ordinalExponent is greater than runningRemainder
	// Decrement the coeff
	// PART FOUR: Now the coeff * Base ^ exponent is the maximal combination, record the coefficient into the place indicated by the ordinalExponent
	// NOTES: The trick to this without knowing in advance seems to be appending the coefficient to a growing slice
	// But remember: If one decrements the coeff to Zero before finding the maximal combination and continuing on the decrement the ordinalExponent
	// Care must be taken to append a 0 in this location.
	// ALT STRATEGY: Alternatively one could initialize a slice of some size (in the end this size must be exactly correct to display properly later on)
	// Then Every time a maximal coeff is found, the coeff is recorded in the exact placce indicated by the ordinalExponent
	// Before decrementing said location and repeating until finished at ordinalExpoent == 0; coeff == {0..Base-1}

	if S.cachedRepr != nil {
		return S.cachedRepr
	}
	/* PART ONE */
	var ordinalExponent, coeff int

	for runningProduct := 1; runningProduct < S.Value; runningProduct *= S.Base {
		ordinalExponent++
	}
	/* PART TWO */
	ordinalExponent--
	for runningRemainder := S.Value; ordinalExponent >= 0; ordinalExponent-- {
		/* PART THREE */
		for coeff = S.Base - 1; coeff > 0 && coeff*powerOf(S.Base, ordinalExponent) > runningRemainder; coeff-- {
			// Nothing needs done here
		}
		/* PART FOUR */
		runningRemainder -= coeff * powerOf(S.Base, ordinalExponent)
		format = append(format, coeff)
	}
	S.cachedRepr = format

	return format
}

// ExampleLogForm sets expectations for behavior of custom Base switching mechaninsm
func ExampleLogForm() [][]int {
	var S = SortableBaseN{Value: 908276523, Base: 8}
	// No clue what this is, but it's of form [A*8^n-1, B*8^n-2, ..., N*8^0]
	// Actually its   6610631453
	// Py3 confirms 0o6610631453
	var resultS []int = S.LogForm()

	var R = SortableBaseN{Value: 16756712, Base: 16}
	// '0xffafe8'	15*16^5 15*16^4 10*16^3 15*16^2 14*16^1 8*16^0
	var resultR []int = R.LogForm()

	return [][]int{resultS, resultR}
}

// Level implemented on SortableString for RadixSortable interface{}
func (S SortableString) Level(this int) interface{} {
	const epsilon = ""

	if this >= len(S) {
		return epsilon
	}
	return S[len(S)-1-this]
}

// Len implemented on SortableString for RadixSortable interface{}
func (S SortableString) Len() int {
	return len(S)
}

func (S SortableString) String() string {
	return string(S)
}

// Helper function powerOf returns the Base10 integer form of the _Base_ ^X expression
func powerOf(Base, exponent int) int {
	var ans int = 1

	for i := 0; i < exponent; i++ {
		ans = ans * Base
	}
	return ans
}

func fill(this *[]int) {
	for i := 0; i < cap(*this); i++ {
		(*this) = append(*this, i)
	}
}

// Help find the number of digits in an integer to decide rounds in RADIX
func psuedoLog10(bigNum int) int {
	ans := 0

	for Base := 10; bigNum >= 1; bigNum /= Base {
		ans++
	}

	return ans
}

func findMax(list []int) int {
	var maxSoFar int = list[0]

	for _, currentValue := range list {
		if currentValue > maxSoFar {
			maxSoFar = currentValue
		}
	}

	return maxSoFar
}

func chooseRoundsGiven(list []RadixSortable) int {
	lens := make([]int, 0, len(list))

	for _, each := range list {
		lens = append(lens, each.Len())
	}

	return findMax(lens)
}

func getOrderedKeys(unsortedPairs map[interface{}][]int) []interface{} {
	var min, max, num int

	for each := range unsortedPairs {
		num = int(each.(int))

		if num > max {
			max = num
		}
	}
	min = max
	for each := range unsortedPairs {
		num = int(each.(int))

		if num < min {
			min = num
		}
	}

	var orderedKeys []interface{} = make([]interface{}, 1+max-min) // nil may become redundant. Sorry not sorry

	for k := range unsortedPairs {
		ordinal := int(k.(int)) - min // normalization step
		orderedKeys[ordinal] = k
	}

	return orderedKeys
}

//TODO: Extend sort to strings

// RadixSort results in ascending sorted slice
// Integers supported
func RadixSort(list []RadixSortable) []RadixSortable {
	var maxRound = chooseRoundsGiven(list)
	var radixTable = make(map[interface{}][]int, 0) // key->bucket of indexes in original list
	var radixProgress = make([]int, 0, len(list))   // holds indices to become permutated during sort
	fill(&radixProgress)

	var examineValue RadixSortable
	var key interface{}

	for R := 0; R < maxRound; R++ {
		// radixTable should be EMPTY every round, which is why we advocated for map holding all round separately

		// This is the fill bucket stage
		for _, index := range radixProgress {
			examineValue = list[index]
			key = examineValue.Level(R)
			radixTable[key] = append(radixTable[key], index)
		}

		// This is the Empty Bucket Stage, important to preserve ordering
		emptyIndex := 0
		var orderedKeys []interface{} = getOrderedKeys(radixTable)

		for _, key := range orderedKeys {
			if key == nil {
				continue
			}

			for _, Value := range radixTable[key] {
				radixProgress[emptyIndex] = Value // OOB Warning
				emptyIndex++
			}

			radixTable[key] = make([]int, 0)
		}
	}
	var ans []RadixSortable = make([]RadixSortable, len(list))
	for after, before := range radixProgress {
		ans[after] = list[before]
	}

	return ans
}

// RadixCollection can be type asserted whereas a slice cannot, it may be useful to users of this API!
type RadixCollection []RadixSortable

// ConcurrentRadixSort Leverages a hypothetical principle of retained sort order and allows thread safe execution of a by-attribute based sort
func ConcurrentRadixSort(group *RadixCollection) *RadixCollection {
	// Radix Sort leverages a fundamental principle of sorting; order is retained in places of ambiguity
	// A stable sort maintains this retained ordering, whereas unstable quicksort for instance would alter the zone of ambiguous ordering among AA AA and AA
	// Radix Sort works just as well with digits as characters as attributes of any entity with ordinal categorical data or enumerated ordinal components 
	// Nil, null, None, "", _, and other such values for an attribute are acceptable and treated with descending priority
	// i.e. apple before applebees and 5 before 10
	//
	// The most common implementations of Radix Sort are Least and Most Significant Digit ordered (LSD, MSD respectively)
	// These take every attribute of a number or string and sort the collection of N elements #Rounds times
	//
	// Given a Runtime of O(#Rounds * N), where we cannot reduce N, we instead reduce #Rounds to see quick gains in performance at the cost of probably available space.
	// The worst part about radix is the #number of rounds. So, it is interesting to discuss relevant ways of reducing the rounds needed to sort any collection.
	// By far the most intuitive, or at least a simple approach is to increase the quantity of bins to decrease the #Rounds
	// For example by increasing the base system from a decimal system with 10 bins to a base with 16, 64, or 1000 bins instead
	// Such as hexadecimal, base64, or concatenating multiple characters/digits as a bin up to the point of reasonable space complexity.
	// What capacity can/should a bin have? is a relevant question to ask
	// Note: bins can be discovered prior to knowledge by guessing A*ll possibilities, or one can use just what's actually there instead.
	// The approach below aligns better with the latter method.
	// The second worst part of radix is the space-complexity. I offer no solace here
	//
	// There is a special case of RadixSort with just a single round. This is O(N) and results from the scenario where 1 >= log_Bins(N)
	// Generally the greater the available memory (and thus possible bins) and the smaller the range of numbers, the better the value of R in our favor
	//
	// Traditional Radix sorts the collection inbetween every round and may use something efficient like Key-Count Sort in place
	// This has some redundancy and does not leverage concurrency because each step is linked sequentially to the last.
	// Here I propose that we stand to gain performance from pre-sorted indices, concurrent, parallel, and lazy execution patterns
	// The trick? First, rather than order the elements, sort their original indices or say UIDs at every attribute level.
	// Second, keep in the metadata where bins begin and end, I simplify this by using maps in place of arrays at the price of space but boon of simplicity
	// My alternative here uses a single data structure which can be updated be separate processes at the same time safely
	//
	// The byRoundThenAttributeBins datastructure maintains one thing; the sortedByAttribute indices of the elements from the original collection
	// It is no great work to demonstrate how this structure can be filled by concurrent routines directed in the proper area
	// Until recently, it has been unclear whether these individual index sorts can be combined to form the desired sort pattern
	// For instance, a successful N attribute Radix Sort applies the ordering of attributes
	//					_ -> A -> B -> ... -> N where A is primary, such as the MSD
	// But if not applied cleverly, the concurrent radix sort final step to be disclosed could apply the original ordering _, between each attribute
	// Resulting in		_ -> A -> _ -> B -> _ -> ... -> _ -> N
	// Or more simply	_ -> N
	//
	// The theory behind this distinction is unclear to me, so we rely on our testsuite to confirm or deny the correctness of the following implementation
	// In order to sidestep the _ order, we rely on a sequential iteration through each attributes bins and some inspiration from Key-Count sort

	// Begin the preparation
	// prepareForN is a concise method to prepare a 2D slice inline later on, this particular implementation benefits us with spatial locality
	var prepareForN = func(rounds int) [][]int {
		var oneBigContinguousSlice = make([]int, len(*group)*rounds)
		var preparation [][]int

		var includeBefore, excludeAfter int
		for r := 0; r < rounds; r++ {
			includeBefore = r * len(*group)
			excludeAfter = includeBefore + len(*group)
			preparation[r] = oneBigContinguousSlice[includeBefore:excludeAfter]
		}

		return preparation
	}

	// trackMinAndMax doe *Exactly* what you think it does, but is suited for clean code
	var trackMinAndMax = func(any int, min, max *int) {
		if any > *max {
			*max = any
		} else if any < *min {
			*min = any
		}
	}

	// A better named function wrapping an existing API
	var getAttribute = func(precedence int, fromThis RadixSortable) int {
		return int(fromThis.Level(precedence).(int))
	}

	// A nifty in-place "sort" which will leave the indices of elements in the original collection in order and in the correct area of our main structure
	var keyCountSort = func(fromThisLevel int, intoHere *[]int) {
		// A map works wonders for counting things, but the attribute "need" to be ordered
		// Which is why this is much simpler with integers instead of any collection of sortable entities, can we get a good ordered map here?
		// Suppose we modify .Level() interface{} to .Level() int, now emptying the map is as simple as counting :)
		// We should just add a method, or better yet, a wrapper function. See getAttribute(_, RadixSortable) int
		var countThis = map[int]int{}
		var and int = getAttribute(fromThisLevel, (*group)[0]) // Min
		var thanks int = and                                   // Max

		// Count this ish up! Oh and save that max value for later
		for _, sortableValue := range *group {
			var please int = getAttribute(fromThisLevel, sortableValue)
			countThis[please]++
			trackMinAndMax(please, &and, &thanks)
		}

		// Make a running sum become the start index for all counted attribute values
		// Start with counts, end with well-placed index into slice for every key.
		// E.g. Given five 3s two 7s and one 4, get {3: 0, 4: 5, 7: 6} resulting in [(3) 3 3 3 3 (4) (7) 7] If placing in array later.
		for gratitude := 0; gratitude < thanks; gratitude++ {
			var lastSupper int = 0 // Tracks predecessor without need for two indices into countThis[]

			// Would normally just count up, but need to account for gaps in keys, like "nope, no 5's here"
			if countThis[gratitude] > 0 {
				countThis[gratitude] += lastSupper
				lastSupper = countThis[gratitude]
			}
		}

		// Now we place the indices of elements from the original group into our structure
		for each, value := range *group {
			var key = getAttribute(fromThisLevel, value)
			(*intoHere)[countThis[key]] = each
			countThis[key]++
		}
	}

	// End of helpful function declarations. Begin CRS Logix
	var rounds int = 2                    // = max length of sortables provided, more the merrier, but less is faster and we're stuck on two...; more accurate to findMax([]RadixSortable).Len()
	var groupInOrder = &RadixCollection{} // make([]RadixSortable, 0, len(group))
	var iterationsOfIndicesOrderedByAttribute = prepareForN(rounds)

	// Go collect the meta data at every level of attribute
	for r := 0; r < rounds; r++ {
		go keyCountSort(r, &iterationsOfIndicesOrderedByAttribute[r])
	}

	// Once all of that is complete... Sync possibly needed... The fun begins
	// We gonna write this part for two attributes and generalize later
	// Idea is to iterate through the meta data of each round and ":filter" through each
	// Finding a corresponding index in every round and only adding it to the final groupInOrder
	// After passing through all the rounds. I am unsure of whether this can be done for more than 2 rounds
	// And don't yet know how to implement for 2 rounds, but can do this by hand for a deck of cards using suite and ordinal index as the attributes
	// Thinking of going to round(0).attributeKey(0) and grabbing all the indices there, then iterating over the next round(1)
	// Then whenever a match is found in round(1), writing that immediately to the next stage
	// Noticed that the placement of the index can happen in the same fashion KeyCountSort places elements exactly where they need to be
	// But while this reducing redundant loops over round's slice, I suspect this strategy will raise issues of complexity or limitation after N>2 rounds
	// If not, it may need intermittent decision-making over the index of some items. Can this be achieved?
	// In the end, the result is a series of indices into the original group.
	// To put the elements in a new groupInOrder, take the element at the group[index iterated upon], and place it simply at the end of the new groupInOrder.
	// Voila! ConcurrentRadixSort

	// Other notes...
	// It appears we could "brute force" this structure, but that does not do the finesse of Radix justice.
	// It appears we can find a minimum key reminiscent of the MSD strategy
	// It appears we could track additional meta data during our keysort which would allow instant locating of the same index in other "columns", i.e. rounds
	// It appears we can do this easily with knowledge of where the start index of the preceeding attribute is,
	//		unfortunately, this is only feasible right now for a maximum of two-attributes, more becomes uncertain.

	// A solution to the above is to flow through the structure looking for group membership within each value bucket
	// We define a listener function which under the base case of recursion send along a shared channel the next value in order
	// Once the first and remaining attributes of the list are sorted, Primary, secondary, etc. sorting attributes, we then iterate over the lists in order
	// The arguments to the listener function are the indices which must be matched in each ordered key group until the base case

	var orderedIndices = make([]int, len(group)) // not going to append, but direct assign instead. Init to 0 everywhere
	var sortByMinKey = func(outsideIteration chan int, forThese []int, attributeOrdinal int, predictedIndex int) {
		// Check base cases here
		if len(forThese) == 0 {
			panic("Why is this list empty?")
		}
		else if len(forThese) == 1 {
			send(forThese[0])
		}
		if len(forThese) > 1 && attributeOrdinal + 1 == len(iterationsOfIndicesOrderedByAttribute) { 
			for _, each := forThese {
				send(each)
			}
		}

		// convert to lookup table for convenience, and provide a method for using it
		var lookoutFor = map[int]bool
		for _, this := range forThese {
			lookoutFor[this] = true
		}
		// This is an endless loop if the channel isn't closed...
		var filter() int {
			for this, ok := <-outsideIteration && ok {
				if lookoutFor[this] {
					return this
				}
			}
			return 0
		}

		// outsideIteration and decide that whenever the value changes, the group must as well
		var toThis int = filter()
		var value interface{} = getAttribute(ordinal, (*group)[toThis])
		var nextGroup []int{toThis}

		for toThis := filter() /* This loop condition is not complete */ {
			if getAttribute(ordinal, (*group)[toThis]) == value {
				nextGroup = append(nextGroup, toThis)
			} else {
				sortByMinKey(outsideIteration, nextGroup, attributeOrdinal+1, predictedIndex+len(nextGroup))
				nextGroup = make([]int, 0)
				value = getAttribute(ordinal, (*group)[toThis])
			}
		}

		// Once channel closes or forThese is emptied, quit
		return
	}

	var count {}
	var send() {
		count++
	}
	// Now initiate a set of recusive calls at level one, why do I feel that we need multiple channels..?
	var iterateToCommunicate chan int
	sortByMinKey(iterateToCommunicate, iterationsOfIndicesOrderedByAttribute[0], 0, 0)

	// Finally iterate over our structure until all the indices have been sorted
	for ordinal, attribute := range iterationsOfIndicesOrderedByAttribute {

		for index, ogIndex := range attribute {
			iterateToCommunicate <- ogIndex
		}
	}

	var mapByIndicesIntoSortedOrder = func(indices []int) {
		for unsorted, sorted := range indices {
			groupInOrder[sorted] = (*group)[unsorted]
		}
	}
	return groupInOrder
}
