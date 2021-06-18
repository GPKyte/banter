package main

import "fmt"

// Goal: Given a series of integers, build a Node-and-Pointer style tree, which
// When traversed, produces all orderings of the given series.
// OrdinalTree
type OrdinalTree struct {
	base []*OTNode // The first in the starter series, yin include the rest of series
	size int
}

type OTNode struct {
	val int
	kin []*OTNode // Descendents of this OTNode
}

func NewOrdinalTree(series []int) *OrdinalTree {
	const Basic int = -1

	var size = 1 // Determine size via factorial
	for fact := len(series); fact > 1; fact -= 1 {
		size *= fact
	}

	var base []*OTNode = buildOT(series)

	return &OrdinalTree{
		base: base,
		size: size,
	}
}

func buildOT(series []int) []*OTNode {
	var rin = make([]*OTNode, 0, len(series))

	for x, s := range series {
		rin = append(rin, &OTNode{
			val: s,
			kin: buildOT(append(series[:x], series[x+1:]...)),
		})
	}
	return rin
}

func traversalRecursionAid(stack *quickStack, output chan []int, nodes []*OTNode) {
	// Recurse down to the lowest level of Tree, adding each node value along way
	// and removing it once moving on to next of kin

	for _, eachNode := range nodes {
		stack.Push(eachNode.val)

		if len(eachNode.kin) == 0 {
			output <- []int(*stack)
		} else {
			traversalRecursionAid(stack, output, eachNode.kin)
		}

		stack.Pop()
	}
}

func (OT *OrdinalTree) Traverse() chan []int {
	var output = make(chan []int, OT.Size())
	var stack quickStack = make(quickStack, 0, len(OT.base))

	var travel = func() {
		// Starting from base perform depth first traversal on OT
		traversalRecursionAid(&stack, output, OT.base)
		close(output)
	}
	go travel()

	return output
}

func (OT *OrdinalTree) Size() int {
	return OT.size
}

type quickStack []int

func (qs *quickStack) Push(this int) {
	*qs = append(*qs, this)
}

func (qs *quickStack) Pop() int {
	const DefaultBubble int = -1
	var thisBubble int = DefaultBubble
	var position = len(*qs)

	if position >= 1 {
		thisBubble = (*qs)[position-1]
		*qs = (*qs)[:position-1]
	}

	return thisBubble
}

type quickQueue []*OTNode

func (qq *quickQueue) EnQ(this *OTNode) {
	*qq = append(*qq, this)
}

func (qq *quickQueue) DeQ() *OTNode {
	var DefaultBubble *OTNode = &OTNode{}
	var thisBubble *OTNode = DefaultBubble

	if len(*qq) > 1 {
		thisBubble = (*qq)[0]
		*qq = (*qq)[1:]
	}

	return thisBubble
}

func factorial(n int) int {
	var fax int = 1

	for i := n; i > 1; i-- {
		fax *= i
	}

	return fax
}

// Permute the series into every possible same-length ordering.
// Do to factorial nature of this operation, yield provides the output
// Prefer channel rather than holding output in a buffer of O(N!) size
func permute(series []int) (yield chan []int) {
	// Generate all permutations of series
	var path = make(quickStack, 0, len(series))
	yield = make(chan []int, 1024)

	// This wrapper allows recursion to kick off and then close without deadlock
	var doWorkThenCleanup = func() {
		defer close(yield)
		permuteRecursion(&path, yield, series...)
	}
	go doWorkThenCleanup()

	return yield
}
func permuteRecursion(path *quickStack, output chan []int, partialSeries ...int) {
	if len(partialSeries) < 1 {
		var holdme []int = copyPermutation(path)
		output <- holdme
	}

	for x, y := 0, len(partialSeries); x < y; x++ {
		path.Push(partialSeries[x])
		var removeOne = copySliceWithRemoval(partialSeries, x)
		permuteRecursion(path, output, removeOne...)
	}
	path.Pop()
}
func copyPermutation(path *quickStack) []int {
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

type digitmodulo []int

func (dm *digitmodulo) Interpret(original int) *[]int {

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
}
// main will test the build and traversal of the ordinal tree
func main() {
	var regularSequence = []int{0, 1, 2, 3, 4, 5, 6, 7}

	var allOrderings = permute(regularSequence)
	for permutation, ok := <-allOrderings; ok; permutation, ok = <-allOrderings {
		fmt.Println(scramble("permuted", permutation))
	}
}
