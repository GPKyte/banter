package flow_test

import (
	"testing"

	"github.com/GPKyte/banter/flow"
)

var integerSet = []int{-9, -4, -1, -0, 1, 2, 3, 5, 9, 12, 17, 42, 222, 2321, 42111}
var words = []string{"app", "battle", "can", "drumstick", "erlang", "fingy", "gravitate", "happy", "intricate", "justified", "king", "lump", "monkey", "nebulous", "om nom nom", "penelope", "quirky"}

type numCollection []int // Implements flow.Collection
func (nc numCollection) Len() int {
	return len(nc)
}
func (nc numCollection) Get(this int) interface{} {
	return nc[this]
}

type wordCollection []string // Implements flow.Collection
func (wc wordCollection) Len() int {
	return len(wc)
}
func (wc wordCollection) Get(this int) interface{} {
	return wc[this]
}

func TestAllOrNothin(t *testing.T) {
	var testCollection = flow.Collection(numCollection(integerSet))
	var testSameLength = flow.Select(testCollection)(flow.All)
	var testZeroLength = flow.Select(testCollection)(flow.None)

	if testSameLength.Len() != testCollection.Len() {
		t.Logf("Expected all, but got %d out of %d", testSameLength.Len(), testCollection.Len())
		t.Log(testSameLength)
		t.Fail()
	}
	if testZeroLength.Len() != 0 {
		t.Logf("Expected none, but got %d out of %d", testZeroLength.Len(), testCollection.Len())
		t.Log(testZeroLength)
		t.Fail()
	}
	if t.Failed() {
		t.Log("The following flow.Collection was used in both tests:")
		t.Log(testCollection)
	}
}

func TestShortie(t *testing.T) {
	var testWords = flow.Collection(wordCollection(words))
	var testNonZeroLesserLength = flow.Select(testWords)(flow.Short)

	if testNonZeroLesserLength.Len() == 0 || testNonZeroLesserLength.Len() >= testWords.Len() {
		t.Fail()
		t.Logf("Expected some of total, got %d out of %d", testNonZeroLesserLength.Len(), testWords.Len())
		t.Log(testNonZeroLesserLength)
		t.Fail()
	}

	for loop := 0; loop < testNonZeroLesserLength.Len(); loop++ {
		var word = string(testNonZeroLesserLength.Get(loop).(string))

		if len(word) > 5 {
			t.Fail()
			t.Logf("The flow.Short function should have filtered away long strings")
			t.Log("Yet, we found a string that exceeding expectations: ", word)
		}
	}

	if !t.Failed() {
		t.Log(testNonZeroLesserLength)
	}
}
