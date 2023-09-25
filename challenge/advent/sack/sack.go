package sack

import (
    "strings"
    "fmt"
    "sort"
)

func New(s string) *Sack {
    sack := newSack(s)
    return &sack
}

func newSack(s string) Sack {
    one, other := splitLineInHalf(s)
    return Sack{
        source: s,
        oneHalf: fillComponentWithItems(one),
        otherHalf: fillComponentWithItems(other),
    }
}

func fillComponentWithItems(s string) *Component {
    c := make(Component, 0, len(s))
    for _, si := range s {
        c = append(c, GetItem(string(si)))
    }
    return &c
}

type ItemType struct {
    Key string
    Priority int
}

var AllItemTypes map[string]*ItemType = func() map[string]*ItemType {
    types := make(map[string]*ItemType, 52)

    // Types are ordered by priority values 1-54, from lower to uppercase English letters
    alpha := "abcdefghijklmnopqrstuvwxyz"
    ALPHA := strings.ToUpper(alpha)
    for priority, key := range alpha+ALPHA {
        types[string(key)] = &ItemType{string(key), priority + 1}
    }
    return types
}()
var NotAnItemType = ItemType{Key: "*", Priority: -1}
var NoItem = Item{Type: &NotAnItemType}

type Item struct {
    Type *ItemType
}

func (i *Item) String() string {
    return i.Type.Key
}

func (i *Item) DifferentThan(this Item) bool {
    return i.Type.Key != this.Type.Key
}

func (i *Item) SameAs(this Item) bool {
    return !i.DifferentThan(this)
}

func (i Item) Priority() int {
    return i.Type.Priority
}

// Len, Less, and Swap implementation for using the sort package over Items
type Items []Item
func (i Items) Len() int {return len(i)}
func (i Items) Less(t, j int) bool {return i[t].Type.Priority < i[j].Type.Priority}
func (i Items) Swap(t, j int) {
    holdme := i[t]
    i[t] = i[j]
    i[j] = holdme
}

// GetItem regardless of whether the string given is recognized among all the written item types
func GetItem(s string) Item {
    i, ok := AllItemTypes[s]
    if !ok {
        i = &NotAnItemType
    }

    return Item{Type:i}
}

type Component Items
func (c *Component) String() string {
    pleaseSort := true
    return "(" + c.joinItemsWith(", ", pleaseSort) + ")"
}

func (c *Component) joinItemsWith(separator string, wantSorted bool) string {
    si := make([]string, 0, len(*c)) // Strings for items in component
    for _, ci := range *c {
        si = append(si, ci.String())
    }
    if wantSorted {
        sort.Strings(si)
    }

    return strings.Join(si, separator)
}

func (c *Component) Contains(i Item) bool {
    for _, ci := range *c {
        if ci.SameAs(i) {
            return true
        }
    }
    return false
}

type Sack struct {
    source string
    oneHalf *Component
    otherHalf *Component
}

func (s *Sack) OneHalf() *Component {return s.oneHalf}
func (s *Sack) OtherHalf() *Component {return s.otherHalf}
func (s *Sack) String() string {
    return fmt.Sprintf("A Sack containing:\n\t%v component, and also\n\t%v component.",
                        s.OneHalf(), s.OtherHalf())
}

// Outlier shared between both components in the sack.
// Return the Item found in both components.
func (s *Sack) Outlier() Item {
    var outlier Item
    pleaseDoNotSort := false
    inThis := s.OneHalf().joinItemsWith("", pleaseDoNotSort)
    ofThose := s.OtherHalf().joinItemsWith("", pleaseDoNotSort)

    if strings.ContainsAny(inThis, ofThose) {
        foundOutlierHere := strings.IndexAny(inThis, ofThose)
        outlier = GetItem(string(inThis[foundOutlierHere]))
    } else {
        outlier = NoItem
    }
    return outlier
}

// InBoth strings may be one or more of the same letters.
// Return a string made from one of every letter that appears in this string andThat string given.
func InBoth(this, andThat string) string {
    return InAll([]string{this,andThat})
}

// InAll given strings is a set of letters which appear at least once in each string.
// Return one of every letter (unique) found across all strings given as a new string.
func InAll(ofThese []string) string {
    crossCheck := map[string][]bool{}
    common := make([]string, 0)

    // Initialize a map of single letters referring to slices the len(ofThese)
    for _, key := range StringToSet(strings.Join(ofThese, "")) {
        crossCheck[key] = make([]bool, len(ofThese))
        // default value of bool is false, but if proven otherwise, here is where to set that value: the initialization step here.
    }
    
    // Creating a checklist of sorts to find which letters are in all strings
    for which, eachString := range ofThese {
        for _, letter := range eachString {
            crossCheck[string(letter)][which] = true
        }
    }

    // Using a key-index checklist table is a potential waste of space in space complexity
    // Yet it makes finding matches quite simple
    for letter, checkEachString := range crossCheck {
        if AllTrue(checkEachString) {
            common = append(common, letter)
        }
    }

    return strings.Join(common, "")
}

// StringToSet provides a list of string letters, perhaps better as runes
// Using sort and then reducing in linear fashion
func StringToSet(s string) []string {
    if len(s) == 0 {
        return make([]string, 0, 0)
    }

    // Choosing two separate slices, to avoid manipulating data in place
    list := make([]string, len(s))
    set := make([]string, 0, len(s))

    // Set up list from input
    for i, e := range s {
        list[i] = string(e)
    }

    // Set up Set from list
    set = append(set, list[0])
    for i := 1; i < len(list); i++ {
        if list[i] != set[len(set)-1] {
            set = append(set, list[i])
        }
    }
    return set
}

func AllTrue(perhaps []bool) bool {
    for _, everyOneTrue := range perhaps {
        if !everyOneTrue {
            return false
        }
    }
    return true
}

// ItemCommonToAllSacks will be found and returned.
// Used to find Elf team's badge Item.
func ItemCommonToAllSacks(s []*Sack) Item {
    sas := make([]string, 0, len(s)) // Sacks as Strings

    for _, eachSack := range s {
        sas = append(sas, eachSack.source)
    }
    return GetItem(InAll(sas))
}

