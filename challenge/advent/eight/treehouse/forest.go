package treehouse

import (
    "io"
    "os"
    "log"
    "bufio"
    "fmt"
    "strings"
    "strconv"
)

var Debug *log.Logger = log.New(os.Stderr, "[DEBUG]: ", 0)

func NewForest(fromInput io.Reader) *Forest {
    f := Forest{
        Trees: make([][]*Tree, 0),
        ScenicScores: make([][]int, 0),
    }
    transformTextIntoTreesWithinForest := func(text string) {
        moreTrees := make([]*Tree, 0)

        for i := 0; i < len(text); i++ {
            height, err := strconv.Atoi(text[i:i+1])
            if err != nil {
                Debug.Fatalln(err, text, "at index", i)
            }
            moreTrees = append(moreTrees, NewTree(height))
        }
        f.Trees = append(f.Trees, moreTrees)
    }
    ScanEachLine(fromInput, transformTextIntoTreesWithinForest)

    // Init status of every tree
    f.InspectFromTheOutside()

    return &f
}

func ScanEachLine(ofThis io.Reader, andDoThat func(this string)) {
    i := bufio.NewScanner(ofThis)
    i.Split(bufio.ScanLines)

    for allData := i.Scan(); allData; allData = i.Scan() {
        withThis := i.Text()
        andDoThat(withThis)
    }
}

type Forest struct {
    Trees [][]*Tree
    ScenicScores [][]int
}

func (f *Forest) showVisibility() string {
    exposedSidesPerTree := make([]string, 0)
    countExposures := func(t *Tree) {
        var count int
        ts := t.Status()

        for _, exp := range []bool{
            ts.NorthernExposure,
            ts.EasternExposure,
            ts.SouthernExposure,
            ts.WesternExposure,
        } {
            if exp {
                count++
            }
        }
        exposedSidesPerTree = append(exposedSidesPerTree, fmt.Sprint(count))
    }
    return f.show(&exposedSidesPerTree, countExposures)
}
func (f *Forest) showHeights() string {
    heights := make([]string, 0)
    viaTreeString := func(t *Tree) {
        heights = append(heights, t.String())
    }
    return f.show(&heights, viaTreeString)
}
func (f *Forest) show(this *[]string, via func(t *Tree)) string {
    f.traverse(via)

    rows := make([]string, 0)
    
    for northToSouth := 0; northToSouth < len(f.Trees); northToSouth++ {
        eastToWest := len(f.Trees[northToSouth])
        start := northToSouth*eastToWest
        end := start + eastToWest
        r := strings.Join((*this)[start:end], "")
        rows = append(rows, r)
    }

    return strings.Join(rows, "\n")
}

func (f *Forest) String() string {return f.showVisibility()}

// Any tree taller than the path of its neighbors from the left, right, top, or bottom until the respective border is considered visible
// Any tree which has a taller tree between itself and the forest border in every direction is considered hidden


// The goal is to find the number of visible trees, i.e. those trees with at least one path of trees smaller than itself.


func (f *Forest) traverse(andDoThisPerTree func(t *Tree)) {
    for northToSouth := range f.Trees {
        for eastToWest := range f.Trees[northToSouth] {
            t := f.Trees[northToSouth][eastToWest]
            andDoThisPerTree(t)
        }
    }
}

func (f *Forest) CountVisible() int {
    var visibles int = 0

    countIfIsVisible := func(t *Tree) {
        if t.Status().isVisible() {
            visibles += 1
        }
    }

    f.traverse(countIfIsVisible)
    return visibles
}

func (f *Forest) InspectFromTheOutside() {
    // For every direction (e.g. Westbound)
    // Let each tree know the tallest tree visible in the opposite direction (e.g. Eastbound)
    // Each tree can share this information with the adjacent tree to save time
    // To remove unwanted condition checking, use a tree of height 0 to start each loop rather than a nil pointer

    // Westbound
    for northToSouth := range f.Trees {
        var neighbor *Tree = NewTree(-1)
        for westbound := range f.Trees[northToSouth] {
            tallestFromEast := neighbor.Surroundings.East.TallestTree
            current := f.Trees[northToSouth][westbound]
            current.Surroundings.East.Notice(neighbor)
            current.Surroundings.East.Notice(tallestFromEast)
            neighbor = current
        }
    }
            
    // Eastbound
    for northToSouth := range f.Trees {
        var neighbor *Tree = NewTree(-1)
        for westbound := range f.Trees[northToSouth] {
            eastbound := len(f.Trees[northToSouth]) - (1 + westbound)
            tallestFromWest := neighbor.Surroundings.West.TallestTree
            current := f.Trees[northToSouth][eastbound]
            current.Surroundings.West.Notice(neighbor)
            current.Surroundings.West.Notice(tallestFromWest)
            neighbor = current
        }
    }

    // Northbound
    for eastToWest := range f.Trees[0] {
        var neighbor *Tree = NewTree(-1)
        for southbound := range f.Trees {
            northbound := len(f.Trees) - (1 + southbound)
            tallestFromSouth := neighbor.Surroundings.South.TallestTree
            current := f.Trees[northbound][eastToWest]
            current.Surroundings.South.Notice(neighbor)
            current.Surroundings.South.Notice(tallestFromSouth)
            neighbor = current
        }
    }

    // Southbound
    for eastToWest := range f.Trees[0] {
        var neighbor *Tree = NewTree(-1)
        for southbound := range f.Trees {
            tallestFromNorth := neighbor.Surroundings.North.TallestTree
            current := f.Trees[southbound][eastToWest]
            current.Surroundings.North.Notice(neighbor)
            current.Surroundings.North.Notice(tallestFromNorth)
            neighbor = current
        }
    }

    //
}
func (f *Forest) MeasureScene(x, y int) (n, e, s, w int) {
    // Include all shorter trees up till the edge,
    // Include the first tree that is >= height, and none beyond it
    centerTree := f.Trees[x][y]

    northEdge := 0
    eastEdge := 0
    southEdge := len(f.Trees) - 1
    westEdge := len(f.Trees[northEdge]) - 1
    inScene := func(ns, ew int) bool {
        outsideForestEdges := ns < northEdge || ns > southEdge || ew < eastEdge || ew > westEdge
        return !outsideForestEdges
    }
    viewBlocked := func(ns, ew int) bool {
        return f.Trees[ns][ew].Height >= centerTree.Height
    }

    for northbound := x-1; inScene(northbound, y); northbound-- {
        n += 1

        if viewBlocked(northbound, y) {
            break
        }
    }
    for eastbound := y-1; inScene(x, eastbound); eastbound-- {
        e += 1

        if viewBlocked(x, eastbound) {
            break
        }
    }
    for southbound := x+1; inScene(southbound, y); southbound++ {
        s += 1

        if viewBlocked(southbound, y) {
            break
        }
    }
    for westbound := y+1; inScene(x, westbound); westbound++ {
        w += 1

        if viewBlocked(x, westbound) {
            break
        }
    }

    return n, e, s, w
}

func (f *Forest) CalculateScenicScores() {
    for northToSouth := range f.Trees {
        f.ScenicScores = append(f.ScenicScores, make([]int, len(f.Trees[northToSouth])))

        for eastToWest := range f.Trees[northToSouth] {
            f.ScenicScores[northToSouth][eastToWest] = product(f.MeasureScene(northToSouth, eastToWest))
        }
    }
}
func (f *Forest) FindMostScenicTreeScore() int {
    var most int
    f.CalculateScenicScores()

    for northToSouth := range f.Trees {
        for eastToWest := range f.Trees[northToSouth] {
            ss := f.ScenicScores[northToSouth][eastToWest]
            if most < ss {
                most = ss
            }
        }
    }

    return most
}

func product(of ...int) int {
    prod := 1
    for _, all := range of {
        prod *= all
    }
    return prod
}
