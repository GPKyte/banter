package main

import (
    "testing"
    "sort"
)

func TestMainExamples(t *testing.T) {
    type testcase struct{
        dataFilePath string
        calorieCount int
    }
    var examples = []testcase{
        {"testdata/example-input-00.txt", 32180},
        {"testdata/example-input-01.txt", 24000},
        {"testdata/example-input-02.txt", 100000016},
    }

    for _, tc := range examples {
        elves := generateElvesFromFileReference(tc.dataFilePath)
        ccount := discernHighestCalorieCountAmongTheseElves(elves)
        if got, want := ccount, tc.calorieCount; got != want {
            t.Logf(
                "Expected %s file would result in %d calories, but found %d instead.",
                tc.dataFilePath, want, got)
            t.Fail()
        }
    }
}

func TestLargerInput(t *testing.T) {
    elves := generateElvesFromFileReference("testdata/full-input.txt")
    calCounts := make([]int, 0, len(*elves))
    for _, e := range *elves {
        calCounts = append(calCounts, e.Pack.TotalCalories())
    }
    sort.Ints(calCounts)
    t.Log(calCounts)
}

