package frequency

type Type uint8

const (
	DAILY Type = iota + 1
	WEEKLY
	MONTHLY
	ANNUALLY
)

// TODO: check if this assumption is ok
var toValue = map[Type]int{
	DAILY:    365,
	WEEKLY:   52,
	MONTHLY:  12,
	ANNUALLY: 1,
}

func (t *Type) Value() int {
	return toValue[*t]
}
