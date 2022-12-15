package crate

import (
    "io"
    "os"
    "log"
    "bufio"
    "strings"
    "strconv"
    "fmt"
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
func createStacks(from io.Reader, howMany int) *Stacks {
    stacks := make(Stacks, howMany)
    defer handleMissingStack()
    for sid := 1; sid <= howMany; sid++ {
        s := Stack{
            id: sid,
        }
        stacks[sid-1] = &s
    }

    // Need some knowledge of puzzle input to help put crates onto correct stacks
    distanceBetweenPlacements := CrateSize + len(" ") // 4
    initialOffset := len("[") // 1
    firstId := 1 // stack id is 1-indexed
    defer handleOutOfBounds()

    cratesGroupedByPlacement := LoadCrates(from) 
    for placement, crates := range cratesGroupedByPlacement {
        sid := firstId + (placement - initialOffset) / distanceBetweenPlacements
        if sid > howMany {
            panic(fmt.Sprintf(
                `Out of bounds exception expected: element %d of length %d slice. Due to translation of placement value (%d) to sid:
                review content of range materials %v`,
            sid, howMany, placement, cratesGroupedByPlacement))
        }
        s := stacks.Get(sid) 
        if s == nil {
            panic(fmt.Sprintf("Stack %d of %d has gone missing in action. Stack referencess: %v", sid, howMany, stacks))
        }
        s.crates = reverseSliceOfCrates(crates)
    }
    return &stacks
}

func handleMissingStack() {
    if err := recover(); err != nil {
        panic(err)
    }
}

func handleOutOfBounds() {
    if err := recover(); err != nil {
        panic(err)
    }
}


func reverseSliceOfCrates(given []Crate) []Crate {
    ccount := len(given)
    reversed := make([]Crate, 0, ccount)

    for i := ccount - 1; i >= 0; i-- {
        reversed = append(reversed, given[i])
    }
    return reversed
}

func brokenCreateStacks(from string, howMany int) *Stacks {
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

func (s *Stack) String() string {
    return fmt.Sprint(*s)
}

func (s *Stacks) quickString() string {
    return fmt.Sprint(*s)
}

func (s *Stacks) fullString() string {
    stacks := *s

    var maxHeight int
    for _, each := range *s {
        if each.Height() > maxHeight {
            maxHeight = each.Height()
        }
    }

    buffer := make([]string, 0, maxHeight)
    for everyHeight := maxHeight; everyHeight > 0; everyHeight-- {
        cratesAtThisHeight := make([]string, 0, len(stacks))
        for _, eachStack := range stacks {
            maybeEmptyCrateLabel := eachStack.LookAt(everyHeight)
            cratesAtThisHeight = append(cratesAtThisHeight, maybeEmptyCrateLabel)
        }
        buffer = append(buffer, strings.Join(cratesAtThisHeight, " "))
    }
    stackLabels := make([]string, len(stacks))
    for i, each := range stacks {
        label := strconv.Itoa(each.id)
        stackLabels[i] = label
    }
    buffer = append(buffer, " "+strings.Join(stackLabels, "   ")+" ")
    return strings.Join(buffer, "\n")
}

func (s *Stacks) String() string {
    return s.fullString()
}

func (s *Stacks) Get(id int) *Stack {
    // id is 1-indexed but underlying array is of course 0-indexed
    return (*s)[id-1] // Really should protect this from OOB, but I don't have an alternative to offer when that case arrises. Better to catch the error elsewhere perhaps.
}

func (s *Stack) Top() Crate {
    var r Crate
    if s.Height() > 0 {
        r = s.crates[s.Height()-1]
    } else {
        panic("Stack is empty, nothing on top")
    }

    return r
}

func (s *Stack) Place(c Crate) {
    s.crates = append(s.crates, c)
}

func (s *Stack) PickUp() Crate {
    topCrate := s.Top()
    s.crates = s.crates[:s.Height()-1] // [c c c c c][:5-1] -> [c c c c] Because end range is exclusive
    return topCrate
}

func (s *Stacks) Transfer(t Transfer) {
    src := s.Get(t.Source)
    dst := s.Get(t.Destination)
    var i int
    defer func(){
        if err := recover(); err != nil {
            panic(fmt.Errorf("%w\nEncountered during move %d, of %d for transfer %v", err, i, t.Quantity, t))
        }
    }()
    for i = 0; i < t.Quantity; i++ {
        src.transferCrate(dst)
    }
}

func (s *Stack) transferCrate(to *Stack) {
    defer func(){
        if err := recover(); err != nil {
            panic(fmt.Errorf("While transferring crate from stack %v to stack %v encountered error %w", s, to, err))
        }
    }()
    to.Place(s.PickUp())
}

func (s *Stack) Height() int {
    return len(s.crates)
}

func (s *Stacks) TopSummary() string {
    defer handlePanicFromTopOfEmptyStack()
    all := make([]string, 0, len(*s))
    for _, each := range *s {
        if each.Empty() {
            continue
        }
        all = append(all, each.Top().ShortString())
    }
    return strings.Join(all, "")
}

func (s *Stack) LookAt(height int) string {
    if s.Height() < height {
        return "   " // CrateSize
    }
    return string(s.crates[height-1])
}

func (s *Stack) Empty() bool {
    return s.Height() == 0
}

func handlePanicFromTopOfEmptyStack() {
    if err := recover(); err != nil {
        panic(err)
    }
}

func (s *Stacks) Rearrange(byThese []Transfer) {
    // doNotPauseAndLogEveryTransferOutcome := bool(false)
    logEveryStep := bool(true)
    s.rearrange(byThese, logEveryStep)
}

func (s *Stacks) rearrange(byThese []Transfer, stepwiseLoggingEnabled bool) {
    stepReference := 1
    defer func() {
        if err := recover(); err != nil {
            panic(fmt.Errorf("%w\nEncountered at step %d, for transfer %v.", err, stepReference, byThese[stepReference]))
        }
    }()

    for _, t := range byThese {
        s.Transfer(t)

        if stepwiseLoggingEnabled {
            fmt.Println(t.String())
            fmt.Println(s)
        }
        stepReference++
    }
}

func pauseForUser() {
    p := bufio.NewReader(os.Stdin)
    pause := true
    for pause {
        if s, _ := p.ReadString('\n'); s == "" {
            pause == false
        }
    }
}
