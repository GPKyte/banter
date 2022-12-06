// package elf contains traits and methods related to the representation of this character
package elf

import (
    "io"
    "sort"
    "fmt"
    "strings"

    "github.com/GPKyte/banter/challenge/advent/food"
)

type Elves []*Elf

// MostCaloriesCarried by one elf returns the sum of all food calories in elf's pack
func (e *Elves) MostCaloriesCarried() int {
    // Collection is not sorted so use sequential search for find max
    var c, mc int

    // Improve by creating a routine for comparing against record
    // And routine(s) for calculating total calories of each elf
    for _, every := range *e {
        c = every.Pack.TotalCalories()

        if c > mc {
            mc = c
        }
    }

    return mc
}

type Elf struct {
    Pack *Inventory
}

func (e *Elf) String() string {
    return fmt.Sprint(e.Pack.TotalCalories())
}

// New expedition of Elves based on a description of each elf, e.g. their inventory details
func New(source io.Reader) *Elves {
    elves := make(Elves, 0)

    for _, elfInventoryString := range groupInventoryDescriptionsByElf(source) {
        foodCollection := food.New(strings.NewReader(elfInventoryString))
        i := &Inventory{Foods: foodCollection}
        elves = append(elves, &Elf{Pack: i})
    }

    return &elves
}

func (e *Elves) Len() int {
    return len(*e)
}

func (e *Elves) Less(i, j int) bool {
    return (*e)[i].Pack.TotalCalories() < (*e)[j].Pack.TotalCalories()
}

func (e *Elves) Swap(i, j int) {
    holdme := (*e)[i]
    (*e)[i] = (*e)[j]
    (*e)[j] = holdme
}

func (e *Elves) TopThreeSnackContributors() *Elves {
    sort.Sort(e)
    topThree := (*e)[len(*e)-3:]
    return &topThree
}

func (e *Elves) TotalCalorieCount() (total int) {
    for _, elf := range *e {
        total += elf.Pack.TotalCalories()
    }
    return total
}

// groupInventoryDescriptionsByElf to align with design of smaller units, e.g. food type
// and to ease edge case of elf with empty inventory/no food
func groupInventoryDescriptionsByElf(source io.Reader) []string {
    // Consider using os package or some other combo to get the best line separator
    return strings.Split(readAll(source), "\n\n")
}

func readAll(ofThis io.Reader) string {
    var all []byte = make([]byte, 0)
    var b = make([]byte, 1)
    _, err := io.ReadFull(ofThis, b)
    for ;err == nil; _, err = io.ReadFull(ofThis, b) {
        all = append(all, b[0])
    }
    return string(all)
}
