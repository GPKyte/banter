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
func TestLoadForest(t *testing.T) {
    
}

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

