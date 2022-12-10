package crate

import (
    "log"
    "bufio"
    "strings"
)

type Stacks []*Stack
type Stack struct {
    crates []Crate
    id  int
}

func newStack(ref int) *Stack {
    return &Stack{
        crates: make([]Crate, 0),
        id: ref,
    }
}

// createStacks takes an ascii art drawing of the stacks
// and returns an ordered collection of stacks filled with some number of labelled crates
func createStacks(from string, howMany int) *Stacks {
    s := make(Stacks, 0) // return value
    for i := 0; i < howMany; i++ {
        s = append(s, newStack(i+1))
    }
    // chooseStack by using rune index of crate in the drawing
    // to return the id of the stack the crate will be placed into
    chooseStack := func(placement int) int {
        // 1   5   9   13  17   index of column letter   
        //[A]         [I] [L]
        //[B]     [G] [H] [K]
        //[C] [D] [E] [F] [J]
        // 1   2   3   4   5    correlates to stack id
        // 0   1   2   3   4    == grid[i]
        return 1 + (placement) / 4
    }

    grid := make([][]Crate, 0) // Will be easier to reference from this compared to a slice or string

    defer func() {
        if err := recover(); err != nil {
            // Good chance we will encounter a OOB error within this method
            log.Println(grid)
            log.Fatal(err)
        }
    }()
    // storeNewCratesInGrid to hide the nested for loop
    storeNewCratesInGrid := func(s string) {
        sub := s[:] // Trim Left portion after adding crate to grid so that strings.Index() works again
        offset := 0 // Keep track of this and use to ensure chooseStack() finds the right one

        for ; strings.Contains(sub, "[") ; {
            placement := offset + strings.Index(sub, "[")
            grid[len(grid)-1][chooseStack(placement)-1] = New(sub[placement:placement+CrateSize])
            sub = sub[placement+CrateSize:]
            offset += placement+CrateSize
        }
    }
    fillRowWithNoCrate := func() {
        for i := range grid[len(grid)-1] {
            grid[len(grid)-1][i] = Crate("")
        }
    }

    // The input string is read from top to bottom, though the reverse may simplify this operation.
    scnr := bufio.NewScanner(strings.NewReader(from))
    scnr.Split(bufio.ScanLines)
    for eachLine := scnr.Scan(); eachLine; eachLine = scnr.Scan() {
        grid = append(grid, make([]Crate, howMany))
        fillRowWithNoCrate() // Default data, counterpart of 0, "", or nil

        line := scnr.Text()
        storeNewCratesInGrid(line)
    }

    // The grid was created for this moment
    for level := 0; level < len(grid); level++ {
        for sid := 1; sid <= len(s); sid++ {
            c := grid[level][sid-1]
            if NoCrate.SameAs(c) {
                continue
            }
            s.Get(sid).Place(c)
        }
    }
    return &s
}

func (s *Stacks) Get(id int) *Stack {
    return (*s)[id] // Really should protect this from OOB, but I don't have an alternative to offer when that case arrises. Better to catch the error elsewhere perhaps.
}

func (s *Stack) Top() Crate {
    return s.crates[s.Height()-1]
}

func (s *Stack) Place(c Crate) {
    s.crates = append(s.crates, c)
}

func (s *Stack) PickUp() Crate {
    topCrate := s.Top()
    s.crates = s.crates[:s.Height()-2] // [c c c c c][:5-2] -> [c c c c]
    return topCrate
}

func (s *Stack) transferCrate(to *Stack) {
    to.Place(s.PickUp())
}

func (s *Stack) Height() int {
    return len(s.crates)
}

func (s *Stacks) TopSummary() string {
    all := make([]string, 0, len(*s))
    for _, each := range *s {
        all = append(all, each.Top().ShortString())
    }
    return strings.Join(all, "")
}
