package main

import (
    "os"
    "testing"
)

func TestGridHistory(t *testing.T){}


func TestGridMovement(t *testing.T){
    g := NewGrid()
    g.MoveRope(Up, 1)
    g.MoveRope(Down, 1)
    g.MoveRope(Right, 1)
    g.MoveRope(Left, 1)

    g.MoveRope(Up, 2)
    g.MoveRope(Right, 3)

    head := g.Rope.Head
    tail := g.Rope.Tail

    if head.Vertical != tail.Vertical {t.Fail()}
    if head.Horizontal != 3 {t.Fail()}
}

func TestHistory(t *testing.T) {
    coordinates := []struct{x, y, uniqueCount int}{
        { 0,  0,  1}, // Start
        { 1,  0,  2}, // Right
        { 1,  1,  3}, // Up
        { 1,  2,  4}, // Up
        { 2,  2,  5}, // Right
        { 2,  1,  6}, // Down
        { 1,  1,  6}, // Left
        { 2,  0,  7}, // Right, Down
        { 2,  1,  7}, // Up
        { 3,  2,  8}, // Right, Up
    }

    g := NewGrid()
    for _, coo := range coordinates {
        g.Rope.Tail.MoveTo(coo.x, coo.y)
        g.RecordHistory()
        locations := g.SummarizeHistory()
        if len(locations) != coo.uniqueCount {
            t.Fail()
            t.Log(locations)
        }
    }
}

func TestExample(t *testing.T) {
    ex, _ := os.Open("testdata/example0")
    defer ex.Close()

    g := NewGrid()
    g.ApplyMovements(ex)
    got := g.SummarizeHistory()
    want := 13

    if len(got) != want {
        t.Fail()
        t.Log(g.History)
    }
}
