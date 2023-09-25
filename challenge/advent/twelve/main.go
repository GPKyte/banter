package main

import (
    "io"

    "github.com/GPKyte/banter/challenge/advent/graph"
)

func main() {}

func SolvePartOne(from io.Reader) int {
    G := graph.NewGrid(from)
    Q := NewCoordinateQueue() // for Breadth First Search
    pathTree := graph.NewPathTree(G)

    neighbors := func(c graph.Coordinate) []graph.Coordinate {
        // Not Visited
        // Within bounds
        // Height at most one greater
        ns := make([]graph.Coordinate, 0, 4)
        for _, delta := range []struct{X, Y int}{{0,1},{1,0},{0,-1},{-1,0}} {
            n := graph.Coordinate{X: c.X + delta.X, Y: c.Y + delta.Y}

            inBounds    := n.X < G.XBound() && n.Y < G.YBound()
            notVisited  := !G.Get(n).Visited
            reachable   := LabelToHeight(G.Get(n).Label) <= (1 + LabelToHeight(G.Get(c).Label))

            if inBounds && notVisited && reachable {
                ns = append(ns, n)
            }
        }
        return ns
    }

    var start graph.Coordinate // Find Starting Location
    var goal string = "E" // Label to look for

    Q.NQ(start)
    var ok bool
    var nomad graph.Coordinate
    for nomad, ok = Q.DQ(); ok; nomad, ok = Q.DQ() {
        G.Get(nomad).Visit()

        if G.Get(nomad).Label == goal {
            break
        }

        for _, n := range neighbors(nomad) {
            Q.NQ(n)
            pathTree.Link(nomad, n)
        }
    }

    graph.Overlay(G, pathTree.At(nomad))

    lengthOfPath := 0
    return lengthOfPath
}

func LabelToHeight(s string) int {
    return int(byte(s[0]))
}

type CoordinateQueue []graph.Coordinate
func NewCoordinateQueue() *CoordinateQueue {
    q := make(CoordinateQueue, 0)
    return &q
}
func (cq *CoordinateQueue) Len() int {
    return len(*cq)
}
func (cq *CoordinateQueue) NQ(c graph.Coordinate) {}
func (cq *CoordinateQueue) DQ() (c graph.Coordinate, ok bool) {
    if len(*cq) == 0 {
        return graph.Coordinate{}, false
    }

    c = (*cq)[0]
    *cq = (*cq)[1:]
    return c, true
}
