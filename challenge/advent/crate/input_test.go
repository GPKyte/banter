package crate

import (
    "testing"
    "strings"
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
    parsedStacks := createStacks(asciiArtCratePlacement, stackCount)

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

    if parsedStacks.TopSummary() != referenceStacks.TopSummary() {t.Fail()}
    if parsedStacks.Get(4).Height() != 4 {t.Fail()}
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
    parsedStacks := createStacks(asciiArtCratePlacement, 5)

    for _, t := range transfers {
        for i := 0; i < t.Quantity; i++ {
            src := parsedStacks.Get(t.Source)
            dst := parsedStacks.Get(t.Destination)
            src.transferCrate(dst)
        }
    }

    if parsedStacks.TopSummary() != "JEIG" {t.Fail()}
}
