package elf

import (
    "strings"
    "testing"

    "github.com/GPKyte/banter/challenge/advent/sack"
)

func TestFindAndSumTopThreeCalorieCarriers(t *testing.T) {
    someElvesInventories := `1

    2

    3

    4

    5

    6

    7

    8

    9`
    topThreeCalorieCounts := []int{9, 8, 7}
    elves := New(strings.NewReader(someElvesInventories))
    if want, got := elves.TopThreeSnackContributors().TotalCalorieCount(), sum(topThreeCalorieCounts);
    want != got {
        t.Logf("Expected to get %d, but got %d instead.", want, got)
        t.Fail()
    }
}

func sum(these []int) (total int) {
    for _, num := range these {
        total += num
    }

    return total
}

func TestDiscernMostCaloriesCarriedAmongElves(t *testing.T) {
    someElvesInventories := `1000
2000
1000

3000
5000
1000

5000`
    elves := New(strings.NewReader(someElvesInventories))
    if len(*elves) != 3 {
        t.Logf("Unexpected length of list, want 3 but got %d", len(*elves))
        t.Fail()
    }
    for i, expectedCalorieSum := range []int{4000,9000,5000} {
        if got, want := (*elves)[i].Pack.TotalCalories(), expectedCalorieSum; got != want {
            t.Logf("Expected %d, but got %d instead", want, got)
            t.Fail()
        }
    }
    if elves.MostCaloriesCarried() != 9000 {
        t.Log("Most Calories Carried was incorrect")
        t.Fail()
    }
}

func TestInputParsing(t *testing.T) {
    integersPerLineSeparatedByEmptyLine := `112
332
100

2000
3124

10000



992
2411
0011`

    integerGroups := groupInventoryDescriptionsByElf(strings.NewReader(integersPerLineSeparatedByEmptyLine))
    for i, ig := range integerGroups {
        t.Logf("%d: %s\n", i, ig)
    }
}

func TestFindingTeamBadge(t *testing.T) {
    a := "QWERTYUI"
    b := "IOPBASDF"
    c := "FGHJKYUI"

    threeElves := &Elves{
        &Elf{
            Pack: &Inventory{
                Sack: sack.New(a),
            },
        },
        &Elf{
            Pack: &Inventory{
                Sack: sack.New(b),
            },
        },
        &Elf{
            Pack: &Inventory{
                Sack: sack.New(c),
            },
        },
    }

    if sack.GetItem("I") != threeElves.TeamBadge() {
        t.Fail()
    }
}
