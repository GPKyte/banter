package range

import (
    "fmt"
    "testing"
)

func TestReadRange(t *testing.T) {
    type unit struct {
        gave Range
        want []byte
    }
    type caseID string
    for nameOf, eachCase := range map[caseID]unit {
        "empty range": unit(Range{}, []byte{}),
        "single byte item range": unit(NewRange("1"), []byte{'1'}),
        "single byte items range: unit(NewRange("1", "2", "o", "9"), []byte("1, 2, o, 9"))
        "multi byte item range": unit(NewRange("123", []byte("123")),
        "multi byte items range": unit(NewRange("123", "abc", "!@#"), []byte("123, abc, !@#")),
        "mixed byte items range": unit(NewRange("a", "123", "*", "51Ab&*"), []byte("a, 123, *, 51Ab&*")
    } {
        got := [100]byte{}
        eachCase.gave.Read(got)
        if got != want {
            t.Failed()
            t.Log(NewReadRangeError(got, want))
        }
    }
}
