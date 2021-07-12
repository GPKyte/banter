package challenge

import (
	"fmt"
)

type QuickQueue []int

func (qq *QuickQueue) enQ(me int) {
	*qq = append(*qq, me)
}
func (qq *QuickQueue) deQ() int {
	var d = (*qq)[0]
	*qq = (*qq)[1:]
	return d
}
func (qq *QuickQueue) empty() bool {
	return len(*qq) == 0
}

type QuickStack []int

func (qs *QuickStack) Push(this int) {
	*qs = append(*qs, this)
}

func (qs *QuickStack) Pop() int {
	const DefaultBubble int = -1
	var thisBubble int = DefaultBubble
	var position = len(*qs)

	if position >= 1 {
		thisBubble = (*qs)[position-1]
		*qs = (*qs)[:position-1]
	}

	return thisBubble
}

func Factorial(n int) int {
	var fax int = 1

	for i := n; i > 1; i-- {
		fax *= i
	}

	return fax
}

// Permute the series into every possible same-length ordering.
// Do to factorial nature of this operation, yield provides the output
// Prefer channel rather than holding output in a buffer of O(N!) size
func Permute(series []int) (yield chan []int) {
	// Generate all permutations of series
	var path = make(QuickStack, 0, len(series))
	yield = make(chan []int, 1024)

	// This wrapper allows recursion to kick off and then close without deadlock
	var doWorkThenCleanup = func() {
		defer close(yield)
		permuteRecursion(&path, yield, series...)
	}
	go doWorkThenCleanup()

	return yield
}
func permuteRecursion(path *QuickStack, output chan []int, partialSeries ...int) {
	if len(partialSeries) < 1 {
		var holdme []int = CopyPermutation(path)
		output <- holdme
	}

	for x, y := 0, len(partialSeries); x < y; x++ {
		path.Push(partialSeries[x])
		var removeOne = copySliceWithRemoval(partialSeries, x)
		permuteRecursion(path, output, removeOne...)
	}
	path.Pop()
}
func CopyPermutation(path *QuickStack) []int {
	var copy []int = make([]int, len(*path))

	for each := range *path {
		copy[each] = (*path)[each]
	}
	return copy
}
func copySliceWithRemoval(from []int, remove int) []int {
	copy := make([]int, 0, len(from)-1)

	for each := range from {
		if each == remove {
			continue
		}
		copy = append(copy, from[each])
	}

	return copy
}

func scramble(word string, bySlice []int) string {
	var scram = []byte(word)

	for x, y := 0, len(scram); x < y; x++ {
		scram[x] = word[bySlice[x]]
	}

	return string(scram)
}

func PermutePlus(series []int) chan string {
	var qq = QuickQueue(series)
	var qs = make(QuickStack, 0, len(series))
	var out = make(chan string)

	go func() {
		permutePlusHelper(&qq, &qs, out, len(series))
		close(out)
	}()

	return out
}

func permutePlusHelper(qq *QuickQueue, qs *QuickStack, out chan string, level int) {

	if qq.empty() {
		out <- fmt.Sprint(CopyPermutation(qs))
		return
	}

	for i := level; i > 0; i-- {
		qs.Push(qq.deQ())
		permutePlusHelper(qq, qs, out, level-1)
		qq.enQ(qs.Pop())
	}
}

// DigitModulo uses it's own sequence as the base for each digit,
// ...The base 10 decimal system is a digit sequence whose bases (and DigitModulo) are represented by the slice: [10000, 10000, 1000, 100, 10, 1, 0.1, 0.01, 0.001]
// Interpret(int) provides a sequence of digits by which to interpret a number
type DigitModulo []int

func (dm *DigitModulo) Interpret(original int) *[]int {

	var tion = make([]int, len((*dm)))
	var og = int(original)

	// Find the biggest thing that *fits*; until biggest thing is also the smallest thing possible: 1
	for biggestplace := range *dm {
		if og == 0 {
			break
		}
		if (*dm)[biggestplace] > og {
			continue
		}

		// Find how many biggest things *fit*
		var biggest = (*dm)[biggestplace]
		var thatmany int
		for howmany := 0; biggest*howmany <= og; howmany += 1 {
			thatmany = howmany // `thatmany` will lag behind `howmany` by one when og < the calculated divisor
		}

		// *Fit* that number in the right place
		var rightplace = biggestplace
		tion[rightplace] = thatmany

		// Remove that many big things
		og -= thatmany * biggest
	} // Repeat until biggest thing is also the smallest thing

	return &tion
}

func FactorialSequenceThing(n int) {
	upperbound := Factorial(n)

	// Want to see n slots, each counts where they are modulo
	slots := make(DigitModulo, n, n)
	for index := n - 1; index >= 0; index-- {
		// Decending order of factorials results in the expected order for DigitModulo
		slots[index] = Factorial(n - index)
	}

	// Given 92736:
	//  8! 7! 6! 5! 4!
	// [2  2  2  4  4] Return

	// Given 1,2,3,4,5,6
	//	001 // 1!
	//	010	// 2!
	//	011	// 2! 1!
	//	020	// 2!*2
	//	021 // 2!*2 1!
	//	100	// 3!

	var output *[]int
	for i := 0; i < upperbound; i++ {
		output = slots.Interpret(i)
		fmt.Println(*output)
	}
}

// main will test the build and traversal of the ordinal tree
func main() {
	var regularSequence = []int{0, 1, 2, 3, 4, 5, 6, 7}
	var counter int = 0
	var specialCounter DigitModulo = make([]int, len(regularSequence))
	for index := len(specialCounter) - 1; index >= 0; index-- {
		// Decending order of factorials results in the expected order for DigitModulo
		specialCounter[index] = Factorial(len(specialCounter) - index)
	}

	var allOrderings = Permute(regularSequence)
	for permutation, ok := <-allOrderings; ok; permutation, ok = <-allOrderings {
		fmt.Printf("%-10d: %3v: %3v\n", counter, permutation, *specialCounter.Interpret(counter))
		counter += 1
	}
}
