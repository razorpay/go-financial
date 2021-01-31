package gofinancial

import "errors"

var (
	ErrPayment          = errors.New("payment not matching interest plus principal")
	ErrUnevenEndDate    = errors.New("uneven end date")
	ErrInvalidFrequency = errors.New("invalid frequency")
	ErrNotEqual         = errors.New("input values are not equal")
)
