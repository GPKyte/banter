package main

import (
    "testing"
    "os"
)

func TestIndexMarkers(t *testing.T) {
    ref := "kzhfinabalzpvmqzohzncljsdlvoq"

    bySizeTwo :=  IndexMarkers(2, ref)
    if len(bySizeTwo) < len(ref) - 1 {
        t.Fail()
        t.Log(len(bySizeTwo))
        t.Log(bySizeTwo)
    }

    if bySizeTwo[0] != 2 {
        t.Fail()
        t.Log(bySizeTwo[0])
    }

    bySizeFour := IndexMarkers(4, ref)
    if len(bySizeFour) < 8  { // stopped counting
        t.Fail()
        t.Log(len(bySizeFour))
        t.Log(bySizeFour)
    }

    if bySizeFour[0] != 4 {
        t.Fail()
        t.Log(bySizeFour[0])
    }

    bySizeNine := IndexMarkers(9, ref)
    if len(bySizeNine) != 1 {
        t.Fail()
        t.Log(len(bySizeNine))
    }

    if bySizeNine[0] != 25 {
        t.Fail()
        t.Log(bySizeNine[0])
    }
}

func TestPuzzleOneExamples(t *testing.T) {
    testfile := "testdata/example-" // prefix
    testcases := []struct{gave string; want int}{
        {gave: testfile+"00", want: 5},
        {gave: testfile+"01", want: 6},
        {gave: testfile+"02", want:10},
        {gave: testfile+"03", want:11},
    }
    for _, tc := range testcases {
        test := func(t *testing.T) {
            from, _ := os.Open(tc.gave)
            defer from.Close()
            ciphertext := LoadPuzzle(from)
            indices := IndexMarkers(FirstMarkerSize, ciphertext)
            firstMarkerIndex := indices[0]
            if firstMarkerIndex != tc.want {
                t.Fail()
                t.Log(ciphertext)
                t.Log(indices)
            }
        }
        t.Run(tc.gave, test)
    }
}

func TestPuzzleTwoExamples(t *testing.T) {
    testfile := "testdata/example-" // prefix
    testcases := []struct{gave string; want int}{
        {gave: testfile+"00", want:23},
        {gave: testfile+"01", want:23},
        {gave: testfile+"02", want:29},
        {gave: testfile+"03", want:26},
        {gave: testfile+"04", want:19},
    }
    for _, tc := range testcases {
        test := func(t *testing.T) {
            from, _ := os.Open(tc.gave)
            defer from.Close()
            ciphertext := LoadPuzzle(from)
            indices := IndexMarkers(MessageMarkerSize, ciphertext)
            firstMarkerIndex := indices[0]
            if firstMarkerIndex != tc.want {
                t.Fail()
                t.Log(ciphertext)
                t.Log(indices)
            }
        }
        t.Run(tc.gave, test)
    }
}

func TestByteSet(t *testing.T) {
    biteme := []byte{'a','a','b','a','t','A','7'}
    setme := ToSet(biteme)

    if len(biteme) <= len(setme) {t.Fail()}
    if UnequalByteSlices(setme, ToSet(setme)) {t.Fail()}
    if len(setme) < 5 {t.Fail()}
}

// UnequalByteSlices is true if the length or the elements of a and b differ
func UnequalByteSlices(a, b []byte) bool {
    var lengthMismatch, unequalElements bool = false, false
    lena := len(a)
    lengthMismatch = lena != len(b)
    
    for i := 0; !lengthMismatch && i < lena && !unequalElements; i++ {
        unequalElements = a[i] != b[i]
    }

    return lengthMismatch || unequalElements
}

func TestMarkerValidation(t *testing.T) {
    is := "abcd"
    not := "abca"

    if !ThisIsAMarker(is) {t.Fail()}
    if ThisIsAMarker(not) {t.Fail()}
}
