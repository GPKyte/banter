package monkey

import (
    "os"
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

func TestChoiceUsage(t *testing.T) {
    zed := Monkey{Decide: func() string {return "0"}}
    choose1or3ifdiv13 := NewChoice(parseChoice(`  Test: divisible by 13
    If true: throw to monkey 1
    If false: throw to monkey 3`))

    trinity := Monkey{Decide: choose1or3ifdiv13}

    if zed.Decide() != "0" {t.Fail()}
    if trinity.Decide() != "3" {t.Fail()}
    reset := WorryLevel
    WorryLevel = 39
    if trinity.Decide() != "1" {t.Fail()}
    WorryLevel = reset
}

func TestMakingMonkies(t *testing.T) {
    exf, err := os.Open("example-00")
    if err != nil {
        t.Fatal("Could not open example file", err)
    }
    defer exf.Close()

    defer func() {
        if r := recover(); r != nil {
            t.Fatal(r)
        }
    }()
    ms := New(exf)
    if len(ms.Group) != 4 {
        t.Fail()
        t.Log(ms.Group)
    }

    mtwo, err := ms.Target("2")
    if err != nil {
        t.Fail()
        t.Log(ms.Group)
        t.Log(err)
    }
    if !cmp.Equal(mtwo.Has, Items{79, 60, 97}) {
        t.Fail()
        t.Log(mtwo)
    }
}

