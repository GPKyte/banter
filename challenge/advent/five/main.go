package main

import (
    "log"

    "github.com/GPKyte/banter/challenge/advent/common"
    "github.com/GPKyte/banter/challenge/advent/crate"
)

func main() {
    puzzleFile := common.OpenFirstArgAsFileReader()
    defer puzzleFile.Close()

    crate.MultiCrateTransfer.Enable()
    cargoShip, instructions := crate.LoadPuzzle(puzzleFile)
    cargoShip.Rearrange(instructions)
    log.Println(cargoShip.TopSummary())
    log.Printf("Review outcome below\n%s", cargoShip)
}

