package monkey

import (
    "testing"

    "github.com/google/go-cmp/cmp"
)

func TestItemSwitching(t *testing.T) {
    aaron := &Monkey{Has: Items{10,20,30}}
    beatrice := &Monkey{Has: Items{13,23,33}}
    candy := &Monkey{Has: Items{15,25,35}}

    aaron.TossAllTo(beatrice)
    candy.TossAllTo(aaron)

    // Care about item order and ownership, but len == 0 is an easy quick check
    if len(candy.Has) != 0  {t.Fatal("candy.Has", candy.Has)}
    want := Items{15,25,35}
    if !cmp.Equal(aaron.Has, want) {
        t.Fail()
        t.Log(cmp.Diff(aaron.Has, want))
        t.Fatal("aaron.Has", aaron.Has)
    }
    if !cmp.Equal(beatrice.Has, Items{13,23,33,10,20,30}) {
        t.Fatal("beatrice.Has", beatrice.Has)
    }
}

func TestOperations(t *testing.T) {
    addSome := NewOperation("new = old + 5")
    double  := NewOperation("new = old + old")
    square  := NewOperation("new = old * old")
    scale   := NewOperation("new = 4 * old")

    reset := WorryLevel
    WorryLevel = 100

    testpath := []struct{operate Operation; expectation int} {
        {double,    200},
        {addSome,   205},
        {addSome,   210},
        {square,    44100},
        {scale,     176400},
        {addSome,   176405},
        {double,    352810},
    }
    for _, test := range testpath {
        test.operate()

        if WorryLevel != test.expectation {
            t.Fatalf("WorryLevel (%d) differs from expectation (%d)", WorryLevel, test.expectation)
        }
    }
    WorryLevel = reset
}

func TestOperationTokenParsing(t *testing.T) {
    got := parseOperationTokens("new = old + 5")
    want := []string{"new", "=", "old", "+", "5"}
    if !cmp.Equal(got, want) {
        t.Fatal(cmp.Diff(got,want))
    }
}

