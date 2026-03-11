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
	v, ok := toValue[*t]
	if !ok {
		return 1 // safe default to avoid division by zero
	}
	return v
}
