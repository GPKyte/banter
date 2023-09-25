package assignment

import (
    "testing"
    "strings"
)

func TestImportFromString(t *testing.T) {
    input := "2-4,8-9"
    p := newPair(input)

    if p.a.s.Len() != 3 || p.z.s.Len() != 2 {t.Fail()}

    severalInputs := `0-9,5-9
1-3,3-8
4-5,9-20
3-6,1-10`
    severalAps := LoadFromPuzzleInput(strings.NewReader(severalInputs))
    if len(*severalAps) != 4 {t.Fail()}
}

func TestRedundancyCheck(t *testing.T) {
    with := newPair("20-30,10-50")
    without := newPair("1-5,4-9")
    equal := newPair("1-8,1-8")

    if !with.HasFullyRedundantSectionOverlap() {t.Fail()}
    if without.HasFullyRedundantSectionOverlap() {t.Fail()}
    if !equal.HasFullyRedundantSectionOverlap() {t.Fail()}
}

func TestOverlap(t *testing.T) {
    with := newPair("20-30,10-50")
    without := newPair("1-5,4-9")
    equal := newPair("1-8,1-8")
    other := newPair("10-35,49-81")

    if !with.AnyOverlap() {t.Fail()}
    if !without.AnyOverlap() {t.Fail()}
    if !equal.AnyOverlap() {t.Fail()}
    if other.AnyOverlap() {t.Fail()}
}
