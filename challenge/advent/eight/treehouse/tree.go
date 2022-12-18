package treehouse

func NewTree(height int) *Tree {
    return &Tree{
        Height: height,
        Surroundings: NewLinesOfSight(),
    }
}


// Tree is a data type
// Created from an integer during a larger loading process
type Tree struct {
    Height int
    Surroundings LinesOfSight
}

func Height(t *Tree) int {
    var height int

    if t != nil {
        height = t.Height
    } else {
        height = 0
    }

    return height
}

func (t *Tree) Status() TreeStatus {
    ts := t.Surroundings
    nx := Height(ts.North.TallestTree) < t.Height
    ex := Height(ts.East.TallestTree) < t.Height
    sx := Height(ts.South.TallestTree) < t.Height
    wx := Height(ts.West.TallestTree) < t.Height

    return TreeStatus {
        NorthernExposure: nx,
        EasternExposure:  ex,
        SouthernExposure: sx,
        WesternExposure:  wx,
    }
}

type TreeStatus struct {
    NorthernExposure bool
    EasternExposure  bool
    SouthernExposure bool
    WesternExposure  bool
}

func (ts TreeStatus) isVisible() bool {
    TrueWhenAnyDirectionExposed := ts.NorthernExposure || ts.EasternExposure || ts.SouthernExposure || ts.WesternExposure
    return TrueWhenAnyDirectionExposed
}

type LineOfSight struct {
    TallestTree *Tree
}

func NewLinesOfSight() LinesOfSight {
    return LinesOfSight{
        North: LineOfSight{},
        East: LineOfSight{},
        South: LineOfSight{},
        West: LineOfSight{},
    }
}

type LinesOfSight struct {
    North LineOfSight
    East  LineOfSight
    South LineOfSight
    West  LineOfSight
}

// Notice whether the given tree is taller than the known tallest tree
// And if it is taller, it becomes the new tallest tree.
func (los *LineOfSight) Notice(t *Tree) {
    if tt := los.TallestTree; tt == nil || t != nil && t.Height > tt.Height {
        los.TallestTree = t
    }
}
