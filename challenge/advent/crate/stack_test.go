package crate

import (
    "testing"
)

func TestReversalEffect(t *testing.T) {
    s, m, i, l, e := New("[S]"), New("[M]"), New("[I]"), New("[L]"), New("[E]")
    src := &Stack{[]Crate{s, m, i, l, e}, 1}
    dst := newStack(2)

    src.transferCrate(dst) // SMIL : E
    if dst.Top().DifferentThan(e) {
        t.Fail()
        t.Log(dst.Top())
    }
    src.transferCrate(dst) // SMI : EL
    src.transferCrate(dst) // SM : ELI
    if dst.Top().DifferentThan(i) {
        t.Fail()
        t.Log(dst.Top())
    }
    src.transferCrate(dst) // S : ELIM
    src.transferCrate(dst) // : ELIMS

    if src.Height() != 0 {
        t.Fail()
        t.Log(src.Height())
    }
    if dst.Height() != 5 {
        t.Fail()
        t.Log(dst.Height())
    }

    if dst.Top().DifferentThan(s) {
        t.Fail()
        t.Log(dst.Top())
    }
    dst.transferCrate(src) // S : ELIM
    if dst.Top().DifferentThan(m) {
        t.Fail()
        t.Log(dst.Top())
    }
}

func TestGetStack(t *testing.T) {
    one := newStack(1)
    two := newStack(2)
    three := newStack(3)
    all := &Stacks{one, two, three}

    if one != all.Get(1) {t.Fail()}
    if two != all.Get(2) {t.Fail()}
    if three != all.Get(3) {t.Fail(); t.Log(all)}
}

func TestReverseCrates(t *testing.T) {
    a := Crate("[A]")
    b := Crate("[B]")
    c := Crate("[C]")

    all := []Crate{a, b, c}
    rev := []Crate{c, b, a}

    if got := reverseSliceOfCrates(all); got[0] != rev[0] || got[2] != all[0] {
        t.Fail()
        t.Log(got)
    }
}
