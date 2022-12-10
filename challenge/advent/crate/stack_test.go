package crate

import (
    "testing"
)

func TestReversalEffect(t *testing.T) {
    s, m, i, l, e := New("[S]"), New("[M]"), New("[I]"), New("[L]"), New("[E]")
    src := &Stack{[]Crate{s, m, i, l, e}, 1}
    dst := newStack(2)

    src.transferCrate(dst)
    if dst.Top().DifferentThan(s) {t.Fail()}
    src.transferCrate(dst)
    src.transferCrate(dst)
    if dst.Top().DifferentThan(i) {t.Fail()}
    src.transferCrate(dst)
    src.transferCrate(dst)

    if src.Height() != 0 {t.Fail()}
    if dst.Height() != 5 {t.Fail()}

    if dst.Top().DifferentThan(e) {t.Fail()}
    dst.transferCrate(src)
    if dst.Top().DifferentThan(l) {t.Fail()}
}
