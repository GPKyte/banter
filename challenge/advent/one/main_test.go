package main

import (
    "testing"
    "sort"
)

func TestMainExamples(t *testing.T) {
    type testcase struct{
        dataFilePath string
        calorieCount int
        sumThreeCals int
    }
    var examples = []testcase{
        {"testdata/example-input-00.txt", 32180, 76021},
        {"testdata/example-input-01.txt", 24000, 45000},
        {"testdata/example-input-02.txt", 100000016, 110291025},
    }

    for _, tc := range examples {
        elves := generateElvesFromFileReference(tc.dataFilePath)
        ccount := discernHighestCalorieCountAmongTheseElves(elves)
        threeSum := combinedCaloriesOfTopThreeContributors(elves)
        if got, want := ccount, tc.calorieCount; got != want {
            t.Logf(
                "Expected %s file would result in %d calories, but found %d instead.",
                tc.dataFilePath, want, got)
            t.Fail()
        }
        if want, got := tc.sumThreeCals, threeSum; want != got {
            t.Logf(
                "Expected %s file's top three would total %d calories, but found %d instead.",
                tc.dataFilePath, want, got)
            t.Log(elves.TopThreeSnackContributors())
            t.Log(elves)
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

