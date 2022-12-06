package tournament

import (
    "fmt"
)

type Player struct {
    score int // non-negative
}

func (p *Player) String() string {
    return fmt.Sprint(p.Score())
}

func (p *Player) RaiseScore(by int) {
    p.score += by
}

func (p *Player) Score() int {
    return p.score
}
