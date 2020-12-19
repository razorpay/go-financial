package interesttype

type Type uint8

const (
	FLAT Type = iota + 1
	REDUCING
)

var toString = map[Type]string{
	FLAT:     "flat",
	REDUCING: "reducing",
}

func (t Type) String() string {
	return toString[t]
}
