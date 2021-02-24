package challenge_test

import (
	"strings"
	"testing"
	"text/scanner"

	"github.com/GPKyte/banter/challenge"
)

type TestNumberScanner struct {
	pos int
	src []int
}

func (s TestNumberScanner) NextInt() int {
	next := s.src[s.pos]
	s.pos++
	return next
}

const BigEndianMinInteger = int32(1 << 30) // 32 bit, signed 1111...1{28} equivalent positive integer 0111...1{28}
const BigEndianMaxInteger = int32(1<<30 - 1)

func TestInitAndFillMatrix(t *testing.T) {
	var height, width int = 4, 4
	goal := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

	MattG := challenge.InitMatrix(height, width)
	(*MattG).Fill(TestNumberScanner{0, goal})
	var expectedTotal int = 17 * 8

	if (*MattG).Total() != expectedTotal {
		t.Fail()
		t.Logf("TestI...FillMatrices: Expected %v, but was %v.", expectedTotal != (*MattG).Total())
		t.Log(*MattG)
	}
	if (*MattG).Get(2, 4) != goal[2*4] {
		t.Fail()
		t.Logf("TestI...FillMatrices: Expected %v, But got %v", goal[2*4], (*MattG).Get(2, 4))
		t.Log(*MattG)
	}
}

func TestBlackBoxKnownResults(t *testing.T) {
	// Set up some maps
	// Calculate their counterpart
	// Compare them using Matrix Subtraction, a feature we could certainly improve on LATER
	// Return total totalSum of water
	MattBefore := *challenge.InitMatrix(4, 3)
	MattBefore.Fill(
		"3 3 3",
		"3 1 3",
		"3 1 3",
		"3 3 3",
	)
	MattAfter := *challenge.InitMatrix(4, 3)
	MattAfter.Fill(
		"3 3 3",
		"3 3 3",
		"3 3 3",
		"3 3 3",
	)
	var totalSumDifference int = MattAfter.Total() - MattBefore.Total()

	TravisAfter := MattAfter.Traverse()
	TravisBefore := MattBefore.Traverse()
	traversals := MattAfter.Size()
	if traversals != MattBefore.Size() {
		t.Log("Before and After Matrices are incongruent because the sizes differ. See below:")
		t.Log(MattBefore)
		t.Log(MattAfter)
	}

	manualSum := 0
	expectedVolumeOrTotalSumOfDifference := (3 /*after*/ - 1 /*before*/) * 2 /*times*/
	for i := 0; i < traversals; i++ {
		manualSum += TravisAfter.Now() - TravisBefore.Now()
		TravisAfter.Next()
		TravisBefore.Next()
	}

	if manualSum != expectedVolumeOrTotalSumOfDifference || manualSum != totalSumDifference {
		t.Fail()
		t.Log("Expected %v, but found %v and %v instead", expectedVolumeOrTotalSumOfDifference, manualSum, totalSumDifference)
	}
}

func TestSingleSolution(t *testing.T) {
	var problemDefinition = strings.NewReader(`7 7
	1 7 7 7 7 7 3
	4 1 1 1 2 1 4
	3 1 1 1 2 1 5
	5 1 1 2 2 1 7
	5 2 8 8 1 1 8
	3 1 1 4 1 1 4
	5 5 5 5 5 5 8`)

	var filledMatrixDefinition = strings.NewReader(
		`1 7 7 7 7 7 3
4 3 3 3 3 3 4
3 3 3 3 3 3 5
5 3 3 3 3 3 7
5 3 8 8 3 3 8
3 3 3 4 3 3 4
5 5 5 5 5 5 8`)
	var s scanner.Scanner
	s.Init(filledMatrixDefinition)

	var matt = challenge.InitMatrix(7, 7)
	(*matt).Fill(challenge.StdNumberScanner{From: &s})

	solution := challenge.SingleSolution(problemDefinition)

	t.Fail()
	t.Log(matt)
	t.Log(solution)
	t.Log(problemDefinition)
}

type OrderedMap interface {
	Keys() []interface{}
}

// Store the set of numbers used to fill map as keys and their input-ordered indices in lists at every value
type EqualityListOrderedMap struct {
	size        func()
	Unordered   map[int][]int
	orderedKeys []int
}

// As keys are based on values, we could sort through min..max checking for each entry
// We could append successful hits into the slice we return
// We could call sort on the numbers using some other package
// Because it doesn't matter until it does, we choose the less interfacing option
func (unsorted *EqualityListOrderedMap) Keys() []interface{} {
	var keys = make([]interface{}, 0)
	var keyCount = len(unsorted.Unordered)

	// TODO: Could implement caching by storing the sorted keys and a "cacheExists" bool "cacheUpToDate" bool

	for i := 0; len(keys) < keyCount && i > 0; i++ {
		if unsorted.Unordered[i] != nil {
			keys = append(keys, int(i))
		}
	}

	return keys
}

func TestMetaLayeringStrategy(t *testing.T) {
	// Take a stream of numbers, use an exported process from elsewhere to convert the input into shared value key pairs
	// The indices will be tracked in place of the numbers, the values copied from the numbers with be keys, the key set will be ordered
	// We will call this service EqualityListOrderedMap
	// Layers are all height greater than 0
	var someNaturalNumbers = []int{1, 1, 1, 2, 2, 3, 4, 4, 4, 4, 4, 5, 5, 5, 6, 7, 8, 9, 9, 9, 11, 11, 1, 2, 3, 4, 3, 2, 1}
	var setOfSameNaturalNumbers = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 11}
	var groupIndicesByTheirSharedElemValues = new(EqualityListOrderedMap)

	// Init EqualityListOrderedMap()
	for each, item := range someNaturalNumbers {
		var mayBeNilList = groupIndicesByTheirSharedElemValues.Unordered[item]
		var notNilList []int

		if len(mayBeNilList) == 0 {
			notNilList = make([]int, 0)
			groupIndicesByTheirSharedElemValues.Unordered[item] = notNilList
		}
		groupIndicesByTheirSharedElemValues.Unordered[item] = append(groupIndicesByTheirSharedElemValues.Unordered[item], each)
	}

	if len(setOfSameNaturalNumbers) != len(groupIndicesByTheirSharedElemValues.Keys()) {
		t.Fail()
		t.Log(groupIndicesByTheirSharedElemValues.Keys())
		t.Log(groupIndicesByTheirSharedElemValues)
	}
}
