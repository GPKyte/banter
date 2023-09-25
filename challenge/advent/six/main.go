package main

import(
    "sort"
    "log"
    "bufio"
    "io"

    "github.com/GPKyte/banter/challenge/advent/common"
)

const FirstMarkerSize int = 4 // Need four non repeating characters to signify a marker
const MessageMarkerSize int = 14

func main() {
    puzzle := common.OpenFirstArgAsFileReader()
    defer puzzle.Close()

    Solve(puzzle)
}

func Solve(this io.Reader) {
    ciphertext := LoadPuzzle(this)
    indices := IndexMarkers(FirstMarkerSize, ciphertext)
    startpoint := indices[0]
    moreIndices := IndexMarkers(MessageMarkerSize, ciphertext)
    nextPoint := moreIndices[0]

    log.Println(startpoint)
    log.Println(nextPoint)
}

func LoadPuzzle(from io.Reader) string {
    s := bufio.NewScanner(from)
    s.Split(bufio.ScanLines)
    s.Scan() // Just one line

    puzzle := s.Text()
    return puzzle
}

func IndexMarkers(ofSize int, from string) []int {
    indices := make([]int, 0, len(from)-ofSize)

    start := 0
    end := start + ofSize
    last := len(from)

    for end < last {
        end = start + ofSize
        if ThisIsAMarker(from[start:end]) {
            indices = append(indices, end)
        }
        start++
    }

    return indices
}

type ByteSet map[byte]bool
func (bs *ByteSet) Add(b byte) {(*bs)[b] = true}
func (bs *ByteSet) Keys() []byte {
    keys := make([]byte, 0)
    for k := range *bs {
        keys = append(keys, k)
    }
    sort.Slice(keys, func(i, j int) bool {return keys[i] < keys[j]})
    return keys
}
func ToSet(please []byte) []byte {
    setme := &ByteSet{}
    for _, b := range please {
        setme.Add(b)
    }
    return setme.Keys()
}

func ThisIsAMarker(orNot string) bool {
    return len(orNot) == len(ToSet([]byte(orNot)))
}
