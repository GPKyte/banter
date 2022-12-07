package tournament

import (
    "testing"
    "fmt"
    "strings"
)

var allRoundVariations string = `A X
A Y
A Z
B Y
B Z
B X
C Z
C X
C Y`
var gestures = []Gesture{Rock, Paper, Scissors}

func TestRoundParsing(t *testing.T) {
    all := loadRounds(strings.NewReader(allRoundVariations))
    t.Log(*all)
    if len(*all) != 9 {
        t.Fail()
    }
    allAsString := fmt.Sprint(*all)
    if len(allAsString) < len("RockPaperScissors") * 3 * 2 {
        t.Fail()
    }
    if strings.Contains(allAsString, "Z") {
        t.Fail()
    }
}

func TestGestureResponses(t *testing.T) {
    for _, g := range gestures {
        notLosing := !g.Beats(LoseTo(g))
        notWinning := !WinAgainst(g).Beats(g)
        notDrawing := g != DrawWith(g)
        atRightTime := false

        if notLosing || notWinning || notDrawing || atRightTime {
            t.Fail()
        }
        if notLosing {
            t.Fail()
            t.Logf("Incorrect Gesture given: %v does not lose to %v.", LoseTo(g), g)
        }
        if notWinning {
            t.Fail()
            t.Logf("Incorrect Gesture given: %v does not win against %v.", WinAgainst(g), g)
        }
        if notDrawing {
            t.Logf("Incorrect Gesture given: %v does not draw with %v.", DrawWith(g), g)
        }
    }
}

func TestCorrectedRoundParsing(t *testing.T) {
    all := correctlyLoadRounds(strings.NewReader(allRoundVariations))

    pointsForAllGestures := Rock.points + Paper.points + Scissors.points
    pointsForAllRoundOutcomes := PointsForLoss + PointsForTie + PointsForWin
    totalPointsForAllGestureCombos := len(gestures) * (pointsForAllGestures + pointsForAllRoundOutcomes)

    var atp int // Player a's total points
    var ztp int // Player z's total points
    for _, round := range *all {
        atp += round.Outcome().a
        ztp += round.Outcome().z
    }
    if atp != ztp || atp != totalPointsForAllGestureCombos {
        t.Fail()
        t.Logf("Total points was expected to be %d, Player a got %d points while Player z got %d points.",
                totalPointsForAllGestureCombos, atp, ztp)
        t.Log(*all)
    }
}

func TestRPSRules(t *testing.T) {
    tcs := []struct{test, expect bool}{
        {Rock.Beats(Paper), false},
        {Rock.Beats(Scissors), true},
        {Rock.Beats(Rock), false},
        {Paper.Beats(Scissors), false},
        {Paper.Beats(Rock), true},
        {Paper.Beats(Paper), false},
        {Scissors.Beats(Rock), false},
        {Scissors.Beats(Paper), true},
        {Scissors.Beats(Scissors), false},
    }
    for i, tc := range tcs {
        if want, got := tc.expect, tc.test; want != got {
            t.Fail()
            t.Logf("On test case %d: Expected %v but got %v", i, want, got)
        }
    }
}

func TestRoundOutcomes(t *testing.T) {
    tcs := []struct{
        round Round
        outcome RoundOutcome
    }{
        {
            round: Round{Rock, Paper},
            outcome: RoundOutcome{1, 8},
        },
        {
            round: Round{Paper, Scissors},
            outcome: RoundOutcome{2, 9},
        },
        {
            round: Round{Scissors, Rock},
            outcome: RoundOutcome{3, 7},
        },
        {
            round: Round{Rock, Rock},
            outcome: RoundOutcome{4, 4},
        },
        {
            round: Round{Paper, Paper},
            outcome: RoundOutcome{5, 5},
        },
        {
            round: Round{Scissors, Scissors},
            outcome: RoundOutcome{6, 6},
        },
    }

    for _, tc := range tcs {
        got, want := tc.round.Outcome(), tc.outcome
        if got.a != want.a || got.z != want.z {
            t.Logf("Wanted %v outcome, but got %v instead.", want, got)
            t.Fail()
        }
    }
}

