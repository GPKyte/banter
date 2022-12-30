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

    beatrice.TossTo(aaron)
    beatrice.TossTo(candy)
    beatrice.TossTo(aaron)
    candy.TossTo(aaron)
}

func TestOperations(t *testing.T) {
    addSome := NewOperation("new = old + 5")
    double  := NewOperation("new = old + old")
    square  := NewOperation("new = old * old")
    scale   := NewOperation("new = 4 * old")

    testpath := []struct{item Item; operate Operation; expectation int} {
        {100,    double,    200},
        {200,    addSome,   205},
        {205,    addSome,   210},
        {210,    square,    44100},
        {44100,  scale,     176400},
        {176400, addSome,   176405},
        {176405, double,    352810},
    }
    for _, test := range testpath {
        test.operate(&test.item)

        if test.item != Item(test.expectation) {
            t.Fatalf("Worry Level (%d) differs from expectation (%d)", test.item, test.expectation)
        }
    }
}

func TestOperationTokenParsing(t *testing.T) {
    got := parseOperationTokens("new = old + 5")
    want := []string{"new", "=", "old", "+", "5"}
    if !cmp.Equal(got, want) {
        t.Fatal(cmp.Diff(got,want))
    }
}

func TestChoiceUsage(t *testing.T) {
    zed := Monkey{Decide: func(int) string {return "0"}}
    choose1or3ifdiv13 := NewChoice(parseChoice(`  Test: divisible by 13
    If true: throw to monkey 1
    If false: throw to monkey 3`))

    trinity := Monkey{Decide: choose1or3ifdiv13}

    divable := 39
    if zed.Decide(divable) != "0" {t.Fail()}
    if trinity.Decide(divable) != "3" {t.Fail()}
    if trinity.Decide(divable+1) != "1" {t.Fail()}
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

func TestRoundOne(t *testing.T) {
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

    ms.GoARound()
    mone, err := ms.Target("1")
    if err != nil {
        t.Fatal(err)
    }
    mzed, err := ms.Target("0")
    if err != nil {
        t.Fatal(err)
    }
    if !cmp.Equal(mzed.Has, Items{0, 23, 27, 26}) {
        t.Fail()
        t.Log(mzed)
    }
    if !cmp.Equal(mone.Has, Items{2080, 25, 167, 207, 401, 1046}) {
        t.Fail()
        t.Log(mone)
    }
}

func TestDivBehavior(t *testing.T) {
    var dividee int = 20
    var divisor int = 3

    if dividend := dividee / divisor; dividend != 6 {
        t.Log(dividend)
    }
}
