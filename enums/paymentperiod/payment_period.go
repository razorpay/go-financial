package paymentperiod

type Type uint8

const (
	BEGINNING Type = iota + 1
	ENDING
)

var value = map[Type]float64{
	BEGINNING: 1.0,
	ENDING:    0.0,
}

func (t Type) Value() float64 {
	return value[t]
}
