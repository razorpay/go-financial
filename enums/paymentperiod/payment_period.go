package paymentperiod

type Type uint8

const (
	BEGINNING Type = iota + 1
	ENDING
)

var value = map[Type]int64{
	BEGINNING: 1,
	ENDING:    0,
}

func (t Type) Value() int64 {
	return value[t]
}
