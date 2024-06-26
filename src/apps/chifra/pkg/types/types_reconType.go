package types

import "strings"

type ReconType int

const (
	Invalid ReconType = 1 << iota
	Genesis
	DiffDiff
	SameSame
	DiffSame
	SameDiff
	ShouldNotHappen
	First
	Last
)

func (r ReconType) String() string {
	l := func(r ReconType, s string) string {
		if r&First != 0 {
			s = "first-" + s
			s = strings.Replace(s, "first-same", "first", 1)
		}
		if r&Last != 0 {
			s = s + "-last"
			s = strings.Replace(s, "same-last", "last", 1)
		}
		return s
	}

	rr := r &^ (First | Last)
	switch rr {
	case Genesis:
		return "genesis"
	case DiffDiff:
		return l(r, "diff-diff")
	case SameSame:
		return l(r, "same-same")
	case DiffSame:
		return l(r, "diff-same")
	case SameDiff:
		return l(r, "same-diff")
	default:
		return "invalid"
	}
}
