package crate

import(
    "io"
    "bufio"
    "strings"
    "strconv"
    "fmt"
)

func parseInstructions(s string) []Transfer {
    instructions := make([]Transfer, 0)

    i := bufio.NewScanner(strings.NewReader(s))
    i.Split(bufio.ScanLines)
    for ok := i.Scan(); ok; ok = i.Scan() {
        var src, dst, qty int
        _, err := fmt.Sscanf(i.Text(), "move %d from %d to %d", &qty, &src, &dst)
        if err != nil {
            panic(fmt.Errorf("Got the error (%w) and thus could not scan instructions from this input: %s\n", err, i.Text())) 
        }
        instructions = append(instructions, Transfer{
            Source: src,
            Destination: dst,
            Quantity: qty,
        })
    }
    return instructions
}


func LoadCrates(from io.Reader) map[int][]Crate {
    // Track where each crate was found
    properPlace := map[int][]Crate{}
    transferCratesToProperPlace := func(m map[int]Crate) {
        for key, crate := range m {
            if _, ok := properPlace[key]; !ok {
                properPlace[key] = make([]Crate, 0)
            }
            properPlace[key] = append(properPlace[key], crate)
        }
    }

    s := bufio.NewScanner(from)
    s.Split(bufio.ScanLines)
    for dataRemaining := s.Scan(); dataRemaining; dataRemaining = s.Scan() {
        if !strings.Contains(s.Text(), "[") {
            break // Because the input no longer serves this function
        }
        transferCratesToProperPlace(
            extractCratesFromLine(s.Text()))
    }

    return properPlace
}

func extractCratesFromLine(s string) map[int]Crate {
    // cratesInSparseSlice := make([]Crate, len(s))
    cratesByInputLocation := map[int]Crate{}
    
    foundCrateLabel := false
    for i, r := range s {
        if foundCrateLabel {
           cratesByInputLocation[i] = Crate("["+string(r)+"]") 
        }
        foundCrateLabel = r == '['
    }
    return cratesByInputLocation
}

func LoadPuzzle(input io.Reader) (*Stacks, []Transfer) {
    puz := bufio.NewReader(input) // Once read in part, can no longer read said part without resetting
    var stackDrawing, stackLabels, transferInstructions string

    stackDrawing, _ = puz.ReadString('1')
    puz.UnreadByte() // Give back the '1' so that it can be found in the stack labels
    stackLabels, _ = puz.ReadString('\n')
    puz.ReadLine() // Empty line before transferInstructions
    transferInstructions, _ = puz.ReadString('-') // Read until EOF or puzzle answer separator
    transferInstructions = string(transferInstructions[:len(transferInstructions) - 1])

    howMany := parseNumberOfStacks(strings.NewReader(stackLabels))
    stacks := createStacks(strings.NewReader(stackDrawing), howMany)
    transfers := parseInstructions(transferInstructions)

    return stacks, transfers
}

func brokenLoadPuzzle(input io.Reader) (*Stacks, []Transfer) {
    // Worried that sharing io.Reader and varying implementation options of subfunctions could result in odd behavior, as though io.Reader may need to be reset.
    copy := make([]byte, 0)
    io.ReadFull(input, copy)

    howMany := parseNumberOfStacks(strings.NewReader(string(copy)))
    stacks := createStacks(strings.NewReader(string(copy)), howMany)

    r := bufio.NewReader(strings.NewReader(string(copy)))
    aid := []byte("1.")
    r.ReadString(aid[0]) // Stack labels are written above the tranfer instructions
    r.ReadLine() // Clear this line
    i, _ := r.ReadString(aid[1]) // Won't find a period, will go to EOF
    transfers := parseInstructions(i)

    return stacks, transfers
}

func parseNumberOfStacks(input io.Reader) int {
    r := bufio.NewReader(input)
    discarded, err := r.ReadString('1')
    if err != nil {
        panic(
            fmt.Errorf(
                `Could not find start of stack labels because %w
                Discarded this input: %s`,
            err, discarded))
    }
    stackLabelsLineRaw, err := r.ReadString('\n')
    if err != nil {
        panic(fmt.Errorf("Could not read line of Stack Labels: %w", err))
    }
    stackLabelsLine := strings.Trim(stackLabelsLineRaw, " \n")
    stackLabels := strings.Split(stackLabelsLine, " ")
    lastLabel, err := strconv.Atoi(last(stackLabels))
    if err != nil {
        panic(fmt.Errorf("Could not parse last stack label due to: %w", err))
    }

    return lastLabel
}

func last(ss []string) string {
    return ss[len(ss)-1]
}

func brokenParseNumberOfStacks(input io.Reader) int {
    bs := bufio.NewScanner(input)
    ok := bs.Scan()
    asLongAsItTakesToFindTheLineWithStackLabels := func () bool {
        continueWhileNotFound := ok
        moveOnWhenReady := !strings.Contains(bs.Text(), "1") // Finding the number one breaks the for loop
        return continueWhileNotFound && moveOnWhenReady
    }

    for asLongAsItTakesToFindTheLineWithStackLabels() {
        ok = bs.Scan()
    }
    if !ok && !strings.Contains(bs.Text(), "1") {
        panic("Could not find number of stacks. Most recent line scanned is: "+bs.Text())
    }
    stackLabels := strings.Split(bs.Text(), " ")
    stackCount, err := strconv.Atoi(stackLabels[len(stackLabels)-2])
    if err != nil {
        panic(fmt.Errorf("Could not read last stack label from %v and find number of stacks. Encountered this error as a result: %w", stackLabels, err))
    }
    return stackCount
}

func loadPuzzleAnswer(input io.Reader) *Stacks {
    s := bufio.NewReader(input)
    s.ReadString('-') // "---" marks end of input and start of puzzle
    s.ReadString('\n') // "--\n" is leftover after previous Read, clear the line
    // Load remaining content as stacks

    stackDrawing, _ := s.ReadString('1')
    s.UnreadByte()
    stackLabels, _ := s.ReadString('\n')

    howMany := parseNumberOfStacks(strings.NewReader(stackLabels))
    return createStacks(strings.NewReader(stackDrawing), howMany)
}
