package main

import (
    "testing"
    "os"
)

func TestExamples(t *testing.T) {
    dir := "testdata/"
    testcases := []struct{filename string; want int} {
        {dir+"example1", 13140},
    }

    for _, tcs := range testcases {
        test := func(t *testing.T) {
            f, _ := os.Open(tcs.filename)
            defer f.Close()

            cpu := New()
            instructionSet := Load(f)
            cpu.Execute(instructionSet)
            sss := cpu.SignalStrengthDuring(ClockCyclesOfInterest)
            total := sum(sss)
            if total != tcs.want {
                t.Fail()
                t.Log(total)
                t.Log(sss)
            }
        }
        t.Run(tcs.filename, test)
    }
}

