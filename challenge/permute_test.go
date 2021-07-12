package challenge_test

import (
	"fmt"
	"strings"
	"testing"
	
	"github.com/GPKyte/banter/challenge"
)

func TestEfficientPermutation(t *testing.T) {
	var series = []int{0, 1, 2, 3, 4, 5}
	var bufferedOutput strings.Builder
	var outstream chan string = permutePlus(series)
	var counter int = 0
	var bufferLimit = 1000

	for permutation := range outstream {
		fmt.Fprintln(&bufferedOutput, permutation)

		if counter >= bufferLimit {
			counter = 0
			t.Log(bufferedOutput.String())
			bufferedOutput.Reset()
		}
	}

	t.Log(bufferedOutput.String())
	bufferedOutput.Reset()
}

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
)

func TestPermuteSizeChanges(t *testing.T) {
	var series = []int{1, 2, 3, 4, 5}
	var allOrderings = challenge.Permute(series)
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
	var allOrderings = challenge.Permute(series)
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
	var stack = make(challenge.QuickStack, 0, 20)

	stack.Push(9)
	stack.Push(9)
	stack.Pop()
	stack.Pop()

	stack.Push(1) // .1 2 3 4 5
	stack.Push(2) // .2 3 4 5
	stack.Push(3) // .3 4 5
	stack.Push(4) // .4 5
	stack.Push(5) // .5
	t.Log(challenge.CopyPermutation(&stack))
	stack.Pop()   //
	stack.Pop()   // 5
	stack.Push(5) // 4 .5
	stack.Push(4) // .4
	t.Log(challenge.CopyPermutation(&stack))
	stack.Pop()   //
	stack.Pop()   // 4
	stack.Pop()   // 4 5
	stack.Push(4) // 3 .4 5
	stack.Push(3) // .3 5
	stack.Push(5) // .5
	t.Log(challenge.CopyPermutation(&stack))
	stack.Pop()   //
	stack.Pop()   // 5
	stack.Push(5) // 3 .5
	stack.Push(3) // .3
	t.Log(challenge.CopyPermutation(&stack))
	stack.Pop()   //
	stack.Pop()   // 3
	stack.Pop()   // 3 5
	stack.Pop()   // 3 4 5
	stack.Push(3) // 2 .3 4 5
	stack.Push(2) // .2 4 5
	stack.Push(4) // .4 5
	stack.Push(5) // .5
	t.Log(challenge.CopyPermutation(&stack))
	stack.Pop()   //
	stack.Pop()   // 5
	stack.Push(5) // 4 .5
	stack.Push(4) // .4
	t.Log(challenge.CopyPermutation(&stack))
	stack.Pop()   //
	stack.Pop()   // 4
	stack.Pop()   // 4 5
	stack.Push(4) // 2 .4 5
	stack.Push(2) // .2 5
	stack.Push(5) // .5
	t.Log(challenge.CopyPermutation(&stack))
	stack.Pop() //
	stack.Pop() // 5
}

func TestFactSequence(t *testing.T) {
	t.Fail()

	challenge.FactorialSequenceThing(5)
}

func TestModuloSliceInterpreter(t *testing.T) {
	var expectations = map[int][]int{
		555: {0, 4, 3, 0, 1, 1}, // 0 480 72 0 2 1
		1:   {0, 0, 0, 0, 0, 1},
		800: {1, 0, 3, 1, 1, 0}, // 720 0 72 6 2 0
	}

	var unequal = func(pea []int, bea []int) bool {
		if len(pea) != len(bea) {
			return true
		}

		for index := len(pea) - 1; index > 0; index-- {
			if pea[index] != bea[index] {
				return true
			}
		}

		return false
	}
	if !unequal([]int{1, 2, 4}, []int{1, 2, 3}) || unequal([]int{4, 4}, []int{4, 4}) {
		t.Fail()
		t.Log("Inequality matches incorrect.")
	}

	var dm = challenge.DigitModulo([]int{720, 120, 24, 6, 2, 1})

	for eachOf := range expectations {
		if reality := dm.Interpret(eachOf); unequal(expectations[eachOf], *reality) {
			t.Fail()
			t.Logf("Did not meet expectation because %v is not %v", reality, expectations[eachOf])
		}
	}

	dm = challenge.DigitModulo([]int{100, 10, 1})
	expectations = map[int][]int{
		555: {5, 5, 5},
		1:   {0, 0, 1},
		800: {8, 0, 0},
	}
	for eachOf := range expectations {
		if reality := dm.Interpret(eachOf); unequal(expectations[eachOf], *reality) {
			t.Fail()
			t.Logf("Did not meet expectation because %v is not %v", reality, expectations[eachOf])
		}
	}

}

// Goal: Learn how to build operator support for specific type, or specific interface
