package main

import (
    "io"
    "os"
    "log"
    "github.com/GPKyte/banter/challenge/advent/monkey"
    "github.com/GPKyte/banter/challenge/advent/common"
)

var Console = log.New(os.Stdout, "Day 11: ", 0)

func main() {
    puzzle := common.OpenFirstArgAsFileReader()
    defer puzzle.Close()
    log.Println(SolvePartOne(puzzle))
}

func SolvePartOne(puzzle io.Reader) int {
    monkey.Config.Verbose = true
    ms := monkey.New(puzzle)
    rounds := 20
    for i := 0; i < rounds; i++ {
        ms.GoARound()
    }
    return ms.MonkeyBusiness()
}

