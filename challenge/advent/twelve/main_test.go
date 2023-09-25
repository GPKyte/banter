package main

import (
    "testing"

    "os"
    "github.com/GPKyte/banter/challenge/advent/graph"
)

var TestDir = "testdata/"
func TestFindRoute(t *testing.T) {
    testcases := []struct{
        problemFileName string
        solutionFileName string
        solutionSummary int
    }{
        {// Example 00 can find a path in 31 steps
            TestDir+"ex00",
            TestDir+"ans00",
            31,
        },
    }

    for _, tc := range testcases {
        // File keeping
        given, err := os.Open(tc.problemFileName)
        if err != nil {
            t.Fatal(err)
        }
        defer given.Close()
        want, err := os.Open(tc.solutionFileName)
        if err != nil {
            t.Fatal(err)
        }
        defer want.Close()
        
        // Begin
        shortAnswer := SolvePartOne(given)
        if shortAnswer != tc.solutionSummary {
            t.Log(shortAnswer, "vs.", tc.solutionSummary)
        }
    }
}

func TestQueueOperationsAndBoundaries(t *testing.T) {
    Q := NewCoordinateQueue()
    const (
        dq = iota
        nq
    )
    type step struct {
        opcode  int // To remove redundancy in writing out dq or nq
        coord   graph.Coordinate
        length  int // After the operation completes
        okay    bool // When empty before Dequeue, not okay (false condition)
    }

    testproceedure := []step{
        {
            opcode: dq,
            coord:  graph.Coordinate{},
            length: 0,
            okay:   false,
        },
        {
            opcode: dq,
            coord:  graph.Coordinate{},
            length: 0,
            okay:   false,
        },
        {
            opcode: nq,
            coord:  graph.Coordinate{X: 2, Y: 2},
            length: 1,
        },
        {
            opcode: dq,
            coord:  graph.Coordinate{X: 2, Y: 2},
            length: 0,
            okay:   true,
        },
        {
            opcode: nq,
            coord:  graph.Coordinate{X: 2, Y: 3},
            length: 1,
        },
        {
            opcode: nq,
            coord:  graph.Coordinate{X: 1, Y: 2},
            length: 2,
        },
        {
            opcode: nq,
            coord:  graph.Coordinate{X: 3, Y: 2},
            length: 3,
        },
        {
            opcode: nq,
            coord:  graph.Coordinate{X: 2, Y: 2},
            length: 4,
        },
        {
            opcode: dq,
            coord:  graph.Coordinate{X: 2, Y: 3},
            length: 3,
        },
        {
            opcode: dq,
            coord:  graph.Coordinate{X: 1, Y: 2},
            length: 2,
        },
        {
            opcode: dq,
            coord:  graph.Coordinate{X: 3, Y: 2},
            length: 1,
        },
        {
            opcode: dq,
            coord:  graph.Coordinate{X: 2, Y: 2},
            length: 0,
            okay:   true,
        },
        {
            opcode: dq,
            coord:  graph.Coordinate{},
            length: 0,
            okay:   false,
        },
    }
    for each, step := range testproceedure {
        t.Logf("Step %d: %v", each, step)

        switch step.opcode {
        case nq:
            Q.NQ(step.coord)
        case dq:
            c, ok := Q.DQ()
            if c != step.coord {
                t.Fail()
                t.Logf("\t%02d: Expected coordinate %v but got %v instead.\n", each, step.coord, c)
            }
            if ok != step.okay {
                t.Fail()
                t.Logf("\t%02d: May have issue with stopping Dequeue at appropriate time", each)
            }
        }

        if len(*Q) != step.length {
            t.Fail()
            t.Logf("\t%02d: Expected length of %d, but found %d in %v.\n",
                    each, step.length, len(*Q), *Q)
        }
    }
}

func TestLabelToHeightConversion(t *testing.T) {
    // Special cases: starting point and end goal

    testcases := [] struct {
        give string
        want int
    } {
        {give: "a", want:  1},
        {give: "S", want:  1},
        {give: "b", want:  2},
        {give: "h", want:  8},
        {give: "o", want: 15},
        {give: "z", want: 26},
        {give: "E", want: 26},
        {give: "e", want:  5},
        {give: "s", want: 19},
    }

    for each, tc := range testcases {
        t.Logf("Convert %s to %d\n", tc.give, tc.want)
        want := tc.want
        got := LabelToHeight(tc.give)

        if want != got {
            t.Fail()
            t.Logf("\t%02d: Wanted %d, but got %d.\n", each, want, got)
        }
    }
}
