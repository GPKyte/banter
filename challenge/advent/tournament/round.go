package tournament

import (
    "io"
    "text/scanner"
)

type Round struct {
    ag, zg Gesture
}

func (r Round) String() string {
    return r.ag.String()+" vs. "+r.zg.String() 
}

type Rounds []Round

const (
    PointsForWin = 6
    PointsForTie = 3
    PointsForLoss = 0
)

func (r *Round) Outcome() RoundOutcome {
    ap := r.ag.points
    zp := r.zg.points

    if r.ag.Beats(r.zg) {
        ap += PointsForWin
        zp += PointsForLoss
    } else if r.zg.Beats(r.ag) {
        zp += PointsForWin
        ap += PointsForLoss
    } else {
        ap += PointsForTie
        zp += PointsForTie
    }

    return RoundOutcome{a: ap, z: zp}
}

type RoundOutcome struct {
    a, z int
}

type Gesture struct {
    points int
    id int
    name string
}

func (g *Gesture) String() string {
    return g.name
}

var (
    NotA = Gesture{points: 0, id: 128, name: "Not A Gesture"}
    Rock = Gesture{points: 1, id: 64, name: "Rock"}
    Paper = Gesture{points: 2, id: 32, name: "Paper"}
    Scissors = Gesture{points: 3, id: 16, name: "Scissors"}
    GestureParseMap = map[string]Gesture{
        "A": Rock,
        "X": Rock,
        "B": Paper,
        "Y": Paper,
        "C": Scissors,
        "Z": Scissors,
    }
)

func (g Gesture) Beats(gg Gesture) bool {
    rockAndScissors := g == Rock && gg == Scissors
    paperAndRock := g == Paper && gg == Rock
    scissorsAndPaper := g == Scissors && gg == Paper

    if rockAndScissors || paperAndRock || scissorsAndPaper {
        return true
    }
    return false
}

// parseGesture from shortform A B C or X Y Z input
func parseGesture(from string) Gesture {
    g, ok := GestureParseMap[from]
    if !ok {
        g = NotA
    }
    return g
}

func loadRounds(src io.Reader) *Rounds {
    rnds := make(Rounds, 0)

    var s scanner.Scanner
    s.Init(src)

    var empty, abc, xyz string // abc will be read a loop prior to xyz by scanner    
    for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
        if abc == empty {
            abc = s.TokenText()
            continue
        }
        xyz = s.TokenText()

        ag := parseGesture(abc)
        zg := parseGesture(xyz)
        rnds = append(rnds, Round{ag, zg})

        xyz, abc = empty, empty
    }

    return &rnds
}
