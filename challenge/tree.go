package main

import (
	"fmt"
)

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

	if position > 1 {
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

// main will test the build and traversal of the ordinal tree
func main() {
	var regularSequence = []int{1, 2, 3, 4, 5}
	var otter *OrdinalTree = NewOrdinalTree(regularSequence)
	var allOrderings = otter.Traverse()

	for permutation, ok := <-allOrderings; ok; permutation, ok = <-allOrderings {
		fmt.Println(permutation)
	}

	fmt.Println(*otter)
}
