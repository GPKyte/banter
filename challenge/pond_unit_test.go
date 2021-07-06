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
