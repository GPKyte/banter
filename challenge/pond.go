package challenge

import (
	"io"

)

type TerrainMap struct {
	terrain Matrix
	layers LayerMap
}

type PondSolver struct {
	model TerrainMap
	volume int
}

type LayerTileMap map[int][]Tile
type LayerClusterMap map[int][]Cluster
type Cluster struct {
	Members []Tile
	anyMemberLeaks bool
	height int
}
type Tile struct {
	rowCoordinate int
	colCoordinate int
}

// FindAdjacent wraps the behavior of navigating tiles.
func FindAdjacent(t Tile) []Tile {
	neighborDirections := map[string][]Tile {
		"north": Tile{-1, 0},
		"east": Tile{0, 1},
		"west": Tile{0, -1},
		"south": Tile{1, 0},
	}

	neighbors := make([]Tile, 0, 4)
	for _, each := range neighborDirections {
		row := t.rowCoordinate + each.rowCoordinate
		col := t.colCoordinate + each.colCoordinate
		neighbors = append(neigbors, Tile{row, col})
	}

	return neighbors
}

func FindCluster(t Tile, lookup TerrainMap) Cluster {
	var leaky bool = false
	var members = make([]Tile, 0)
	var height = func(t Tile) int {
		return lookup.model.Get(t.rowCoordinate, t.col.Coordinate)
	}

	// Traverse neighbors of tile
	// How will we avoid revisiting?
	// A small map we discard later could do the trick
	var visited = map[Tile]bool
	var queue = t.FindNeighbors()
	var sameHeight = height(t)

	for len(queue) > 0 {
		var qt = queue[0]
		if visited[qt] {
			continue
		}
		var this = height(qt)

		if this < sameHeight {
			leaky = true
		} else if this == sameHeight {
			queue = queue.append(qt.FindNeighbors())
			members = append(members, qt)
		}

		visited[qt] = true
		queue = dropFirst(queue)
	}

	return Cluster{members, leaky, }	
}

type Matrix interface {
	Get(int, int) int
	Set(int, int, int)
}
// InitMatrix prepares a contiguous block of memory to serve as the underlying structure of a BasicMatrix and returns the address for the struct
func InitMatrix(height, width int) *Matrix {
	M := new(BasicMatrix)
	contiguousBlock := make([]int, height * width)

	for h := 0; h<height; h++ {
		M.twoDimCollection[h] = contiguousBlock[width*h:width*(h+1)]		
	}
	return M
}

// BasicMatrix is the Barebones implementation for trying the pond problem
type BasicMatrix [][]int

// Get returns a value from the Tile at the coordinates given
func (m *BasicMatrix) Get(row, col int) (value int) {
	// Check boundaries ofcourse
	return m[row][col]
}
// Set the value for the Tile at the coordinates given to the value given
func (m *BasicMatrix) Set(row, col, value int) {
	// Check boundaries ofcourse
	m[row][col] = value
}

// Fill a BasicMatric with Numbers until m full
// Error check this process by confirming src NumberReader returns EOF after this process if desired.
func (m BasicMatrix) Fill(src NumberReader) {
	for i := 0; i<len(m); i++ {
		for ii := 0; ii<len(m[i]); ii++ {
			m.Set(i, ii, src.ReadNextInt())
		}
	}
}

type NumberReader interface {
	ReadNextInt()
}
type AscendingNumberReader int
type TestNumberReader []int
type IONumberReader io.Reader


func (r *AscendingNumberReader) ReadNextInt() int {
	*r += 1
	return *r-1
}

func (r *TestNumberReader) ReadNextInt() (int, error) {
	var next int = r[0]
	*r = dropFirst(r)

	if EOF := recover(); EOF != nil {
		return 0, error("EOF")
	}

	return next, nil
}

func dropFirst(maybeEmptySlice *[]int) []int {
	if len(*maybeEmptySlice) <= 1 {
		return []int{}
	}
	return (*maybeEmptySlice)[1:len(*maybeEmptySlice)]
}

func (r *IONumberReader) ReadNextInt() int {
	// We need to read byte by byte and after encountering 0-9 anything other than 0-9 indicates end of current int
	var place = make([]byte, 0)
	var wordHolder = make([]byte, 0)

	r.Read(place)
	for each, b := range place {

		switch b {
		case byte(0):
		case byte(1):
		case byte(2):
		case byte(3):
		case byte(4):
		case byte(5):
		case byte(6):
		case byte(7):
		case byte(8):
		case byte(9):

		}
		
	}
}

