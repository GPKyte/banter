package graph

import(
    "testing"
    "github.com/google/go-cmp/cmp"
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

func TestPathTreeLinking(t *testing.T) {
    g := NewGrid(strings.NewReader(g5x5))

    pt := NewPathTree(g)
    hereToThere := []struct{here, there Coordinate}{
        {Coordinate{X: 0, Y: 0}, Coordinate{X: 0, Y: 1}},
        {Coordinate{X: 0, Y: 1}, Coordinate{X: 0, Y: 2}},
        {Coordinate{X: 0, Y: 2}, Coordinate{X: 1, Y: 2}},
        {Coordinate{X: 1, Y: 2}, Coordinate{X: 2, Y: 2}},
        {Coordinate{X: 2, Y: 2}, Coordinate{X: 2, Y: 3}},
        {Coordinate{X: 2, Y: 3}, Coordinate{X: 1, Y: 3}},
        {Coordinate{X: 1, Y: 3}, Coordinate{X: 1, Y: 4}},
        {Coordinate{X: 1, Y: 4}, Coordinate{X: 2, Y: 4}},
        {Coordinate{X: 2, Y: 4}, Coordinate{X: 3, Y: 4}},
        {Coordinate{X: 3, Y: 4}, Coordinate{X: 4, Y: 4}},
    }
    for each, pair := range hereToThere {
        t.Logf("%02d: Linking pair %v", each, pair)
        pt.Link(pair.here, pair.there)
    }
    wantTheseCoordsInTrace := []Coordinate{
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

    lastCoordPair := hereToThere[len(hereToThere)-1]
    gotTheseCoords := pt.TracePath(lastCoordPair.there)

    if !cmp.Equal(wantTheseCoordsInTrace, gotTheseCoords) {
        t.Fail()
        t.Log("Undesireable outcome of path tree trace. Comparison below:\n"+cmp.Diff(wantTheseCoordsInTrace, gotTheseCoords))
    }
}

