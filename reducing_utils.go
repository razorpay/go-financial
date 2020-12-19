package calculator

import (
	"math"

	"github.com/razorpay/go-financial/enums/paymentperiod"
)

// TODO: add more documentation.
/*
Compute payment
 by solving
fv + pv*(1 + rate)**nper +pmt*(1 + rate*when)/rate*((1 + rate)**nper - 1) == 0
rate: rate of interest charged annually in decimal system
nper: number of periods it is compounded.
pv : present value
fv: future value
when: amount collection time
*/
func pmt(rate float64, nper int64, pv float64, fv float64, when paymentperiod.Type) float64 {
	factor := math.Pow(1.0+float64(rate), float64(nper))
	secondFactor := (factor - 1) * (1 + rate*when.Value()) / rate
	return -(pv*factor + fv) / secondFactor
}

// Compute interest payment
// TODO: add more test cases.
func ipmt(rate float64, per int64, nper int64, pv float64, fv float64, when paymentperiod.Type) float64 {
	totalPmt := pmt(rate, nper, pv, fv, when)
	ipmt := rbl(rate, per, totalPmt, pv, when) * rate
	if when == paymentperiod.BEGINNING {
		if per < 1 {
			return math.NaN()
		} else if per == 1 {
			return 0
		} else {
			// paying at the beginning, so discount it.
			return ipmt / (1 + rate)
		}
	} else {
		if per < 1 {
			return math.NaN()
		} else {
			return ipmt
		}
	}
}

// Compute principal payment
func ppmt(rate float64, per int64, nper int64, pv float64, fv float64, when paymentperiod.Type, round bool) float64 {
	total := pmt(rate, nper, pv, fv, when)
	ipmt := ipmt(rate, per, nper, pv, fv, when)
	if round {
		return math.Round(total) - math.Round(ipmt)
	} else {
		return total - ipmt
	}
}

// rbl: remaining balance
func rbl(rate float64, per int64, pmt float64, pv float64, when paymentperiod.Type) float64 {
	return fv(rate, (per - 1), pmt, pv, when)
}

/*
fv (future value) is computed by solving the equation::
 fv +
 pv*(1+rate)**nper +
 pmt*(1 + rate*when)/rate*((1 + rate)**nper - 1) == 0
*/
func fv(rate float64, nper int64, pmt float64, pv float64, when paymentperiod.Type) float64 {
	factor := math.Pow(1.0+float64(rate), float64(nper))
	secondFactor := (1 + rate*when.Value()) * (factor - 1) / rate
	return -pv*factor - pmt*secondFactor
}
