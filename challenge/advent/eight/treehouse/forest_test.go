package treehouse

import (
    "testing"
    "strings"
)

var ExampleULE = `1918274629
8126124510
9127371021
2190317202
3728211914
0841272722
4453354444`
var ExampleU5x5E = `53353
22003
02410
30401
22552`

func TestForest_ShowHeights(t *testing.T) {
    f := NewForest(strings.NewReader(ExampleU5x5E))
    heights := f.showHeights()
    if heights != ExampleU5x5E {
        t.Fail()
        t.Log(heights)
        t.Log("vs.")
        t.Log(ExampleU5x5E)
    }
}

func TestForest_CountVisible(t *testing.T) {
    f := NewForest(strings.NewReader(ExampleU5x5E))

    if cv := f.CountVisible(); cv != 20 {
        t.Fail()
        t.Log("5x5 Forest visible count is", cv)
    }
}

func TestForest_showVisible(t *testing.T) {
    f := NewForest(strings.NewReader(ExampleU5x5E))

    expectation := `32122
10003
11311
20201
21322`
    vis := f.showVisibility()
    if vis != expectation {
        t.Fail()
        t.Log(vis[0:5],   expectation[0:5],   ExampleU5x5E[0:5])
        t.Log(vis[6:11],  expectation[6:11],  ExampleU5x5E[6:11])
        t.Log(vis[12:17], expectation[12:17], ExampleU5x5E[12:17])
        t.Log(vis[18:23], expectation[18:23], ExampleU5x5E[18:23])
        t.Log(vis[24:29], expectation[24:29], ExampleU5x5E[24:29])
    }
}

var ExampleRTX = `30373
25512
65332
33549
35390`

// - - - - - Scenic Score for ExampleRTX
// - 1 4 1 -
// - 6 1 2 -
// - 1 8 3 -
// - - - - -

func TestVisibleCountInAnotherExample(t *testing.T) {
    f := NewForest(strings.NewReader(ExampleRTX))
    vis := f.CountVisible()
    if vis != 21 {
        t.Fail()
        t.Log(vis)
    }
}

func TestScenicScore(t *testing.T) {
    // Each direction from a given tree,
    // There are zero or more trees shorter than the given tree
    // Until reaching the forest edge,
    // or a tree of equal height or taller
    // Including all shorter trees and the first same or greater height tree and multiplying the count of trees in each direction yields the scenic score of a tree

    f := NewForest(strings.NewReader(ExampleRTX))
    s := f.FindMostScenicTreeScore()
    if s != 8 {
        t.Fail()
        t.Log(s)
        t.Log(f.ScenicScores)
        n, e, s, w := f.MeasureScene(3,2)
        t.Log(n,e,s,w)
    }
}

func TestMeasureScene(t *testing.T) {
    f := NewForest(strings.NewReader(ExampleRTX))
    n,e,s,w := f.MeasureScene(3,0)
    if e != 0 {
        t.Fail()
        t.Log(n,e,s,w)
    }
    n,e,s,w = f.MeasureScene(1,1)
    if n*e*s*w != 1 {
        t.Fail()
        t.Log(n,e,s,w)
    }
    n,e,s,w = f.MeasureScene(2,1)
    if w*s != 6 {
        t.Fail()
        t.Log(w,s)
    }
}

