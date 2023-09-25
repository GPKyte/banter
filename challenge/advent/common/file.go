package common

import (
    "os"
)

func OpenFirstArgAsFileReader() *os.File {
    if len(os.Args) < 2 {
        panic("Specify which file to read from as first argument.")
    }
    fileAtRelativePath := os.Args[1]

    file, err := os.Open(fileAtRelativePath)
    if err != nil && os.IsNotExist(err) {
        panic("File does not exist at "+fileAtRelativePath)
    } else if err != nil {
        panic("Could not open file "+fileAtRelativePath)
    }

    return file
}
