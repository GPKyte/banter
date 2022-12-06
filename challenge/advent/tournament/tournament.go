package tournament

import "log"

type Tournament struct {
    ap *Player
    zp *Player
    roundCount int
}
func (t *Tournament) TrackScore(ro RoundOutcome) {
    t.ap.RaiseScore(ro.a)
    t.zp.RaiseScore(ro.z)
    t.roundCount += 1
}

func (t *Tournament) TiedPlayers() []*Player {
    _, _, ties := t.results()
    return ties
}

func (t *Tournament) WinningPlayer() *Player {
    winner, _, _ := t.results()
    return winner
}

func (t *Tournament) LosingPlayer() *Player {
    _, loser, _ := t.results()
    return loser
}

// results summarizes tournament status of the winners, losers, and tied players.
// Rather than repeat the implementation of comparing player scores,
// Create this single method and use return values in aptly named methods.
func (t *Tournament) results() (winner, loser *Player, ties []*Player) { 
    aps := t.ap.Score()
    zps := t.zp.Score()

    if aps == zps {
        winner = nil
        loser = nil
        ties = []*Player{t.ap, t.zp}

    } else if aps < zps {
        winner = t.zp
        loser = t.ap
        ties = nil

    } else if aps > zps {
        winner = t.ap
        loser = t.zp
        ties = nil
 
    } else {
        log.Printf("Reached unlikely conclusion within tournament results, aps: %d, zps: %d", aps, zps)
        winner = nil
        loser = nil
        ties = []*Player{t.ap, t.zp}
    }

    return winner, loser, ties
}

func (t *Tournament) RoundsPlayed() int {
    return t.roundCount
}

