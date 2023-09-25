package crate

import (
    "fmt"
)

type Transfer struct {
    Source int // id of Stack
    Destination int // id of Stack
    Quantity int // number of Crates
}

func (t *Transfer) String() string {
    return fmt.Sprintf("move %d from %d to %d", t.Quantity, t.Source, t.Destination)
}
