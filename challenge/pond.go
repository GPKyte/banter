package challenge

import (
	"fmt"
	"io"
	"strconv"
	"strings"
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

// Q ueue of Tiles found via the BFS method on the prestructured data (e.g. Matrix, graph)
type Q struct {
	fifo *[]Tile
}

// Serve will bring the next Tile out from a wait state
func (q *Q) serve() Tile {
	// Choose to allow OOB to panic
	var beingServed Tile = (*q.fifo)[0]
	*q.fifo = (*q.fifo)[1:]

	return beingServed
}

// Wait will enqueue
func (q *Q) wait(here Tile) {
	// Will this lack of pointing to references lead to a frontline issue?
	*q.fifo = append(*q.fifo, here)
}

func (q *Q) empty() bool {
	return len(*q.fifo) == 0
}

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

type LayerClusterMap map[int][]Cluster

// Visitor is an interface to wrap the Visit function needed now for matrix traversal
type Visitor interface {
	Visit(interface{}) bool
}

// BasicVisitor maintains a history of Tiles visited at every layer to assist navigation of tile neighbors
type BasicVisitor struct {
	history map[interface{}]bool
}

// Visit a Tile to mark that area and avoid duplicate work later, this is a version of shortcircuiting in practice.
func (v *BasicVisitor) Visit(time interface{}) (beenHereBefore bool) {
	// We know that we have not
	beenHereBefore = false
	// except...
	if v.history[time] {
		beenHereBefore = true
	}

	v.history[time] = true

	// It may be useful to know whether we have
	return beenHereBefore
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

type NumberScanner interface {
	NextInt() int
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

// Total is the sum of the heights recorded in the matrix
func (m *BasicMatrix) Total() int {
	var sum int
	for i := 0; i < len(*m); i++ {
		for ii := 0; ii < len((*m)[i]); ii++ {
			sum += m.Get(i, ii)
		}
	}
	return sum
}

// Equals references deep values rather than identity
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

// ClusterTogether connectedTiles from the given Tiles based on adjacency
// Clusters are Useful for addressing several connected Tiles
func clusterTogether(unClustered []Tile, src Matrix) []Cluster {
	var clusters = make([]Cluster, 0)
	// A fully disconnected set of tiles will result in as many clusters
	// One layer of tiles all connected and on the same layer results in one cluster
	var friend Visitor = &BasicVisitor{history: map[interface{}]bool{}}
	var clusterLeaks bool = false
	var shortestTallerNeighbor int
	var middle = func(a, b, c int) int {
		for combo := range Permute([]int{a, b, c}) {
			if combo[0] >= combo[1] && combo[1] > combo[2] {
				return combo[1]
			}
		}
		return a
	}
	var height = func(t Tile) int {
		return src.Get(t.rowCoordinate, t.colCoordinate)
	}
	var queueStorage = make([]Tile, 0, len(unClustered))
	var q Q = Q{fifo: &queueStorage}
	var visitFriendsOf = func(t Tile) {
		for _, each := range FindAdjacent(t) {
			if height(each) < height(t) {
				clusterLeaks = true
				continue
			} else if height(each) > height(t) {
				shortestTallerNeighbor = middle(height(t), shortestTallerNeighbor, height(each))
			} else if friend.Visit(each) {
				continue
			} else if height(each) == height(t) {
				q.wait(each)
			}
		}
	}

	for _, uc := range unClustered {
		if friend.Visit(uc) {
			continue
		}
		var buddies = make([]Tile, 0, len(unClustered))
		buddies = append(buddies, uc)
		friend.Visit(uc)
		visitFriendsOf(uc)

		for !q.empty() {
			var mine = q.serve()
			buddies = append(buddies, mine)
			visitFriendsOf(mine)
		}

		var newCluster = Cluster{Members: buddies, anyMemberLeaks: clusterLeaks, height: height(uc)}
		clusterLeaks = false
		clusters = append(clusters, newCluster)
	}
	return clusters
}

// SingleSolution takes a problem definition via the input parameter.
// Given the dimensions of a numerical matrix, and the numbers to fill it,
// Return the expected volume of water which would be trapped in a landscape
// 	 described by the matrix. Refer to the file beginning for the details
func SingleSolution(input io.Reader) int {
	tm := SetupPondProblem(input)
	return SolvePondProblem(tm)
}

func SetupPondProblem(input io.Reader) TerrainMap {
	var problemDefinition *scanner.Scanner = &scanner.Scanner{}
	problemDefinition = problemDefinition.Init(input)
	var get = StdNumberScanner{From: problemDefinition}

	var matrixHeight, matrixWidth int
	/* Open this comment tag to not take height and width from scanner */
	matrixHeight = get.NextInt()
	matrixWidth = get.NextInt()

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

	return TerrainMap{matt, tilesByHeight}
}

func SolvePondProblem(tm TerrainMap) (volumeOfRainWater int) {
	var clustersForLayer LayerClusterMap = map[int][]Cluster{} // Layer each group of clusters by their height
	/* Important to start from the bottom of terrain to simulate rain water filling land vessels like ponds and valleys.
	   It is important because the next iteration will build on top of the previous layer */
	for layz := range tm.layers.Keys() {

		clustersForLayer[layz] = clusterTogether(tm.layers[layz], tm.terrain)

		for _, each := range clustersForLayer[layz] {
			if each.anyMemberLeaks {
				continue
			}

			var heightDifference int = 1
			var fillHeight = each.height + heightDifference
			for _, tile := range each.Members {
				// TODO share the top-height value among cluster members, and update it on the fly.
				tm.terrain.Set(tile.rowCoordinate, tile.colCoordinate, fillHeight)
			}
			volumeOfRainWater += heightDifference * len(each.Members)
		}
	}

	return volumeOfRainWater
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

func (m *BasicMatrix) String() string {
	var b strings.Builder
	for _, row := range *m {
		fmt.Fprintln(&b, row)
	}
	return b.String()
}

func (m *BasicMatrix) Visualize() {
	fmt.Print(m.String())
}
