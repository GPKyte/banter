package section

type Section int
func (s Section) PriorTo(this Section) bool {
    return s < this
}
func (s Section) SameAs(this Section) bool {
    return s == this
}
func (s Section) After(this Section) bool {
    return s > this
}

type SectionRange []Section
// Range is inclusive of both start and end, e.g. Range(3, 22) has 19 sections: 3, 4, 5, ..., 21, 22.
func Range(start, end int) *SectionRange {
    r := make(SectionRange, 1+end-start)

    for s := range r {
        r[s] = Section(s+start)
    }

    return &r
}
func (s *SectionRange) First() Section {return (*s)[0]}
func (s *SectionRange) Last() Section {return (*s)[s.Len()-1]}
func (s *SectionRange) empty() bool {return (*s).Len() == 0}
func (s *SectionRange) Len() int {return len(*s)}
func (s *SectionRange) Contains(r *SectionRange) bool {
    if s.empty() || r.empty() {
        return false
    }
    if s.Len() < r.Len() {
        // Shortcut
        return false
    }

    s1 := s.First()
    sn := s.Last()
    r1 := r.First()
    rn := r.Last()

    beginningIncludes := s1.SameAs(r1) || s1.PriorTo(r1)
    endingIncludes := sn.SameAs(rn) || sn.After(rn)

    whetherSContainsR := beginningIncludes && endingIncludes
    return whetherSContainsR
}

func EitherContainsTheOther(one, another *SectionRange) bool {
    return one.Contains(another) || another.Contains(one)
}

func (s *SectionRange) Overlap(r *SectionRange) *SectionRange {
    noOverlap := r.First().After(s.Last()) || s.First().After(r.Last())
    if noOverlap {
        return &SectionRange{}
    }
    if s.Contains(r) {
        return r
    }
    if r.Contains(s) {
        return s
    }
    
    start := max(s.First(), r.First())
    end := min(s.Last(), r.Last())

    return Range(int(start), int(end))
}

func DoAnyOverlap(s, r *SectionRange) bool {
    return s.Overlap(r).Len() > 0
}

func min(a, b Section) Section {
    var min Section
    if a <= b {
        min = a
    } else {
        min = b
    }
    return min
}
func max(a, b Section) Section {
    var max Section
    if a >= b {
        max = a
    } else {
        max = b
    }
    return max
}
