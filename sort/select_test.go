package sort_test

import (
	"fmt"
	"testing"

	"github.com/GPKyte/banter/sort"
)

var input = []int{7, 6, 4, 2, 5, 7, 8, 4, 2, 5, 6, 8, 9, 4, 2, 5, 8, 6, 7, 7, 9, 0, 4, 3, 2, 5, 4, 6, 5, 3, 4, 2, 3, 1, 3, 4, 7, 8, 9}

func TestCountTo3(t *testing.T) {
	var outputBuffer = []int{}
	var sortSrc = sort.ReadFirstToN(input, 3)

	for next, ok := <-sortSrc; ok; next, ok = <-sortSrc {
		outputBuffer = append(outputBuffer, next)
	}
	fmt.Println("out: ", outputBuffer)
}

func TestNoCommonMembers(t *testing.T) {
	var betterBeEmpty []int = sort.SelectCommonMembers([][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	})

	if len(betterBeEmpty) > 0 {
		t.Logf("Expected No Common Members, found %d: %v", len(betterBeEmpty), betterBeEmpty)
		t.FailNow()
	}
}

func TestSomeCommonMembers(t *testing.T) {
	var justACouple []int = sort.SelectCommonMembers([][]int{
		{1, 2, 3},
		{1, 9, 3},
		{1, 2, 3},
	})

	if len(justACouple) != 2 {
		t.Logf("Expected More or Less Members, found %d: %v", len(justACouple), justACouple)
		t.FailNow()
	}
}

func TestAllCommonMembers(t *testing.T) {
	var everybody []int = sort.SelectCommonMembers([][]int{
		{1, 2, 3},
		{1, 2, 3},
		{1, 2, 3},
	})

	if len(everybody) != 3 {
		t.Logf("Expected Only Common Members, found %d: %v", len(everybody), everybody)
		t.Fail()
	}

	var justNobody []int = sort.SelectCommonMembers([][]int{
		{0},
		{0},
		{0},
	})

	if len(justNobody) != 1 {
		t.Logf("Expected Only Common Members, found %d: %v", len(justNobody), justNobody)
		t.FailNow()
	}
}
