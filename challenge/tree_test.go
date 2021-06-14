package main

import (
	"testing"
)

func TestBigPermutation(t *testing.T) {
	var series = []int{9, 1, 2, 4, 7, 9, 0, 24, 645, 3675, 343, 356, 23, 267, 431}
	var yieldPermutations = permute(series)
	var expectedOutput int = factorial(len(series))
	var counter int = 0

	var handleOutOfRange = func() {
		if something := recover(); something != nil {
			t.Log(something)
		}
	}
	defer handleOutOfRange()

	for p, ok := <-yieldPermutations; ok; p, ok = <-yieldPermutations {
		counter++
		t.Log(p)
	}

	if counter != expectedOutput {
		t.Fail()
		t.Logf("Expected %d, Generated %d", expectedOutput, counter)
	}
}

func TestPermuteSizeChanges(t *testing.T) {
	var series = []int{1, 2, 3, 4, 5}
	var allOrderings = permute(series)
	var counter int = 0

	for permutation, ok := <-allOrderings; ok; permutation, ok = <-allOrderings {
		counter += 1

		if len(permutation) != len(series) {
			t.Logf("Permutation #%d failed size check, debug from here.", counter)
			t.Logf(" - Expected #%d to be length %d, was %d and looks like %v.", counter, len(series), len(permutation), permutation)
			t.Fail()
		}
	}
}

func sum(totalThis []int) int {
	var total int = 0

	for _, x := range totalThis {
		total += x
	}

	return total
}
func TestPermuteYieldsWrongNumbers(t *testing.T) {
	var series = []int{1, 2, 3, 4, 5}
	var seriesSum = sum(series)
	var allOrderings = permute(series)
	var counter int = 0

	for permutation, ok := <-allOrderings; ok; permutation, ok = <-allOrderings {
		counter += 1

		if sum(permutation) != seriesSum {
			t.Logf("Permutation #%d is not an accurate permuation because the numbers differ from the series", counter)
			t.Logf(" - Expected #%d to have the numbers %v, but had %v instead.", counter, series, permutation)
			t.Fail()
		}
	}
}

func TestStack(t *testing.T) {
	var stack = make(quickStack, 0, 20)

	stack.Push(9)
	stack.Push(9)
	stack.Pop()
	stack.Pop()

	stack.Push(1) // .1 2 3 4 5
	stack.Push(2) // .2 3 4 5
	stack.Push(3) // .3 4 5
	stack.Push(4) // .4 5
	stack.Push(5) // .5
	t.Log(copyPermutation(&stack))
	stack.Pop()   //
	stack.Pop()   // 5
	stack.Push(5) // 4 .5
	stack.Push(4) // .4
	t.Log(copyPermutation(&stack))
	stack.Pop()   //
	stack.Pop()   // 4
	stack.Pop()   // 4 5
	stack.Push(4) // 3 .4 5
	stack.Push(3) // .3 5
	stack.Push(5) // .5
	t.Log(copyPermutation(&stack))
	stack.Pop()   //
	stack.Pop()   // 5
	stack.Push(5) // 3 .5
	stack.Push(3) // .3
	t.Log(copyPermutation(&stack))
	stack.Pop()   //
	stack.Pop()   // 3
	stack.Pop()   // 3 5
	stack.Pop()   // 3 4 5
	stack.Push(3) // 2 .3 4 5
	stack.Push(2) // .2 4 5
	stack.Push(4) // .4 5
	stack.Push(5) // .5
	t.Log(copyPermutation(&stack))
	stack.Pop()   //
	stack.Pop()   // 5
	stack.Push(5) // 4 .5
	stack.Push(4) // .4
	t.Log(copyPermutation(&stack))
	stack.Pop()   //
	stack.Pop()   // 4
	stack.Pop()   // 4 5
	stack.Push(4) // 2 .4 5
	stack.Push(2) // .2 5
	stack.Push(5) // .5
	t.Log(copyPermutation(&stack))
	stack.Pop() //
	stack.Pop() // 5
}
