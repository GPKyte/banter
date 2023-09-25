package main

import (
    "os"
    "log"

    "github.com/GPKyte/banter/challenge/advent/tournament"
)

func main() {
    defer func(){
        if err := recover(); err != nil {
            log.Println(err)
        }
    }()

    if len(os.Args) < 2 {
        panic("Ooph! No filename given for reading data on Rock Paper Scissors competition. Please provide as first argument")
    }
    relFilePath := os.Args[1]
    rpsData, err := os.Open(relFilePath)
    if err != nil {
        panic("Could not open the specified file "+relFilePath)
    }
    defer rpsData.Close()
    
    tourny := tournament.New(rpsData)
    tourny.PlayAll()
    log.Println(tourny)
}
