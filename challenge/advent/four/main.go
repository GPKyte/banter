package main

import (
    "log"

    "github.com/GPKyte/banter/challenge/advent/assignment"
    "github.com/GPKyte/banter/challenge/advent/common"
)

func main() {
    assignmentData := common.OpenFirstArgAsFileReader()
    defer assignmentData.Close()

    assignmentPairs := assignment.LoadFromPuzzleInput(assignmentData)
    var numberOfFullyRedundantSectionAssignments int
    var numberOfAssignmentPairsWithOverlappingSections int

    for _, ap := range *assignmentPairs {
        if ap.HasFullyRedundantSectionOverlap() {
            numberOfFullyRedundantSectionAssignments++
        }

        if ap.AnyOverlap() {
            numberOfAssignmentPairsWithOverlappingSections++
        }
    }

    log.Printf("Found %d cases where one elf was assigned to at least all the same sections as their partner.\n", numberOfFullyRedundantSectionAssignments)

    log.Printf("Found %d cases where there was any overlapping sections in the elves' assignments per pair.\n", numberOfAssignmentPairsWithOverlappingSections)
}
