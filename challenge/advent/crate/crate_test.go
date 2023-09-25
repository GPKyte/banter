package crate

import (
    "testing"
    "os"
)

func exampleHelper(t *testing.T, fn string) func(t *testing.T) {
    return func(t *testing.T) {
        ex, err := os.Open(fn)
        if err != nil {
            panic(err)
        }
        defer ex.Close()

        stacks, usingTransfers := LoadPuzzle(ex)
        stacks.rearrange(usingTransfers, true)
        ex.Seek(int64(0), 0) // Reset reader
        expectedStacks := loadPuzzleAnswer(ex)

        if have, want := stacks.TopSummary(), expectedStacks.TopSummary(); have != want {
            t.Fail()
            t.Log(have, want)
            t.Log(stacks)
            t.Log(expectedStacks)
        }
    }
}

func TestPuzzleExample(t *testing.T) {

    MultiCrateTransfer.Disable()

    for _, fn := range []string{
        "testdata/example-00.txt",
        "testdata/example-01.txt",
    } {
        t.Run(fn, exampleHelper(t, fn))
    }

    MultiCrateTransfer.Enable()

    for _, fn := range []string{
        "testdata/example-10.txt",
        "testdata/example-11.txt",
    } {
        t.Run(fn, exampleHelper(t, fn))
    }
}

