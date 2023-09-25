package main

import (
    "io"
    "bufio"
    "strconv"
    "os"
    "log"
)

var Debug = log.New(os.Stderr, "[DEBUG]:", 0)

type Direction int

const (
    Up Direction = iota
    Down
    Left
    Right
)

type Grid struct {
    History map[RopeEnd]bool    // Using map rather than [][]int to avoid resizing slice during rope movements
    Rope    Rope                // Represents the location of the rope
    MaxVert int                 // Caching this data for string
    MaxHori int                 // Caching this data for string
}

// Rope describes the placement of the head and tail of a rope
// Making it easier to move as a unit and track elsewhere
type Rope struct {
    Head    RopeEnd
    Tail    RopeEnd
}

// RopeEnd loses some information (head/tail detail) but helps store location
type RopeEnd struct {
    Vertical   int
    Horizontal int
}

func NewGrid() *Grid {return &Grid{History: make(map[RopeEnd]bool)}}
func (g *Grid) String() string {return ""}

var TextToDirection = map[string]Direction{
    "U": Up,
    "D": Down,
    "R": Right,
    "L": Left,
}
func (g *Grid) ApplyMovements(from io.Reader) {
    the := bufio.NewScanner(from)
    the.Split(bufio.ScanLines)

    for allLines := the.Scan(); allLines; allLines = the.Scan() {
        line := the.Text() // D 42
        dir := TextToDirection[line[0:1]]
        dist, err := strconv.Atoi(line[2:])
        if err != nil {
            Debug.Printf("Error while parsing %s, please review:\n%v", line, err)
        }
        g.MoveRope(dir, dist)
    }
}
func (g *Grid) MoveRope(dir Direction, dist int) {
    cover := 0
    entireDistance := func() bool {
        cover++
        return cover <= dist
    }

    for entireDistance() {
        if Overlapping(g.Rope.Head, g.Rope.Tail) {
            g.Rope.Head.Move(dir)
        } else if Adjacent(g.Rope.Head, g.Rope.Tail) {
            body := g.Rope.Head
            g.Rope.Head.Move(dir)

            // Ensure tail is adjacent to head with one of eight moves
            if !Adjacent(g.Rope.Head, g.Rope.Tail) {
                g.Rope.Tail.MoveTo(body.Horizontal, body.Vertical)
            }
        }

        g.RecordHistory()
    }
}

func (g *Grid) SummarizeHistory() []RopeEnd {
    res := make([]RopeEnd, 0)
    for re := range g.History {
        res = append(res, re)
    }
    return res
}
func (g *Grid) RecordHistory() {
    tailCopy := RopeEnd(g.Rope.Tail)
    g.History[tailCopy] = true
    g.AdjustBorderEdge(g.Rope.Tail.Horizontal, g.Rope.Tail.Vertical)
    g.AdjustBorderEdge(g.Rope.Head.Horizontal, g.Rope.Head.Vertical)
}
func (g *Grid) AdjustBorderEdge(hori, vert int) {
    type update struct{current *int; vsNew int}
    var apply = func(u update) {
        if *u.current < u.vsNew {
            *u.current = u.vsNew
        }
    }
    var updates = []update{
        {&g.MaxHori, hori},
        {&g.MaxVert, vert},
    }
    
    for _, potentialUpdate := range updates {
        apply(potentialUpdate)
    }
}

// There are two cases for a rope
// The head overlaps with the tail,
// The tail is in one of eight locations, adjacent to the head
// In the case of overlap, the head can move to one of the eight spots
// In the case of adjacency, the head can move to overlap,
// Or the head can move to another adjacent position
// Or the head can move to a non-adjacent position
// Of all the adjacent positions,
// A move to any of the non-adjacent positions
// requires the tail to assume the heads position prior to the last move.
// . . . . . . . . 
// . . h h h . . . 
// . h H H H h . . 
// . h H T H h . . 
// . h H H H h . . 
// . . h h h . . . 
// . . . . . . . . 
func Overlapping(head, tail RopeEnd) bool {
    return head == tail
}
func Adjacent(head, tail RopeEnd) bool {
    vdist := head.Vertical - tail.Vertical
    hdist := head.Horizontal - tail.Horizontal

    return !(hdist > 1 || hdist < -1 || vdist > 1 || vdist < -1)
}

var OneMove = map[Direction]struct{hori, vert int}{
    Up:     { 0,  1},
    Down:   { 0, -1},
    Left:   {-1,  0},
    Right:  { 1,  0},
}

func (re *RopeEnd) Move(dir Direction) {
    re.Horizontal += OneMove[dir].hori
    re.Vertical += OneMove[dir].vert
}

func (re *RopeEnd) MoveTo(h, v int) {
    re.Horizontal = h
    re.Vertical = v
}
