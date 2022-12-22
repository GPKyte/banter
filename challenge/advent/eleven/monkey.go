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

type Operation func()
type Choice func() *Monkey

var WorryLevel int = 1
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

