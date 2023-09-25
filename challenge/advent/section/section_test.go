package section

import (
    "testing"
)

func TestContains(t *testing.T) {
    small := Range(2, 4)
    big := Range(1, 5)
    smaller := Range(2,2)

    if !small.Contains(smaller) {t.Fail()}
    if !big.Contains(small) {t.Fail()}
    if small.Contains(big) {t.Fail()}
    if !EitherContainsTheOther(small, big) {t.Fail()}
}

func TestNeitherContains(t *testing.T) {
    partialOverlapWithAmer := Range(7,20)
    amer := Range(5,9)
    completelyDifferentThanAmer := Range(70, 90)
    empty := Range(2,1)
    
    if EitherContainsTheOther(amer, completelyDifferentThanAmer) {t.Fail()}
    if EitherContainsTheOther(amer, partialOverlapWithAmer) {t.Fail()}
    if EitherContainsTheOther(amer, empty) {t.Fail()}
}
