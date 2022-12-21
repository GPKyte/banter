package main

import (
    "fmt"

    "github.com/GPKyte/banter/challenge/advent/common"
)

func main() {
    from := common.OpenFirstArgAsFileReader()
    defer from.Close()

    g := NewGrid()
    g.ApplyMovements(from)
    distinctPlacesRopeLanded := g.SummarizeHistory()

    fmt.Println(len(distinctPlacesRopeLanded))
}
