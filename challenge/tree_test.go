package main

import (
	"testing"
)

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
