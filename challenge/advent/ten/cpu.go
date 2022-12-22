package main

import (
    "io"
)

func New() *CPU {return &CPU{}}
type CPU struct {
    X int // register value
    Clock int    
}
var clockCyclesOfInterest = []int {
    20, 60, 100, 140, 180, 220,
}

func (c *CPU) Cycle() {
    c.Clock++
}

func (c *CPU) Execute(ops Operations) {}

func (c *CPU) RegisterValueDuringCyclesOfInterest() []int {
    return []int{0}
}
func sum(ofThese []int) int {total := 0; return total}
func SignalStrength(clockCycle, registerValue int) int {ss := 1; return ss}

type Operations []Operation
type Operation  struct{
    Action func(*CPU)
    CycleCost int
}
func Load(from io.Reader) Operations {ops := make(Operations, 0); return ops}
