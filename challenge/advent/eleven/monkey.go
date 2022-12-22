package monkey

import (
    "io"
    "strings"
    "strconv"
)

func New(from io.Reader) *Monkeys {return nil}

type Monkey struct {
    ID string
    Has Items
    HandleItems Operation
    Decide Choice
    Group *Monkeys
}

type Monkeys struct {
    Group []*Monkey
}

type Choice func() string
func NewChoice(test func() bool, monkeyIdIfTestPasses, monkeyIdIfTestFails string) Choice {
    return func() string {
        if test() {
            return monkeyIdIfTestPasses
        } else {
            return monkeyIdIfTestFails
        }
    }
}
func parseChoice(from string) (func() bool, string, string) {
    // Given three lines, first is test, second is when true, third when false
    lines := strings.Split(from, "\n")
    if len(lines) != 3 {
        panic("Can not parse Choice from "+from)
    }
    first := lines[0]
    testdescription := strings.Trim(first[strings.Index(first, ":") + 1:], " ") // "..." from "  Test: ..."
    testparts := strings.Split(testdescription, " ")
    if testparts[0] != "divisible" {
        panic("Can only test divisibility. But this test description says: "+testdescription)
    }
    test := func() bool {
        divisor, err := strconv.Atoi(testparts[len(testparts)-1])
        if err != nil {
            panic(err)
        }
        return WorryLevel % divisor == 0
    }

    second := lines[1]
    ifTrueLabel := second[strings.IndexAny(second, "0123456789"):]
    third := lines[2]
    ifFalseLabel := third[strings.IndexAny(third, "0123456789"):]

    return test, ifTrueLabel, ifFalseLabel
}

var WorryLevel int = 1
type Operation func()
func NewOperation(from string) Operation {
    // Operation: var = Expression
    // Expression: var | var Operator Expression
    // Two vars are recognized, old, new
    // new and old are both &WorryLevel, just at different times
    // Operators include + and *, but could include /, -, or others
    var tokens []string = parseOperationTokens(from)

    operands, err := chooseOperands(tokens[2], tokens[4])
    if err != nil {
        panic(err)
    }
    operator  := chooseOperator(tokens[3]) // + * / -

    operation := func() {
        WorryLevel = operator(*operands[0], *operands[1])
    }
    return operation
}

func chooseOperands(from ...string) ([]*int, error) {
    ops := make([]*int, 0)

    for _, label := range from {
        var o *int

        switch label {
        case "new":
            o = &WorryLevel
        case "old":
            o = &WorryLevel
        default:
            num, err := strconv.Atoi(label)
            if err != nil {
                return nil, err
            }
            o = &num
        }
        ops = append(ops, o)
    }
    return ops, nil
}

func chooseOperator(from string) operator {
    switch from {
    case "+":
        return add
    case "*":
        return multiply
    }
    return noop
}

type operator func(int, int) int
var add = func(a, b int) int {return a + b}
var multiply = func(a, b int) int {return a * b}
var noop = func(a, b int) int {return 1}

func parseOperationTokens(from string) []string {
    tokens := strings.Split(strings.Trim(from, " "), " ")
    return tokens
}

func (m *Monkey) TossAllTo(that *Monkey) {
    that.Has = append(that.Has, m.Has...)
    m.Has = m.Has[:0] // Nothing
}

type Items []Item
type Item int

