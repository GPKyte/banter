package sack

import (
    "testing"
    "strings"
)
// Sack components will be equal size for now
// Sack components could be unequal to create a pair where no shared items occur
// Sack components cannot have any of the same items
// Duplicates within one component of a sack are acceptable


func TestFindSharedOutlierWithinSackComponents(t *testing.T) {
    tcs := map[Item]*Sack{
        GetItem("a"): New("alhsdfpjbnacwein"),
        GetItem("G"): New("AbCdEfGGhIjKLm"),
        GetItem("b"): New("bbBb"),
        GetItem("z"): New("asdfghjklASDFGzHJKLqwertyuiopQWERTYUIOPzxxcvbnmZXCVBNM"),
    }
    for item, sack := range tcs {
        if item.DifferentThan(sack.Outlier()) {
            t.Fail()
            t.Log(item.String(), sack.String())
        }
        itemNotInBothComponents := !(sack.OneHalf().Contains(item) && sack.OtherHalf().Contains(item))
        if itemNotInBothComponents {
            t.Fail()
            t.Log("One half or the other half does not contain the expected item", sack)
        }
    }
}

func TestNoOutlier(t *testing.T) {
    s := newSack("asdFGHJK")
    o := s.Outlier()
    if o.DifferentThan(NoItem) {
        t.Fail()
        t.Logf("Found an outlier %v where there were none. This was %s", o.String(), strings.ToLower(s.String()))
    }
}

func TestWhetherContainsIsAccurate(t *testing.T) {
    this := newSack("asdfghjklASDFGHJKLZXCVBNMzxcvbnm")
    this.OneHalf().Contains(GetItem("h"))
    this.OtherHalf().Contains(GetItem("M"))
    this.OneHalf().Contains(GetItem("M"))
}

func TestGetPriorityOfOutlier(t *testing.T) {
    tcs := []struct{want int; sack *Sack}{
        {15, New("ASDFGHJKoasdfghjko")},
        {1, New("aBaC")},
        {-1, New("")},
        {-1, New("abc")},
        {52, New("ZabcxyzZ")},
    }
    for _, tc := range tcs {
        if have := tc.sack.Outlier().Type.Priority; have != tc.want {
            t.Fail()
            t.Logf("Wanted to see priority of %d, but have %d instead. From %s",
                    tc.want, have, tc.sack.String())
        }
    }
}

func TestCommonToAllThree(t *testing.T) {
    a := "QWERTYUI"
    b := "IOPBASDF"
    c := "FGHJKYUI"

    common := "I"
    if common != InAll([]string{a,b,c}) {
        t.Fail()
    }

    got := ItemCommonToAllSacks([]*Sack{
        New(a), New(b), New(c)})
    if got != GetItem(common) {
        t.Fail()
    }
}
