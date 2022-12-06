package tournament

import (
    "fmt"
)

type Player struct {
    score int // non-negative
    name string
}

func (p *Player) String() string {
    return fmt.Sprint(p.name)
}

func (p *Player) RaiseScore(by int) {
    p.score += by
}

func (p *Player) Score() int {
    return p.score
}
