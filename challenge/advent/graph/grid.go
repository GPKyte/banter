package graph

import (
    "io"
    "bufio"
    "strings"
)

type Coordinate struct {
    X, Y int
}

type Grid struct {
    self []*[]*GridNode
}

type GridNode struct {
    Label   string
    Overlay string
    Height  int
    Visited bool
}

// Path is an acyclic directed subgraph built on Grid
// Start <- pn <- pn <- pn <- End
type Path struct {
    End *PathNode
}

// PathNode can be traversed back to the beginning
type PathNode struct {
    Stem    *PathNode
    At      Coordinates   
}

type Coordinates struct {
    X int
    Y int
}

func NewGrid(from io.Reader) *Grid {
    ls := GetLines(from)
    nodes := make([]*[]*GridNode, len(ls)) // Remember to assign rather than append

    for i, l := range ls {
        ns := LineToGridNodes(l)
        nodes[i] = ns
    }

    return &Grid{self: nodes}
}

func LineToGridNodes(s string) *[]*GridNode {
    ns := make([]*GridNode, len(s)) // assign rather than append
    for i, si := range s {
        ns[i] = &GridNode{
            Label:  string(si),
            Height: 0,
            Overlay:"",
            Visited:false,
        }
    }
    return &ns
}

// GetLines returns slice of all nonempty strings from input
func GetLines(from io.Reader) []string {
    s := bufio.NewScanner(from)
    lines := make([]string, 0)
    for ok := s.Scan(); ok; ok = s.Scan() { 
        l := s.Text()
        if len(l) == 0 {
            continue
        }
        lines = append(lines, l)
    }
    return lines
}

func (g *Grid) Get(c Coordinate) *GridNode {
    return (*g.self[c.Y])[c.X]
}

func (g *Grid) String() string {
    buffer := make([]string, 0, (g.XBound()+1)*+g.YBound()) // +1 for new line

    for _, line := range g.self {
        for _, gn := range *line {
            buffer = append(buffer, gn.String())
        }
        buffer = append(buffer, "\n")
    }

    return strings.Join(buffer[:len(buffer)-1], "") // Ditch last new line
}

func (gn *GridNode) String() string {
    s := gn.Label

    if gn.OverlayActive() {
        s = gn.Overlay
    }

    return s
}

func (gn *GridNode) OverlayActive() bool {
    return gn.Overlay != ""
}

func (g *Grid) XBound() int {
    return len(*g.self[0])
}

func (g *Grid) YBound() int {
    return len(g.self)
}

func (gn *GridNode) Visit() {
    gn.Visited = true
}

type PathTree struct {
    Links [][][]Coordinate
}

func NewPathTree(g *Grid) *PathTree {
    return nil
}

func (pt *PathTree) Link(a, b Coordinate) {}
func (pt *PathTree) At(a Coordinate) *PathNode {return nil}

func Overlay(g *Grid, pt *PathNode) {

}
