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
	secondFactor := (factor - 1) * (1 + rate*when.Value()) / rate
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

 rate	: a discount rate compounded once per period
 values	: the value of the cash flow for that time period. Values provided here must be an array of float64

References:
	L. J. Gitman, “Principles of Managerial Finance, Brief,” 3rd ed., Addison-Wesley, 2003, pg. 346.

*/
func Npv(rate float64, values []float64) float64 {
	internal_npv := float64(0.0)
	current_rate_t := float64(1.0)
	for _, current_val := range values {
		internal_npv += (current_val / current_rate_t)
		current_rate_t *= (1 + rate)
	}
	return internal_npv
}
