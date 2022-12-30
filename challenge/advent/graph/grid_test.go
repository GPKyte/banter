package graph

import(
    "testing"
    "strings"
)

var g5x5 = `
aaaaa
aaaaa
aaaaa
aaaaa
aaaaa`[1:]

func TestGridOverlayLooksRight(t *testing.T) {
    g := NewGrid(strings.NewReader(g5x5))
    coords := []Coordinate{
        {X: 0, Y: 0},
        {X: 0, Y: 1},
        {X: 0, Y: 2},
        {X: 1, Y: 2},
        {X: 2, Y: 2},
        {X: 2, Y: 3},
        {X: 1, Y: 3},
        {X: 1, Y: 4},
        {X: 2, Y: 4},
        {X: 3, Y: 4},
        {X: 4, Y: 4},
    }
    for _, c := range coords {
        g.Get(c).Overlay = "X"
    }
    expectation := `
Xaaaa
Xaaaa
XXXaa
aXXaa
aXXXX`[1:]
    if expectation != g.String() {
        t.Fail()
        t.Log("Overlay Invalid, see below\n"+g.String())
    }
}
