package challenge

import (
	"testing"
)

func TestClusterTogether(t *testing.T) {
	t.Fail()
	var src = [][]int{
		{3, 4, 4, 5, 3, 3}, //		{e, i, i, j, f, f,},
		{5, 2, 2, 3, 1, 3}, //		{k, c, c, g,_a, f,},
		{4, 1, 3, 3, 2, 6}, //		{h, b, g, g, d, l,},
	}
	var batt Matrix = (*BasicMatrix)(&src)
	var tiles = make([]Tile, 0, len(src)*len(src[0]))
	for dimOne, v := range src {
		for dimTwo := range v {
			tiles = append(tiles, Tile{rowCoordinate: dimOne, colCoordinate: dimTwo})
		}
	}

	// Assert clusters a..l found
	t.Log(clusterTogether(tiles, batt))
	t.Log(clusterTogether([]Tile{
		{rowCoordinate: 1, colCoordinate: 1},
		{rowCoordinate: 1, colCoordinate: 2},
	}, batt))
}

func TestLoopTillQEmpty(t *testing.T) {
	var example = []Tile{
		{rowCoordinate: 1, colCoordinate: 1},
		{rowCoordinate: 1, colCoordinate: 2},
		{rowCoordinate: 2, colCoordinate: 1},
		{rowCoordinate: 2, colCoordinate: 2},
	}
	var queueStorage = make([]Tile, 0, len(example))
	var q Q = Q{fifo: &queueStorage}

	q.wait(example[0])
	q.wait(example[1])
	q.wait(example[2])
	q.wait(example[3])

	var countMustBeFour = 0

	for !q.empty() {
		countMustBeFour += 1
		t.Log(q.serve())
	}

	if countMustBeFour != 4 {
		t.Errorf("Count was %d", countMustBeFour)
	}
}

func TestMiddle(t *testing.T) {
	var middle = func(a, b, c int) int {
		for combo := range Permute([]int{a, b, c}) {
			if combo[0] >= combo[1] && combo[1] > combo[2] {
				return combo[1]
			}
		}
		return a
	}
	var tests = map[int][]int{
		7: {7, 9, 6}, // a
		2: {1, 2, 3}, // b
		3: {8, 2, 3}, // c
		5: {5, 5, 5}, // a,b,c
	}

	for mid, param := range tests {
		result := middle(param[0], param[1], param[2])
		if mid != result {
			t.Errorf("Found %d in middle, but expected %d", result, mid)
		}
	}
}
