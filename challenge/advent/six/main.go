package main

import(
    "log"
    "bufio"
    "io"

    "github.com/GPKyte/banter/challenge/advent/common"
)

const FirstMarkerSize int = 4 // Need four non repeating characters to signify a marker

func main() {
    source := common.OpenFirstArgAsFileReader()
    defer source.Close()

    ciphertext := LoadPuzzle(source)
    indices := IndexMarkers(FirstMarkerSize, ciphertext)
    startpoint := indices[0]

    log.Println(startpoint)
}

func LoadPuzzle(from io.Reader) string {
    r := bufio.NewReader(from)
    puzzle, err := r.ReadString('\n')
    if err != nil {
        panic(err)
    }
    return puzzle
}

func IndexMarkers(ofSize int, from string) []int {
    indices := make([]int, 0, len(from)-ofSize)

    for start, end := 0, ofSize; end <= len(from); start, end = start + 1, end + 1 {
        if ThisIsAMarker(ofSize, from[start:end]) {
            indices = append(indices, end)
        }
    }

    return indices
}

type ByteSet map[byte]bool
func (bs *ByteSet) Add(b byte) {(*bs)[b] = true}
func (bs *ByteSet) Keys() []byte {return []byte{}}
func ToSet(please []byte) []byte {
    setme := &ByteSet{}
    for _, b := range please {
        setme.Add(b)
    }
    return setme.Keys()
}

func ThisIsAMarker(ofSize int, orNot string) bool {return false}
