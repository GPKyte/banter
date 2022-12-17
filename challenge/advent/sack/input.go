package sack

import (
    "io"
    "text/scanner"
)
// ReadSackDescriptions to return a list of Sacks
// Given a file reader or a string reader with even numbers of english letters on every line
// Creates a series of Sacks made of two equal size components full of items.
func ReadSackDescriptions(From io.Reader) *[]*Sack {
    sacks := make([]*Sack, 0)

    var s scanner.Scanner
    s.Init(From)
    for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
        sacks = append(sacks, New(s.TokenText()))
    }

    return &sacks
}

// splitLineInHalf gives two strings of equal length when given string is even.
// But when given string is odd, the median is included in the left/first half.
func splitLineInHalf(s string) (string, string) {
    len := len(s)
    extraStaysInFirstHalf := len % 2
    halfway := len / 2
    firstAnd, secondHalf := splitStringAtIndex(s, halfway + extraStaysInFirstHalf)
    return firstAnd, secondHalf
}

func splitStringAtIndex(s string, index int) (string, string) {
    first := s[:index]
    second := s[index:]
    return first, second
}
