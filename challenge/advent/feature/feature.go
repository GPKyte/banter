package feature

type Status int

const (
    Enabled = iota
    Disabled
)

// Feature is useful for uncertain or evolving requirements while maintaining ease of porting to previous versions or use cases.
type Feature struct {
    Description string
    Status Status
    Name string
    Version Version
}

type Version string

func New(desc string, s Status) *Feature {
    return &Feature{
        Description: desc,
        Status: s,
    }
}

// Enabled features will use this method in conditional expression to use choose feature dependent code.
func (f *Feature) Enabled() bool {
    return f.Status == Enabled
}

// Disabled feature's status may suggest to use code which uses alternative strategy without dependence on this feature's related implementation.
func (f *Feature) Disabled() bool {
    return f.Status == Disabled
}

func (f *Feature) Enable() {
    f.Status = Enabled
}

func (f *Feature) Disable() {
    f.Status = Disabled
}

func (f *Feature) String() string {
    return f.Description + "("+string(f.Status)+")"
}

func (s Status) String() string {
    if s == Enabled {
        return "Enabled"
    } else {
        return "Disabled"
    }
}

