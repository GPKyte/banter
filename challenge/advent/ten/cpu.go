package main

import (
    "strconv"
    "io"
    "bufio"
)

func New() *CPU {
    return &CPU{
        X: 1,
        Clock: 0,
        XHistory: make([]int, 0),
    }}
type CPU struct {
    X int // register value
    Clock int
    XHistory []int
}
var ClockCyclesOfInterest = []int {
    20, 60, 100, 140, 180, 220,
}

func (c *CPU) Cycle() {
    c.Clock++
    c.XHistory = append(c.XHistory, c.X)
}

func (c *CPU) Execute(ops Operations) {
    for _, o := range ops {
        for tick := 0; tick < o.CycleCost; tick++ {
            c.Cycle()
        }
        o.Action(c)
    }
}

func (c *CPU) RegisterValueDuringCyclesOfInterest() []int {
    vals := make([]int, len(ClockCyclesOfInterest))
    for i := range vals {
        vals[i] = c.XHistory[ClockCyclesOfInterest[i] - 1]
    }
    return vals
}
func (c CPU) SignalStrengthDuring(cycles []int) []int {
    vals := make([]int, len(cycles))
    for i := range vals {
        vals[i] = SignalStrength(cycles[i], c.XHistory[cycles[i] - 1])
    }
    return vals
}
func sum(ofThese []int) int {
    var total int
    for _, each := range ofThese {
        total += each
    }
    return total
}
func SignalStrength(clockCycle, registerValue int) int {
    ss := clockCycle * registerValue
    return ss
}

type Operations []Operation
type Operation  struct{
    Action func(*CPU)
    CycleCost int
}
func Load(from io.Reader) Operations {
    ops := make(Operations, 0)
    oneCycle := 1
    twoCycles := 2

    s := bufio.NewScanner(from)
    for loading := s.Scan(); loading; loading = s.Scan() {
        line := s.Text()
        opname := line[0:4]
        var op Operation

        switch {
        case "addx" == opname:
            value, _ := strconv.Atoi(line[5:])
            addxAction := func(c *CPU) {c.X += value}
            op = Operation{addxAction, twoCycles}
        case "noop" == opname:
            noopAction := func(c *CPU) {}
            op = Operation{noopAction, oneCycle}
        }
        ops = append(ops, op)
    }
    return ops
}
