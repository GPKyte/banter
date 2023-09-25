package crate

import (
    "fmt"
    "testing"
    "strings"
    "runtime/debug"
)

func TestStackCreationWithCrates(t *testing.T) {
    first := newStack(1)
    second := newStack(2)
    third := newStack(3)
    fourth := newStack(4)
    fifth := newStack(5)
    referenceStacks := &Stacks{first, second, third, fourth, fifth}
    stackCount := len(*referenceStacks)

    four    := "            [I]    "
    three   := "            [H]    " 
    two     := "[E]         [F] [G]"
    one     := "[A]     [B] [C] [D]"
 // zero    := " 1   2   3   4   5 "
    asciiArtCratePlacement := strings.Join([]string{four, three, two, one}, "\n")
    parsedStacks := createStacks(strings.NewReader(asciiArtCratePlacement), stackCount)
    if parsedStacks.String() == asciiArtCratePlacement {
        t.Log("Failed to parse or print these stacks\n"+parsedStacks.String())
    }

    first.Place(New("[A]"))
    first.Place(New("[E]"))
    //second
    third.Place(New("[B]"))
    fourth.Place(New("[C]"))
    fourth.Place(New("[F]"))
    fourth.Place(New("[H]"))
    fourth.Place(New("[I]"))
    fifth.Place(New("[D]"))
    fifth.Place(New("[G]"))

    if have, want := parsedStacks.TopSummary(), referenceStacks.TopSummary(); have != want {
        t.Fail()
        t.Log(have, want)
        t.Log(parsedStacks.String())
    }
    if got := parsedStacks.Get(4).Height(); got != 4 {
        t.Fail()
        t.Log(got)
    }
}

func handleMemoryError(t *testing.T) {
    if err := recover(); err != nil {
        t.Log(err)
        t.Fail()
        debug.PrintStack()
    }
}

func TestInstructionParsing(t *testing.T) {
    instructionInput := `move 1 from 1 to 2
move 2 from 2 to 3
move 1 from 3 to 1`

    transfers := parseInstructions(instructionInput)
    if len(transfers) != 3 {t.Fail()}

    four    := "            [I]    "
    three   := "            [H]    " 
    two     := "[E]         [F] [G]"
    one     := "[A] [J] [B] [C] [D]"
    asciiArtCratePlacement := strings.Join([]string{four, three, two, one}, "\n")
    parsedStacks := createStacks(strings.NewReader(asciiArtCratePlacement), 5)

    for _, t := range transfers {
        for i := 0; i < t.Quantity; i++ {
            src := parsedStacks.Get(t.Source)
            dst := parsedStacks.Get(t.Destination)
            src.transferCrate(dst)
        }
    }

    if parsedStacks.TopSummary() != "JEIG" {t.Fail()}
}

func TestCrateLoading(t *testing.T) {
    four    := "            [I]    "
    three   := "            [H]    "
    two     := "[E]         [F] [G]"
    one     := "[A]     [B] [C] [D]"
    mic := LoadCrates(strings.NewReader(strings.Join([]string{four, three, two, one}, "\n")))

    t.Log(mic)
}

func TestStackNumberParsing(t *testing.T) {
    inOneLine := " 1 2 3 4 5 6 7 8 9 \n"
    beforeInstructions := inOneLine+"move 2 from 5 to 1\n"
    afterStacks := "[A] [B] [C] [D] [E] [F] [G] [H] [I] [J]\n"+inOneLine
    emptyLine := "\n"
    afterStacksAtEOF := fmt.Sprint(afterStacks+inOneLine)
    afterStacksBeforeInstructions := fmt.Sprint(afterStacks+inOneLine+emptyLine+beforeInstructions)

    atEOF := strings.NewReader(afterStacksAtEOF)
    atMore := strings.NewReader(afterStacksBeforeInstructions)
    snAtEOF := parseNumberOfStacks(atEOF)
    snAtMore := parseNumberOfStacks(atMore)

    if snAtEOF != 9 {
        t.Fail()
        t.Log("Trouble parsing stack count from line at EOF")
    }
    if snAtMore != 9 {
        t.Fail()
        t.Log("Trouble parsing stack count in middle of data")
    }

}
