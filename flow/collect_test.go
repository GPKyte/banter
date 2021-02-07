package flow_test

import (
	"testing"

	"github.com/GPKyte/banter/flow"
)

var integerSet = []int{-9, -4, -1, -0, 1, 2, 3, 5, 9, 12, 17, 42, 222, 2321, 42111}

type numCollection []int // Implements flow.Collection
func (nc numCollection) Len() int {
	return len(nc)
}
func (nc numCollection) Get(this int) {
	return nc[this]
}

func TestAllOrNothin(t *testing.T) {
	var testCollection = numCollection(integerSet)
	var testSameLength = flow.Select(testCollection)(flow.All)
	var testZeroLength = flow.Select(testCollection)(flow.None)

	if len(testSameLength) != len(testCollection) {
		t.Logf("Expected all, but got %d out of %d", len(testSameLength), len(testCollection))
	} else if len(testZeroLength) != 0 {
		t.Logf("Expected none, but got %d out of %d", len(testZeroLength), len(testCollection))
	}
}
