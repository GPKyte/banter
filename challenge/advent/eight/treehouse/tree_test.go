package treehouse

import(
    "testing"

    "github.com/google/go-cmp/cmp"
)

func TestFormSurroundings(t *testing.T) {
    main := NewTree(5)
    easternTaller := NewTree(9)
    westernTaller := NewTree(6)
    southernTaller := NewTree(8)
    northernTaller := NewTree(8)

    main.Surroundings.North.Notice(northernTaller)
    main.Surroundings.East.Notice(easternTaller)
    main.Surroundings.South.Notice(southernTaller)
    main.Surroundings.West.Notice(westernTaller)

    if main.Status().isVisible() {
        t.Log("Expected this tree to be hidden on all sides")
    }
}

func TestPathOfTrees(t *testing.T) {
    treePath := make([]*Tree, 0)
    for _, h := range []int{1,2,3,4,5,7,6,5,4,3} {
        treePath = append(treePath, NewTree(h))
    }

    // Check East to West along treePath
    for i := 0; i+1 < len(treePath); i++ {
        eastNeighbor := treePath[i]
        ett := eastNeighbor.Surroundings.East.TallestTree
        current := treePath[i+1]
        current.Surroundings.East.Notice(ett)
        current.Surroundings.East.Notice(eastNeighbor)
    }
    // Get tallest Tree Result
    tallestTreeToTheEastHeightsPerTree := make([]int, len(treePath))
    for i, t := range treePath {
        tt := t.Surroundings.East.TallestTree
        tallestTreeToTheEastHeightsPerTree[i] = Height(tt)
    }
    if !cmp.Equal(tallestTreeToTheEastHeightsPerTree, []int{0,1,2,3,4,5,7,7,7,7}) {
        t.Fail()
        t.Logf(`Tallest Tree to the East per tree did not match expectation.
        Found %v given this tree path %v`, tallestTreeToTheEastHeightsPerTree, treePath)
    }

    // Check whether Exposed to the East
    exposureExpectation := []bool{true,true,true,true,true,true,false,false,false,false}
    for i, tree := range treePath {
        if exposureExpectation[i] != tree.Status().EasternExposure {
            t.Fail()
            t.Log("Exposure expectation different than reality.", i, tree)
        }
    }
}
