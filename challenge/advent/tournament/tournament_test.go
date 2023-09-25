package tournament

import (
    "testing"
)

func TestRPSTournamentScoreTracking(t *testing.T) {
    tourny := &Tournament{
        ap: &Player{},
        zp: &Player{},
    }

    aWon := RoundOutcome{a: 7, z: 0}
    zWon := RoundOutcome{a: 0, z: 9}
    draw := RoundOutcome{a: 5, z: 5}

    // aWon +28 +0
    tourny.TrackScore(aWon)
    tourny.TrackScore(aWon)
    tourny.TrackScore(aWon)
    tourny.TrackScore(aWon)
    // draw +20, +20
    tourny.TrackScore(draw)
    tourny.TrackScore(draw)
    tourny.TrackScore(draw)
    tourny.TrackScore(draw)
    // zWon +0 +36
    tourny.TrackScore(zWon)
    tourny.TrackScore(zWon)
    tourny.TrackScore(zWon)
    tourny.TrackScore(zWon)

    tourny.WinningPlayer().Score()
    if tourny.WinningPlayer() != tourny.zp {
        t.Fail()
    }
    tourny.LosingPlayer().Score()
    if tourny.LosingPlayer() != tourny.ap {
        t.Fail()
    }
    tourny.RoundsPlayed()
}
