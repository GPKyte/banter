package graph

import (
    "log"
    "os"
)

type Settings struct {
    Verbose bool
    Debug   bool
    Version string
}

var Config = Settings {
    Verbose:    false,
    Debug:      false,
    Version:    Alpha,
}

var Alpha string = "0.1.0"

var verbose = log.New(os.Stdout, "[Graph package]", 0)
func verbosePrintf(s string, args ...interface{}) {
    if Config.Verbose {
        verbose.Printf(s, args...)
    }
}
