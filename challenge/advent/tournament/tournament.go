package tournament

import (
    "log"
    "fmt"
    "io"
)

type Tournament struct {
    ap *Player
    zp *Player
    rounds *Rounds
    roundsPlayed int
}
func (t *Tournament) TrackScore(ro RoundOutcome) {
    t.ap.RaiseScore(ro.a)
    t.zp.RaiseScore(ro.z)
    t.roundsPlayed += 1
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
    return t.roundsPlayed
}

func (t *Tournament) PlayAll() {
    for _, r := range *t.rounds {
        t.TrackScore(r.Outcome())
    }
}

func (t *Tournament) String() string {
    var msg string

    if t.RoundsPlayed() <= 0 {
        msg = "Tournament has not begun."
    } else if t.RoundsPlayed() < len(*t.rounds) {
        msg = "Tournament is in progress."
    } else if t.RoundsPlayed() >= len(*t.rounds) {
        if t.TiedPlayers() != nil {
            msg = fmt.Sprintf(
                "Tournament is over. Players tied with a whopping %d points each",
                t.TiedPlayers()[0].Score())
        } else {
            msg = fmt.Sprintf(
                "Tournament is over. Winning player %s scored %d points compared to their opponent %s who scored %d points.\n",
                t.WinningPlayer(), t.WinningPlayer().Score(), t.LosingPlayer(), t.LosingPlayer().Score())
        }
    } else {
        msg = `Tournament has reached the end of its Life Cycle and guarantees of observable behavior are no longer active.`
    }

    return msg
}

func New(from io.Reader) *Tournament {
    return &Tournament{
        ap: &Player{name: "Elfuli"},
        zp: &Player{name: "Elfammer"},
        rounds: correctlyLoadRounds(from),
        roundsPlayed: 0,
    }
}

