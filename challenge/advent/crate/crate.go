package crate

type Crate string
const CrateSize int = len("[C]") // characters needed to label a crate in the puzzle's input e.g. '[C]'

func New(crate string) Crate {
    return Crate(crate)
}

var NoCrate = Crate("")

func (c Crate) DifferentThan(this Crate) bool {return true}
func (c Crate) SameAs(this Crate) bool {return true}
func (c Crate) ShortString() string {return string(c[1:2])}
