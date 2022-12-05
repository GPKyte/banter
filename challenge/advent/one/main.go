// The first day of the Advent of Code 2022 Challenge features five elves managing their food on a journey by foot.
// The elves want to distribute food from the elf with the greatest number of calories. To aid them, this program interprets a list of integers describing the calorie content of each food item in each elf's inventory.
// After reading the food details, the greatest total calories carried by one elf will be printed to the console.
package main

import(
    "os"
    "fmt"

    "github.com/GPKyte/banter/challenge/advent/elf"
)

func main() {
    defer func(){
        if err := recover(); err != nil {
            fmt.Println(err)
        }
    }()

    if len(os.Args) < 2 {
        panic("Include filename of list describing calories of the elves' food.")
    }
    relativePathToFoodFile := os.Args[1]

    elves := generateElvesFromFileReference(relativePathToFoodFile)
    fmt.Printf(
        "The elf carrying the most food by calorie is packing a total of %d C, wow!\n",
        discernHighestCalorieCountAmongTheseElves(elves))
}

func generateElvesFromFileReference(f string) *elf.Elves {
    ff, err := os.Open(f)
    defer ff.Close()
    if err != nil {
        panic("Cannot open file to get information about the elves.")
    }

    return elf.New(ff)
}

func discernHighestCalorieCountAmongTheseElves(elves *elf.Elves) int {
    return elves.MostCaloriesCarried()
}
