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

// What would finding the cluster of an area look like?
// Say I look at Tile_0 from the input reel and it has two neighbors, Tile_01, Tile_02.
// I found the neighbors via a helper function "FindAdjacent".
// Tile_0 is a border tile in the upperleft corner of a Matrix or Terrain Map in this example.
// Tile_0 is first into the Q, followed by 01 and 02, 01's neighbors 01a 01b 01c then 02's neighbors 02a 01b 02c
// Exploring and revisiting these Tiles in a certain order via a Q affords well-known benefits from the BFS

// Q.ueue of Tiles found via the BFS method on the prestructured data (e.g. Matrix, graph)
type Q struct {
	frontLine *[]Tile
}
// Serve will bring the next Tile out from a wait state
func (q *Q) serve() Tile {
	const outOfRange = -99999
	const errorTile = Tile{outOfRange, outOfRange}
	var lenQ = len(*(*q).frontline) // This number appears several times locally

	// cannot return a Tile, error condition
	if lenQ<=0 {
		return errorTile
	}

	beingServed  = *(*q.frontline)[0]
	*q.frontline = *(*q.frontline)[1:]

	return beginServed
}
// Wait will enqueue
func (q *Q) wait(here Tile) {
	// Will this lack of pointing to references lead to a frontline issue?
	q.frontLine = append(q.frontline, here)
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

func clusterTogether(maybeConnected []Tile, indis Matrix) ([]Cluster, error) {
	// Prefer error to quiet? Can I use validator on Matrix instead of this check? Perhaps TODO
	if len(maybeConnected) == 0 || indis == nil {
		return nil, fmt.Errorf("Tiles are a required argument. Try again.")
	}
	// Given the tiles above,
	// Return the same tiles, but group any series of adjacent tiles together as a cluster.
	// It is convenient to decide here whether the cluster is leaky or not. Granted, we are lacking needed information to do so.
	var clusterHeight = maybeConnected[0]
	var lookupAid map[Tile][]Tile
	for _, seems := range maybeConnected {
		lookupAid[seems] = make([]Tile, 0, len(maybeConnected)-1) // Prep work, these will be lost if redundancy exists. Might exist?
	}

	// Background daemon will determine if clusterIsLeaky or not
	var maybeNeighbors chan Tile
	var clusterIsLeaky bool = false
	go func() {
		var now = true
		for now {
			var maybeLower = <-maybeNeighbors
			var heightMaybe = indis.Get(maybeLower.rowCoordinate, maybeLower.colCoordinate)
			// Is this cluster leaky? Find out while we wait
			if heightMaybe < clusterHeight {
				clusterIsLeaky = true
				break
			}
		}
		// When does this close??
	}()


	// Look up the neighbors and build the adjacency map to determine connectivity relavent to cluster analysis
	for _, maybe := range maybeConnected {
		var may []Tile = lookupAid[maybe] // Just an empty list to fill with neighbors that were found in the first cover loop
		var be []Tile = FindAdjacent(maybe) // The neighbors of this input tile which need matched against the other input tiles

		for _, bNeighbor := range be {
			var matches bool
			// Either the Tiles will be in both the input and the neighbor results, or just the latter
			if len(lookupAid[bNeighbor]) != 0 {

			}
			// Recall we return Clusters, i.e. []Tiles + metadata
			lookupAid[bNeighbor] = append(lookupAid[bNeighbor], maybe) // Add tile to neighbor
			lookupAid[maybe] = append(lookupAid[maybe], bNeighbor) // Add neighbor to tile

			newDistributedListOfMembers = append(lookupAid[bNeighbor], lookupAid[maybe]...)
		}

		// share reference to the conjoined slices that store neighbors, still use append


		//add matching neighbors to growing list, then append that to the matche's neighbor list?
		// what holds the list of neighbors? Edges rather are kept looked up by each vertex's list of neighbors, bidirection is maintained.
		// lookupAid[vertex] = []Tile
		lookupAid[maybe] 
	)
	return nil, nil
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

// SingleSolution takes a problem definition via the input parameter.
// Given the dimensions of a numerical matrix, and the numbers to fill it,
// Return the expected volume of water which would be trapped in a landscape
// 	described by the matrix. Refer to the file beginning for the details
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
	for iRow := 0; iRow < matrixHeight; iRow++ {
		for iiCol := 0; iiCol < matrixWidth; iiCol++ {
			heightOfTile := get.NextInt() // SECTION 210 H
			tilesByHeight[heightOfTile] = append(tilesByHeight[heightOfTile], Tile{rowCoordinate: iRow, colCoordinate: iiCol})
			saveNumbers = append(saveNumbers, heightOfTile)
		}
	}

	// Build a matrix describing a landscape with the saved numbers from SECTION 210 H
	// TODO Decide whether it's necessary to build a copy matrix just so we can modify one matrix without concern
	var matt = BasicMatrix(generateTwoDimArray(matrixHeight, matrixWidth))
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

	// A layer includes connected tiles of the same current height
	// A cluster *could* have different height tiles by design, but is intended for height-differentiation.
	//	 A cluster *could* alternatively grow upward instead, but this is not preferred
	var recordClustersPerLayer LayerClusterMap // Store found clusters here by height
	// Use Breadth first search, BFS, to find a cluster whilst we record the min height of the perimeter.
	// Start the BFS from the lowest height, tiles at minHeight

	for iHeight := minHeight; iHeight < maxHeight; iHeight++ {
		// Every tile at this height is either at it's original height or was promoted to this height
		// Every tile at this height is either connected to at least one tile of lesser height or none

		// Knowing this layer's clusters is helpful because the perimeter of a cluster decides
		//	whether or not the entire cluster of tiles stays the same height or is raised by the "rain water"
		recordClustersPerLayer[iHeight] = clusterTogether(tilesByHeight[iHeight])
		for _, cluster := range recordClustersPerLayer[iHeight] {
			if cluster.retainsRainWater() {
				// Promote tiles by raising them to next layer
				tilesByHeight[iHeight+1] = append(tilesByHeight[iHeight+1], cluster.Members...)
			}
		}
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

// Expand a pond cluster searches for neig	// One candidate for this cluster is lower than necessary for pond water retention per rules.b				// One candidate for this cluster is lower than necessary for pond water retention per rules.ring 	// One candidate for this cluster is lower than necessary for pond water retention per rules.land and 	// One candidate for this cluster is lower than necessary for pond water retention per rules.includes	// One candidate for this cluster is lower than necessary for pond water retention per rules. them i // One candidatp foo thns clusfuncis(l wer *hPn nec)ssary for pond waExp retentann ped rules.
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
	var twoDim = make([][]int, 0, firstDim)
	var slidingStart, slidingFinish int

	for iFirst := 0; iFirst < firstDim; iFirst++ {
		slidingStart = secondDim * iFirst
		slidingFinish = secondDim*iFirst + secondDim
		twoDim[iFirst] = oneDim[slidingStart:slidingFinish]
	}

	return twoDim
}
