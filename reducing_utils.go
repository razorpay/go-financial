/*
This package is a go native port of the numpy-financial package with some additional helper
functions.

The functions in this package are a scalar version of their vectorised counterparts in
the numpy-financial(https://github.com/numpy/numpy-financial) library.

Currently, only some functions are ported, the remaining will be ported soon.
*/
package gofinancial

import (
	"math"

	"github.com/razorpay/go-financial/enums/paymentperiod"
)

/*
Pmt compute the fixed payment(principal + interest) against a loan amount ( fv = 0).

It can also be used to calculate the recurring payments needed to achieve a certain future value
given an initial deposit, a fixed periodically compounded interest rate, and the total number of periods.

It is obtained by solving the following equation:

	fv + pv*(1 + rate)**nper + pmt*(1 + rate*when)/rate*((1 + rate)**nper - 1) == 0

Params:
 rate	: rate of interest compounded once per period
 nper	: total number of periods to be compounded for
 pv	: present value (e.g., an amount borrowed)
 fv	: future value (e.g., 0)
 when	: specification of whether payment is made
	  at the beginning (when = 1) or the end
	  (when = 0) of each period

References:
	[WRW] Wheeler, D. A., E. Rathke, and R. Weir (Eds.) (2009, May).
	Open Document Format for Office Applications (OpenDocument)v1.2,
	Part 2: Recalculated Formula (OpenFormula) Format - Annotated Version,
	Pre-Draft 12. Organization for the Advancement of Structured Information
	Standards (OASIS). Billerica, MA, USA. [ODT Document].
	Available:
	http://www.oasis-open.org/committees/documents.php?wg_abbrev=office-formula
	OpenDocument-formula-20090508.odt
*/
func Pmt(rate float64, nper int64, pv float64, fv float64, when paymentperiod.Type) float64 {
	factor := math.Pow(1.0+float64(rate), float64(nper))
	var secondFactor float64
	if rate == 0 {
		secondFactor = float64(nper)
	} else {
		secondFactor = (factor - 1) * (1 + rate*when.Value()) / rate
	}
	return -(pv*factor + fv) / secondFactor
}

/*
IPmt computes interest payment for a loan under a given period.

Params:

 rate	: rate of interest compounded once per period
 per	: period under consideration
 nper	: total number of periods to be compounded for
 pv	: present value (e.g., an amount borrowed)
 fv	: future value (e.g., 0)
 when	: specification of whether payment is made
	  at the beginning (when = 1) or the end
	  (when = 0) of each period

References:
	[WRW] Wheeler, D. A., E. Rathke, and R. Weir (Eds.) (2009, May).
	Open Document Format for Office Applications (OpenDocument)v1.2,
	Part 2: Recalculated Formula (OpenFormula) Format - Annotated Version,
	Pre-Draft 12. Organization for the Advancement of Structured Information
	Standards (OASIS). Billerica, MA, USA. [ODT Document].
	Available:
	http://www.oasis-open.org/committees/documents.php?wg_abbrev=office-formula
	OpenDocument-formula-20090508.odt
*/
func IPmt(rate float64, per int64, nper int64, pv float64, fv float64, when paymentperiod.Type) float64 {
	totalPmt := Pmt(rate, nper, pv, fv, when)
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

/*
PPmt computes principal payment for a loan under a given period.

Params:

 rate	: rate of interest compounded once per period
 per	: period under consideration
 nper	: total number of periods to be compounded for
 pv	: present value (e.g., an amount borrowed)
 fv	: future value (e.g., 0)
 when	: specification of whether payment is made
	  at the beginning (when = 1) or the end
	  (when = 0) of each period

References:
	[WRW] Wheeler, D. A., E. Rathke, and R. Weir (Eds.) (2009, May).
	Open Document Format for Office Applications (OpenDocument)v1.2,
	Part 2: Recalculated Formula (OpenFormula) Format - Annotated Version,
	Pre-Draft 12. Organization for the Advancement of Structured Information
	Standards (OASIS). Billerica, MA, USA. [ODT Document].
	Available:
	http://www.oasis-open.org/committees/documents.php?wg_abbrev=office-formula
	OpenDocument-formula-20090508.odt
*/
func PPmt(rate float64, per int64, nper int64, pv float64, fv float64, when paymentperiod.Type, round bool) float64 {
	total := Pmt(rate, nper, pv, fv, when)
	ipmt := IPmt(rate, per, nper, pv, fv, when)
	if round {
		return math.Round(total) - math.Round(ipmt)
	} else {
		return total - ipmt
	}
}

// Rbl computes remaining balance
func rbl(rate float64, per int64, pmt float64, pv float64, when paymentperiod.Type) float64 {
	return Fv(rate, (per - 1), pmt, pv, when)
}

/*
Nper computes the number of periodic payments by solving the equation:

 fv +
 pv*(1 + rate)**nper +
 pmt*(1 + rate*when)/rate*((1 + rate)**nper - 1) = 0


Params:

 rate	: an interest rate compounded once per period
 pmt	: a (fixed) payment, paid either
	  at the beginning (when =  1) or the end (when = 0) of each period
 pv	: a present value
 when	: specification of whether payment is made
	  at the beginning (when = 1) or the end
	  (when = 0) of each period
 fv: a future value
 when	: specification of whether payment is made
	  at the beginning (when = 1) or the end
	  (when = 0) of each period

*/
func Nper(rate float64, pmt float64, pv float64, fv float64, when paymentperiod.Type) float64 {
	z := pmt * (1 + rate*when.Value()) / rate
	return math.Log((-fv+z)/(pv+z)) / math.Log(1+rate)
}

/*
Fv computes future value at the end of some periods(nper) by solving the following equation:

 fv +
 pv*(1+rate)**nper +
 pmt*(1 + rate*when)/rate*((1 + rate)**nper - 1) == 0


Params:

 pv	: a present value
 rate	: an interest rate compounded once per period
 nper	: total number of periods
 pmt	: a (fixed) payment, paid either
	  at the beginning (when =  1) or the end (when = 0) of each period
 when	: specification of whether payment is made
	  at the beginning (when = 1) or the end
	  (when = 0) of each period

References:
	[WRW] Wheeler, D. A., E. Rathke, and R. Weir (Eds.) (2009, May).
	Open Document Format for Office Applications (OpenDocument)v1.2,
	Part 2: Recalculated Formula (OpenFormula) Format - Annotated Version,
	Pre-Draft 12. Organization for the Advancement of Structured Information
	Standards (OASIS). Billerica, MA, USA. [ODT Document].
	Available:
	http://www.oasis-open.org/committees/documents.php?wg_abbrev=office-formula
	OpenDocument-formula-20090508.odt
*/
func Fv(rate float64, nper int64, pmt float64, pv float64, when paymentperiod.Type) float64 {
	factor := math.Pow(1.0+float64(rate), float64(nper))
	secondFactor := (1 + rate*when.Value()) * (factor - 1) / rate
	return -pv*factor - pmt*secondFactor
}

/*
Pv computes present value by solving the following equation:

 fv +
 pv*(1+rate)**nper +
 pmt*(1 + rate*when)/rate*((1 + rate)**nper - 1) == 0


Params:

 fv	: a future value
 rate	: an interest rate compounded once per period
 nper	: total number of periods
 pmt	: a (fixed) payment, paid either
	  at the beginning (when =  1) or the end (when = 0) of each period
 when	: specification of whether payment is made
	  at the beginning (when = 1) or the end
	  (when = 0) of each period

References:
	[WRW] Wheeler, D. A., E. Rathke, and R. Weir (Eds.) (2009, May).
	Open Document Format for Office Applications (OpenDocument)v1.2,
	Part 2: Recalculated Formula (OpenFormula) Format - Annotated Version,
	Pre-Draft 12. Organization for the Advancement of Structured Information
	Standards (OASIS). Billerica, MA, USA. [ODT Document].
	Available:
	http://www.oasis-open.org/committees/documents.php?wg_abbrev=office-formula
	OpenDocument-formula-20090508.odt
*/
func Pv(rate float64, nper int64, pmt float64, fv float64, when paymentperiod.Type) float64 {
	factor := math.Pow(1.0+float64(rate), float64(nper))
	secondFactor := (1 + rate*when.Value()) * (factor - 1) / rate
	return (-fv - pmt*secondFactor) / factor
}

/*
Npv computes the Net Present Value of a cash flow series

Params:

 rate	: a discount rate applied once per period
 values	: the value of the cash flow for that time period. Values provided here must be an array of float64

References:
	L. J. Gitman, “Principles of Managerial Finance, Brief,” 3rd ed., Addison-Wesley, 2003, pg. 346.

*/
func Npv(rate float64, values []float64) float64 {
	internalNpv := float64(0.0)
	currentRateT := float64(1.0)
	for _, current_val := range values {
		internalNpv += (current_val / currentRateT)
		currentRateT *= (1 + rate)
	}
	return internalNpv
}

/*
This function computs the ratio that is used to find a single value that sets the non-liner equation to zero

Params:
 nper 	: number of compounding periods
 pmt	: a (fixed) payment, paid either
	  	  at the beginning (when = 1) or the end (when = 0) of each period
 pv		: a present value
 fv		: a future value
 when 	: specification of whether payment is made
		  at the beginning (when = 1) or the end (when = 0) of each period
 curRate: the rate compounded once per period rate
*/
func getRateRatio(pv, fv, pmt, curRate float64, nper int64, when paymentperiod.Type) float64 {
	f0 := math.Pow((1 + curRate), float64(nper))
	f1 := f0 / (1 + curRate)
	y := fv + pv*f0 + pmt*(1.0+curRate*when.Value())*(f0-1)/curRate
	derivative := (float64(nper) * f1 * pv) + (pmt * ((when.Value() * (f0 - 1) / curRate) + ((1.0 + curRate*when.Value()) * ((curRate*float64(nper)*f1 - f0 + 1) / (curRate * curRate)))))

	return y / derivative
}

/*
Rate computes the Interest rate per period by running Newton Rapson to find an approximate value for:

y = fv + pv*(1+rate)**nper + pmt*(1+rate*when)/rate*((1+rate)**nper-1)
(0 - y_previous) /(rate - rate_previous) = dy/drate {derivative of y w.r.t. rate}


Params:
 nper 	: number of compounding periods
 pmt	: a (fixed) payment, paid either
	  	  at the beginning (when = 1) or the end (when = 0) of each period
 pv		: a present value
 fv		: a future value
 when 	: specification of whether payment is made
		  at the beginning (when = 1) or the end (when = 0) of each period
 params	: optional parameters for maxIter, tolerance, and initialGuess
References:
	Wheeler, D. A., E. Rathke, and R. Weir (Eds.) (2009, May). Open Document
    Format for Office Applications (OpenDocument)v1.2, Part 2: Recalculated
    Formula (OpenFormula) Format - Annotated Version, Pre-Draft 12.
    Organization for the Advancement of Structured Information Standards
    (OASIS). Billerica, MA, USA. [ODT Document]. Available:
    http://www.oasis-open.org/committees/documents.php?wg_abbrev=office-formula
    OpenDocument-formula-20090508.odt
*/

func Rate(pv, fv, pmt float64, nper int64, when paymentperiod.Type, params ...float64) (float64, bool) {
	initialGuess := 0.1
	tolerance := 1e-6
	maxIter := 100

	for index, value := range params {
		switch index {
		case 0:
			maxIter = int(value)
		case 1:
			initialGuess = value
		case 2:
			tolerance = value
		default:
			//no more values to be read
		}
	}

	var nextIterRate, currentIterRate float64 = initialGuess, initialGuess

	for iter := 0; iter < maxIter; iter++ {
		currentIterRate = nextIterRate
		nextIterRate = currentIterRate - getRateRatio(pv, fv, pmt, currentIterRate, nper, when)
	}

	if math.Abs(nextIterRate-currentIterRate) > tolerance {
		return nextIterRate, false
	}

	return nextIterRate, true
}
