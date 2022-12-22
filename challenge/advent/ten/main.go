package main

import (
    "io"
    "fmt"
    "github.com/GPKyte/banter/challenge/advent/common"
)

func main() {
    instructionFile := common.OpenFirstArgAsFileReader()
    defer instructionFile.Close()

    fmt.Println(solvePuzzle(instructionFile))
}

func solvePuzzle(from io.Reader) int {
    ops := Load(from)

    cpu := New()
    cpu.Execute(ops)
    sss := cpu.SignalStrengthDuring(ClockCyclesOfInterest)
    return sum(sss)
}
