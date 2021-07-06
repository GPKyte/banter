package challenge

import (
	"fmt"
	"io"
	"strconv"
	"text/scanner"
)

const OutOfBoundsHeight = 0

type Matrix interface {
	Get(int, int) int
	Set(int, int, int)
	Total() int
	Equals(Matrix) bool
	Fill(NumberScanner)
}

type TerrainMap struct {
	terrain Matrix
	layers  LayerTileMap
}

type PondSolver struct {
	model  TerrainMap
	volume int
}

type LayerTileMap map[int][]Tile

// MaxMapSize is an arbitrary limit
const MaxMapSize = 99999

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

// What would finding the cluster of an area look like?
// Say I look at Tile_0 from the input reel and it has two neighbors, Tile_01, Tile_02.
// I found the neighbors via a helper function "FindAdjacent".
// Tile_0 is a border tile in the upperleft corner of a Matrix or Terrain Map in this example.
// Tile_0 is first into the Q, followed by 01 and 02, 01's neighbors 01a 01b 01c then 02's neighbors 02a 01b 02c
// Exploring and revisiting these Tiles in a certain order via a Q affords well-known benefits from the BFS

// Q.ueue of Tiles found via the BFS method on the prestructured data (e.g. Matrix, graph)
type Q struct {
	fifo *[]Tile
}

// Serve will bring the next Tile out from a wait state
func (q *Q) serve() Tile {
	const outOfRange = -MaxMapSize
	var errorTile = Tile{outOfRange, outOfRange}
	var lenQ = len(*(*q).fifo) // This number appears several times locally

	// cannot return a Tile, error condition
	if lenQ <= 0 {
		return errorTile
	}

	var beingServed Tile = (*q.fifo)[0]
	*q.fifo = (*q.fifo)[1:]

	return beingServed
}

// Wait will enqueue
func (q *Q) wait(here Tile) {
	// Will this lack of pointing to references lead to a frontline issue?
	*q.fifo = append(*q.fifo, here)
}

type LayerClusterMap map[int][]Cluster

type Cluster struct {
	Members        []Tile
	anyMemberLeaks bool
	height         int
}

// retainsRainWater is a nicety relevant for the original problem definition which reads well in an If statement.
func (c *Cluster) retainsRainWater() bool {
	// Will need to find whether anyMemberLeaks by searching the perimeter beyond the cluster
	// The cluster does not have access to this information so it must be decided elsewhere
	return !c.anyMemberLeaks
}

type Tile struct {
	rowCoordinate int
	colCoordinate int
}

// FindAdjacent wraps the behavior of navigating tiles.
func FindAdjacent(t Tile) []Tile {
	neighborDirections := map[string]Tile{
		"north": {-1, 0},
		"east":  {0, 1},
		"west":  {0, -1},
		"south": {1, 0},
	}

	neighbors := make([]Tile, 0, 4)
	for _, each := range neighborDirections {
		row := t.rowCoordinate + each.rowCoordinate
		col := t.colCoordinate + each.colCoordinate
		neighbors = append(neighbors, Tile{row, col})
	}

	return neighbors
}

func FindCluster(t Tile, lookup Matrix) Cluster {
	var maybeLeaky bool = false
	var members = make([]Tile, 0)
	var height = func(t Tile) int {
		return lookup.Get(t.rowCoordinate, t.colCoordinate)
	}

	// Traverse neighbors of tile
	// How will we avoid revisiting?
	// A small map we discard later could do the trick
	var alreadyVisited = map[Tile]bool{}
	var queue = FindAdjacent(t)
	var sameHeight = height(t)

	// Here we explore immediate neighbors to the tile
	// Same height tiles are added to the Queue to be considered in the cluster.
	for len(queue) > 0 {
		var qt = queue[0]
		if alreadyVisited[qt] {
			continue
		}

		var this = height(qt)
		if this < sameHeight {
			maybeLeaky = true
		} else if this == sameHeight {
			queue = append(queue, FindAdjacent(qt)...)
			members = append(members, qt)
		}

		alreadyVisited[qt] = true
		queue = queue[1:]
	}

	return Cluster{members, maybeLeaky, sameHeight}
}

// InitMatrix prepares a contiguous block of memory
// to serve as the underlying structure of a BasicMatrix
// and returns the address for the struct
func InitMatrix(height, width int) *BasicMatrix {
	m := new(BasicMatrix)
	contiguousBlock := make([]int, height*width)

	for h := 0; h < height; h++ {
		(*m) = append((*m), contiguousBlock[width*h:width*(h+1)])
	}
	return m
}

// BasicMatrix is the Barebones implementation for trying the pond problem
type BasicMatrix [][]int

// Get returns a value from the Tile at the coordinates given
func (m *BasicMatrix) Get(row, col int) (value int) {
	// Check boundaries. Beyond the Matrix is defined as an endless expanse at the unchanging height of 0.
	if row < 0 || col < 0 || row >= len(*m) || len(*m) <= 0 || col >= len((*m)[0]) {
		return OutOfBoundsHeight
	}

	return (*m)[row][col]
}

// Set the value for the Tile at the coordinates given to the value given
func (m *BasicMatrix) Set(row, col, value int) {
	// Check boundaries
	if row < len(*m) && col < len((*m)[0]) && row >= 0 && col >= 0 {
		(*m)[row][col] = value
	}
}

// Fill a BasicMatrix with Numbers until m full
// Error check this process by confirming src NumberScanner returns EOF after this process if desired.
func (m *BasicMatrix) Fill(src NumberScanner) {
	for i := 0; i < len(*m); i++ {
		for ii := 0; ii < len((*m)[i]); ii++ {
			m.Set(i, ii, src.NextInt())
		}
	}
}

func (m *BasicMatrix) Total() int {
	var sum int
	for i := 0; i < len(*m); i++ {
		for ii := 0; ii < len((*m)[i]); ii++ {
			sum += m.Get(i, ii)
		}
	}
	return sum
}

// Equals referenes deep values
func (one *BasicMatrix) Equals(another Matrix) bool {
	var other BasicMatrix = *another.(*BasicMatrix) // I love this line -GK

	if len(*one) != len(other) {
		return false
	}
	// Use the one BasicMatrix and the API of Matrix to find equality by catching OOB error over the other Matrix
	for iCountdown := len(*one); iCountdown > 0; iCountdown-- {
		if len((*one)[iCountdown]) != len(other[iCountdown]) {
			return false
		}
		for iiCountdown := len((*one)[iCountdown]); iiCountdown > 0; iiCountdown-- {
			this := (*one).Get(iCountdown, iiCountdown)
			that := (other).Get(iCountdown, iiCountdown)

			if this != that {
				return false
			}
		}
	}
	// One Matrix and the other Matrix must share the same size, shape, and values.
	return true // They are Equivalent
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

// ClusterTogether connectedTiles from the given Tiles based on adjacency
// Clusters are Useful for addressing several connected Tiles
func clusterTogether(unClustered []Tile, src Matrix) []Cluster {
	var clusters = make([]Cluster, 0)
	// A fully disconnected set of tiles will result in as many clusters
	// One layer of tiles all connected and on the same layer results in one cluster

	for _, tile := range unClustered {
		clusters = append(clusters, Cluster{Members: []Tile{tile}})
	}

	return clusters
}

// SingleSolution takes a problem definition via the input parameter.
// Given the dimensions of a numerical matrix, and the numbers to fill it,
// Return the expected volume of water which would be trapped in a landscape
// 	described by the matrix. Refer to the file beginning for the details
func SingleSolution(input io.Reader) int {
	var problemDefinition *scanner.Scanner = &scanner.Scanner{}
	problemDefinition = problemDefinition.Init(input)
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
	for iRow := 0; iRow < matrixHeight; iRow++ {
		for iiCol := 0; iiCol < matrixWidth; iiCol++ {
			heightOfTile := get.NextInt() // SECTION 210 H
			tilesByHeight[heightOfTile] = append(tilesByHeight[heightOfTile], Tile{rowCoordinate: iRow, colCoordinate: iiCol})
			saveNumbers = append(saveNumbers, heightOfTile)
		}
	}

	// Build a matrix describing a landscape with the saved numbers from SECTION 210 H
	// TODO Decide whether it's necessary to build a copy matrix just so we can modify one matrix without concern
	var bm BasicMatrix = BasicMatrix(generateTwoDimArray(matrixHeight, matrixWidth))
	var matt Matrix = &bm
	for iRow := 0; iRow < matrixHeight; iRow++ {
		for iiCol := 0; iiCol < matrixWidth; iiCol++ {
			recallHeight := saveNumbers[iRow*matrixWidth+iiCol]
			matt.Set(iRow, iiCol, recallHeight)
		}
	}

	// Find max to decide when to stop searching later
	var maxHeight int
	var minHeight int = 0 // Per the problem definition, this could change
	for _, newHeightValue := range saveNumbers {
		if newHeightValue > maxHeight {
			maxHeight = newHeightValue
		}
	}
	var tm TerrainMap = TerrainMap{matt, tilesByHeight}
	// A layer includes connected tiles of the same current height
	// A cluster *could* have different height tiles by design, but is intended for height-differentiation.
	//	 A cluster *could* alternatively grow upward instead, but this is not preferred
	var recordClustersPerLayer LayerClusterMap // Store found clusters here by height
	// Use Breadth first search, BFS, to find a cluster whilst we record the min height of the perimeter.
	// Start the BFS from the lowest height, tiles at minHeight
	for iHeight := minHeight; iHeight < maxHeight; iHeight++ {
		recordClustersPerLayer[iHeight] = clusterTogether(tilesByHeight[iHeight], tm.terrain)
	}

	// TODO share the top-height value among cluster members, and update it on the fly.

	// TODO Find answer
	var volumeOfRainWater int
	return volumeOfRainWater
}

type Pond struct {
	perimeter []Tile
	interior  []Tile
}

// Expand a pond cluster searches for neighbor
// One candidate for this cluster is lower than necessary for pond water retention per rules.b				// One candidate for this cluster is lower than necessary for pond water retention per rules.ring 	// One candidate for this cluster is lower than necessary for pond water retention per rules.land and 	// One candidate for this cluster is lower than necessary for pond water retention per rules.includes	// One candidate for this cluster is lower than necessary for pond water retention per rules. them i // One candidatp foo thns clusfuncis(l wer *hPn nec)ssary for pond waExp retentann ped rules.
func (p *Pond) Expand() (isChanged bool) {
	const (
		interior  = iota
		perimeter = iota
		exterior  = iota
	)
	var validation map[Tile]int
	var needValidation []Tile
	var newPerimeterCount int

	for _, each := range p.interior {
		// Do not add the interior to the exterior
		validation[each] = interior
	}
	for _, per := range p.perimeter {
		// When an adjacent tile is outside the pond, it is either in the perimeter or the exterior
		// The exterior becomes the new perimeter
		validation[per] = perimeter
		needValidation = append(needValidation, FindAdjacent(per)...)
	}
	for _, needy := range needValidation {
		// This step reduces redundant entries and helps filter
		if validation[needy] != interior && validation[needy] != perimeter {
			validation[needy] = exterior
			newPerimeterCount++
		}
	}
	// Now that we are shifting from old to new, include the old perimeter in the current interior
	p.interior = append(p.interior, p.perimeter...) // And then,
	p.perimeter = make([]Tile, 0, newPerimeterCount)

	for tile, location := range validation {
		// This step filters the explored tiles by their location
		if location == exterior {
			p.perimeter = append(p.perimeter, tile)
		}
	}

	return false
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
	var twoDim = make([][]int, firstDim)
	var slidingStart, slidingFinish int

	for iFirst := 0; iFirst < firstDim; iFirst++ {
		slidingStart = secondDim * iFirst
		slidingFinish = secondDim*iFirst + secondDim
		twoDim[iFirst] = oneDim[slidingStart:slidingFinish]
	}

	return twoDim
}
