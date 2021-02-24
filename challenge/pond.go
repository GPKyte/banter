package challenge

import (
	"fmt"
	"io"
	"strconv"
	"text/scanner"
)
type Matrix interface {
	Get(int, int) int
	Set(int, int, int)
	Total() int
	Equals(*Matrix) bool
	Fill([]byte)
}

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
	var visited map[Tile]bool
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
	return (*m)[row][col]
}
// Set the value for the Tile at the coordinates given to the value given
func (m *BasicMatrix) Set(row, col, value int) {
	// Check boundaries ofcourse
	(*m)[row][col] = value
}

// Fill a BasicMatric with Numbers until m full
// Error check this process by confirming src NumberScanner returns EOF after this process if desired.
func (m BasicMatrix) Fill(src NumberScanner) {
	for i := 0; i<len(m); i++ {
		for ii := 0; ii<len(m[i]); ii++ {
			m.Set(i, ii, src.NextInt())
		}
	}
}

type NumberScanner interface {
	NextInt() int
}

type AscendingNumberScanner int
func (s AscendingNumberScanner) NextInt() int {
	s += 1
	return int(s)-1
}

type DefaultNumberScanner struct {
	pos int
	src []int
}
func (s DefaultNumberScanner) NextInt() int {
	next := s.src[s.pos]
	s.pos++
	return next
}

type StdNumberScanner struct {
	from *scanner.Scanner
}
func (s StdNumberScanner) NextInt() int {
	var tok = s.from.Scan() // Parse the next token, but token is not readily usable

	if tok == scanner.EOF {
		return 0
	}

	var num int
	var err error
	if num, err = strconv.Atoi(s.from.TokenText()); err != nil {
		panic(fmt.Errorf("Expected integer, but received %v\n%v", s.from.TokenText(), err))
	}
	return num
}

func SingleSolution(input io.Reader) {
	var problemDefinition *scanner.Scanner
	problemDefinition.Init(input)
	var get = StdNumberScanner{from: problemDefinition}

	var matrixHeight, matrixWidth int
	matrixHeight = get.NextInt()
	matrixWidth = get.NextInt()

	// Take this opportunity and excuse to build a layer map
	var tilesByHeight = make(LayerTileMap)
	var saveNumbers = make([]int, 0)
	for i := 0; i < matrixHeight; i++ {
		for ii := 0; ii < matrixWidth; ii++ {
			heightOfTile := get.NextInt()
			tilesByHeight[heightOfTile] = append(tilesByHeight[heightOfTile], Tile{rowCoordinate: i, colCoordinate: ii})
			saveNumbers = append(saveNumbers, heightOfTile)
		}
	}

	var matt BasicMatrix
	matt.Fill(DefaultNumberScanner{src: saveNumbers, pos: 0})

	// TODO Now that we've read input and made our two data structures, we will navigate the layers of Tiles and
	// both classify clusters of tiles as open or closed
	// we use a BFS method to find a cluster, whilst we record the min height of the neighbors not in the cluster
	// if the min height is less than the layer, the cluster is open, otherwise the cluster is closed and
	// the top height could become the height of the lowest neighbor.

	// TODO share the top-height value among cluster members, and update it on the fly.
	// A cluster could stay in the layer or grow upward
}