package tournament

import (
    "testing"
)

func TestPlayerScoreCases(t *testing.T) {
    p := &Player{}

    p.RaiseScore(0)
    p.Score()
    p.RaiseScore(3)
    p.Score()
    p.RaiseScore(6)

    if p.Score() != 9 {
        t.Fail()
    }
}
