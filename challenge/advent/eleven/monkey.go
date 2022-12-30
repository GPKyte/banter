package monkey

import (
    "fmt"
    "io"
    "bufio"
    "errors"
    "strings"
    "strconv"
)

type Configuration struct {
    Verbose bool
    Debug   bool
    Version string
}

var Config = Configuration{
    Verbose:    false,
    Debug:      false,
    Version:    VersAlpha,
}

var VersAlpha = "0.1.0"

func New(from io.Reader) *Monkeys {
    group := make([]*Monkey, 0)
    all := Monkeys {}
    s := bufio.NewScanner(from)

    for ok := true; ok; ok = s.Scan() {
        nextMonkey := parseOneMonkey(s)
        nextMonkey.Group = &all
        group = append(group, nextMonkey)
    }

    all.Group = group
    return &all
}

// parseOneMonkey by scanning through 6 input lines
// Starting at the label line, ending at the ifFalse line
// Expecting empty new line in the next scan
func parseOneMonkey(via *bufio.Scanner) *Monkey {
    var ok bool = true
    ok = via.Scan()
    labelLine := via.Text()
    label := labelLine[len("Monkey "):strings.Index(labelLine, ":")]
    ok = via.Scan()
    itemsLine := via.Text()
    itemsLine = itemsLine[strings.Index(itemsLine, ":")+1:]
    itemsLine = strings.Trim(itemsLine, " ")
    itemsLine = strings.Replace(itemsLine, ",", "", -1)
    items := make(Items, 0)
    for _, itemstr := range strings.Split(itemsLine, " ") {
        itemnum, err := strconv.Atoi(itemstr)
        if err != nil {
            panic(err)
        }
        items = append(items, Item(itemnum))
    }
    ok = via.Scan()
    opText := via.Text()
    opText = opText[strings.Index(opText, ":")+1:]
    op := NewOperation(opText)
    ok = via.Scan()
    choiceText := via.Text()
    ok = via.Scan()
    choiceText += "\n" + via.Text()
    ok = via.Scan()
    choiceText += "\n" + via.Text()
    choice := NewChoice(parseChoice(choiceText))

    if !ok {
        return nil
    }
    return &Monkey{
        ID: label,
        Has: items,
        HandleItem: op,
        Decide: choice,
    }
}


type Monkey struct {
    ID string
    Has Items
    HandleItem Operation
    Decide Choice
    Group *Monkeys
}

type Monkeys struct {
    Group []*Monkey
}

func (m *Monkeys) Target(id string) (*Monkey, error) {
    for _, monk := range (*m).Group {
        if monk.ID == id {
            return monk, nil
        }
    }
    return nil, errors.New("No monkey was targetted by id "+id)
}

type Choice func(int) string
func NewChoice(test func(int) bool, monkeyIdIfTestPasses, monkeyIdIfTestFails string) Choice {
    return func(this int) string {
        if test(this) {
            return monkeyIdIfTestPasses
        } else {
            return monkeyIdIfTestFails
        }
    }
}
func parseChoice(from string) (func(int) bool, string, string) {
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
    test := func(worryLevel int) bool {
        divisor, err := strconv.Atoi(testparts[len(testparts)-1])
        if err != nil {
            panic(err)
        }
        return worryLevel % divisor == 0
    }

    second := lines[1]
    ifTrueLabel := second[strings.IndexAny(second, "0123456789"):]
    third := lines[2]
    ifFalseLabel := third[strings.IndexAny(third, "0123456789"):]

    return test, ifTrueLabel, ifFalseLabel
}

var WorryLevel int = 1
type Operation func(*Item)
func NewOperation(from string) Operation {
    // Operation: var = Expression
    // Expression: var | var Operator Expression
    // Two vars are recognized, old, new yet both pull from the same place
    // Operators include + and *, but could include /, -, or others
    var tokens []string = parseOperationTokens(from)


    operation := func(i *Item) {
        operator  := chooseOperator(tokens[3]) // + * / -
        operands, err := chooseOperands(i, tokens[2], tokens[4])
        if err != nil {
            panic(fmt.Errorf("%w Caused by line '%s'", err, from))
        }

        *i = Item(operator(*operands[0], *operands[1]) / 3)
    }
    return operation
}

func chooseOperands(itemRef *Item, from ...string) ([]*int, error) {
    ops := make([]*int, 0)

    for _, label := range from {
        var o *int

        switch label {
        case "new":
            o = (*int)(itemRef)
        case "old":
            o = (*int)(itemRef)
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

func (m *Monkey) TossTo(that *Monkey) {
    that.Has = append(that.Has, m.Has[0])
    m.Has = m.Has[1:]
}

func (m *Monkey) TossItems() {
    for len(m.Has) > 0 {
        item := &m.Has[0]
        m.HandleItem(item)
        receiver, err := m.Group.Target(m.Decide(int(*item)))
        if err != nil {
            panic(err)
        }
        m.TossTo(receiver)
    }
}

type Items []Item
type Item int

func (ms *Monkeys) GoARound() {
    // In the initial order:
    for _, m := range ms.Group {
        m.TossItems()
    }
}
