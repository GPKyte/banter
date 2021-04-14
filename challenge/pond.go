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
	Fill(NumberScanner)
}

type TerrainMap struct {
	terrain Matrix
	layers  LayerMap
}

type PondSolver struct {
	model  TerrainMap
	volume int
}

type LayerTileMap map[int][]Tile

// Keys are ordered ascending from 0.
// From the problem scope we can say that Keys are greater than or equal to 0
func (orderedMap LayerTileMap) Keys() []int {
	var keys = make([]int, 0)

	for i := 0; len(keys) < len(orderedMap); i++ {
		if orderedMap[i] != nil {
			keys = append(keys, i)
		}
	}
	return keys
}

type LayerClusterMap map[int][]Cluster

type Cluster struct {
	Members        []Tile
	anyMemberLeaks bool
	height         int
}
type Tile struct {
	rowCoordinate int
	colCoordinate int
}

// FindAdjacent wraps the behavior of navigating tiles.
func FindAdjacent(t Tile) []Tile {
	neighborDirections := map[string][]Tile{
		"north": Tile{-1, 0},
		"east":  Tile{0, 1},
		"west":  Tile{0, -1},
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

	return Cluster{members, leaky}
}

// InitMatrix prepares a contiguous block of memory to serve as the underlying structure of a BasicMatrix and returns the address for the struct
func InitMatrix(height, width int) *Matrix {
	M := new(BasicMatrix)
	contiguousBlock := make([]int, height*width)

	for h := 0; h < height; h++ {
		M.twoDimCollection[h] = contiguousBlock[width*h : width*(h+1)]
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
func (m *BasicMatrix) Fill(src NumberScanner) {
	for i := 0; i < len(*m); i++ {
		for ii := 0; ii < len((*m)[i]); ii++ {
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
	return int(s) - 1
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
	From *scanner.Scanner
}

func (s StdNumberScanner) NextInt() int {
	var tok = s.From.Scan() // Parse the next token, but token is not readily usable

	if tok == scanner.EOF {
		return 0
	}

	var num int
	var err error
	if num, err = strconv.Atoi(s.From.TokenText()); err != nil {
		panic(fmt.Errorf("Only integers allowed in input. Received %v\n%v", s.From.TokenText(), err))
	}
	return num
}

func SingleSolution(input io.Reader) int {
	var problemDefinition *scanner.Scanner
	problemDefinition.Init(input)
	var get = StdNumberScanner{From: problemDefinition}

	var matrixHeight, matrixWidth int
	/* Open this comment tag to not take height and width from scanner *
	matrixHeight = get.NextInt()
	matrixWidth = get.NextInt()
	/* Close this comment tag to force set Height and Width */
	matrixHeight = 5
	matrixWidth = 5
	/**/

	// If optimization needed for input, consider buffered reading rows of matrix
	// Take this opportunity to build a layer map
	var tilesByHeight = make(LayerTileMap)
	var saveNumbers = make([]int, 0)
	for i := 0; i < matrixHeight; i++ {
		for ii := 0; ii < matrixWidth; ii++ {
			heightOfTile := get.NextInt()
			tilesByHeight[heightOfTile] = append(tilesByHeight[heightOfTile], Tile{rowCoordinate: i, colCoordinate: ii})
			saveNumbers = append(saveNumbers, heightOfTile)
		}
	}

	var matt = BasicMatrix(generateTwoDimArray(matrixHeight, matrixWidth))
	for iRow := 0; iRow < matrixHeight; iRow++ {
		for iiCol := 0; iiCol < matrixWidth; iiCol++ {
			recallHeight := saveNumbers[iRow*matrixWidth+iiCol]
			matt.Set(iRow, iiCol, recallHeight)
		}
	}

	// TODO Now that we've read input and made our two data structures, we will navigate the layers of Tiles and
	// both classify clusters of tiles as open or closed
	// We need to hold our clusters somewhere, do it by height
	var getClustersPerLayer = make(LayerClusterMap)
	// we use a BFS method to find a cluster, whilst we record the min height of the neighbors not in the cluster
	// if the min height is less than the layer, the cluster is open, otherwise the cluster is closed and
	// the top height could become the height of the lowest neighbor.
	for thisHeight := range tilesByHeight.Keys() {
		var clusterCandidates []Tile = tilesByHeight[thisHeight]
		var tomo TileVisitor
		// We upgrade connected tiles to the next layer once verified the do not reside adjacent to a lower tile
		// Because we start at the bottom and find all the neighbors of the cluster, we avoid redundant or recursive checks down slope
		for _, cc := range clusterCandidates {
			var tomoFeltDejavu = tomo.Visit(cc)

			if tomoFeltDejavu {
				continue
			}
			// Otherwise do this instead
			matt.Get(cc.rowCoordinate, cc.colCoordinate)

		}
	}

	// TODO share the top-height value among cluster members, and update it on the fly.
	// A cluster could stay in the layer or grow upward
}

// Visitor is an interface to wrap the Visit function needed now for matrix traversal
type Visitor interface {
	Visit() bool
}

// TileVisitor maintains a history of Tiles visited at every layer to assist navigation of tile neighbors
type TileVisitor struct {
	history map[Tile]bool
}

// Visit a Tile to mark that area and avoid duplicate work later, this is a version of shortcircuiting in practice.
func (tv *TileVisitor) Visit(time Tile) (beenHereBefore bool) {
	// We know that we have not
	beenHereBefore = false
	// except...
	if tv.history[time] {
		beenHereBefore = true
	}

	tv.history[time] = true

	// It may be useful to know whether we have
	return beenHereBefore
}

func generateTwoDimArray(firstDim, secondDim int) [][]int {
	var oneDim = make([]int, firstDim*secondDim)
	var twoDim = make([][]int, 0, firstDim)
	var slidingStart, slidingFinish int

	for iFirst := 0; iFirst < firstDim; iFirst++ {
		slidingStart = secondDim * iFirst
		slidingFinish = secondDim*iFirst + secondDim
		twoDim[iFirst] = oneDim[slidingStart:slidingFinish]
	}

	return twoDim
}
