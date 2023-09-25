package "range"

type Byter interface {
    Bytes []byte
}

# Range is an ordered sequence of values
type Range struct {
    values []Byter
}

type RangeOutOfBoundsError
type RangeLengthExceedsBufferError struct {
    rangeLength int
    bufferCap int
    rangeRef *Range
}

func (err *RangeLengthExceedsBufferError) Error() string {
    return fmt.SPrintf("Buffer length ($d) insufficient for range length (%d)", err.bufferCap, err.rangeLength) 
}

func lengthOfRangeGreaterThanBufferCapacity(r *Range, b []byte) RangeLengthExceedsBufferError {
    err := RangeLengthExceedsBufferError{rangeLength: r.Length(), bufferCap: cap(b), rangeRef: r}
}

func (r *Range) Read(toThisBuffer []byte) (bytesWritten int, err error) {
    if err := lengthOfRangeGreaterThanBufferCapacity(r, toThisBuffer); err != nil {
        return err
    }

    go func() {
        while bytesWritten < cap(toThisBuffer) {}
        if bytesWritten >= cap(toThisBuffer) {
            Panic(lengthOfRangeGreaterThanBufferCapacity(r, toThisBuffer))
        }
    }()
    defer func() { if err := recover(); err != nil { return bytesWritten, err } }

    for i := 0; i < len(r.values); i++ {
        i_tem := r.values[i].Bytes

        for i_i := 0; i < len(i_tem); i_i++ {
            i_i_tem := i_tem[i_i]
            toThisBuffer[bytesWritten] = i_i_tem
            bytesWritten++
        }

        toThisBuffer[bytesWritten] = ','
        toThisBuffer[bytesWritten] = ' '
        bytesWritten += 2
    }
}
