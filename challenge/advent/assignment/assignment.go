package assignment

import (
    "io"
    "fmt"
    "strings"
    "strconv"
    "bufio"

   "github.com/GPKyte/banter/challenge/advent/section"
   "github.com/GPKyte/banter/challenge/advent/elf"
)

type Assignment struct {
    s *section.SectionRange
    e *elf.Elf
}

type AssignmentPair struct {
    a, z *Assignment
}

func (ap *AssignmentPair) String() string {
    return fmt.Sprintf("%d-%d,%d-%d",
        ap.a.s.First(), ap.a.s.Last(), ap.z.s.First(), ap.z.s.Last())
}

func LoadFromPuzzleInput(from io.Reader) *[]*AssignmentPair {
    aps := make([]*AssignmentPair, 0)

    s := bufio.NewScanner(from)
    s.Split(bufio.ScanLines) // Set token scan behavior to read line by line
    for dataRemaining := s.Scan(); dataRemaining; dataRemaining = s.Scan() {
        line := s.Text()
        aps = append(aps, newPair(line))
    }

    return &aps
}
func newPair(s string) *AssignmentPair {
    first, second := splitAtThe(",", s)

    return &AssignmentPair{
        a: New(first),
        z: New(second),
    }
}
func splitAtThe(sep, s string) (string, string) {
    ss := strings.Split(s, sep)

    if len(ss) != 2 {
        panic("Could not properly split the input string: "+s)
    }

    return ss[0], ss[1]
}
func New(s string) *Assignment {
    one, two := splitAtThe("-", s)
    start, err := strconv.Atoi(one)
    if err != nil {
        panic("Could not parse the "+s+" string for new Assignment.")
    }
    end, err := strconv.Atoi(two)
    if err != nil {
        panic("Could not parse the "+s+" string for new Assignment.")
    }
    return &Assignment{
        s: section.Range(start, end),
        e: &elf.Elf{},
    }
}
func (ap *AssignmentPair) AnyOverlap() bool {
    return section.DoAnyOverlap(ap.a.s, ap.z.s)
}

func (ap *AssignmentPair) HasFullyRedundantSectionOverlap() bool {
    return section.EitherContainsTheOther(ap.a.s, ap.z.s)
}
