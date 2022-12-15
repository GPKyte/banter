package main

import (
    "log"

    "github.com/GPKyte/banter/challenge/advent/elf"
    "github.com/GPKyte/banter/challenge/advent/sack"
    "github.com/GPKyte/banter/challenge/advent/common"
)

func main() {
    puzzleInput := common.OpenFirstArgAsFileReader()
    defer puzzleInput.Close()
    sacks := sack.ReadSackDescriptions(puzzleInput)

    // In part one, the goal is to find an Item which appears in each component of a each sack
    var sumOfAllSackOutlierPriorities int
    for _, s := range *sacks {
        sumOfAllSackOutlierPriorities += s.Outlier().Priority()
    }

    log.Printf("Total of outlier priorities is %d.\n",
                sumOfAllSackOutlierPriorities)

    // In part two, groups of three elves, each with a sack will find the Item common amongst all of their sacks
    groupsOfThree := elf.EquipSacksOntoElvenGroupsOfSize(3, sacks)
    var sumOfBadgePriorities int
    log.Println(len(*groupsOfThree), "groups and", len(*sacks), "sacks")
    for _, g := range *groupsOfThree {
        sumOfBadgePriorities += g.TeamBadge().Priority()
    }

    log.Printf("Total badge priority is %d.\n",
                sumOfBadgePriorities)
}
