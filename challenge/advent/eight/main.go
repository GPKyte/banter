package main

import (
    "fmt"

    "github.com/GPKyte/banter/challenge/advent/common"
    "github.com/GPKyte/banter/challenge/advent/eight/treehouse"
)

func main() {
    inputFile := common.OpenFirstArgAsFileReader()
    defer inputFile.Close()

    forest := treehouse.NewForest(inputFile)

    fmt.Println(forest)
    fmt.Println(forest.CountVisible())
    fmt.Println(forest.FindMostScenicTreeScore())
}
