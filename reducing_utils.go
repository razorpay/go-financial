/*
This package is a go native port of the numpy-financial package with some additional helper
functions.

The functions in this package are a scalar version of their vectorised counterparts in
the numpy-financial(https://github.com/numpy/numpy-financial) library.

Currently, only some functions are ported, the remaining will be ported soon.
*/
package gofinancial

// TODO: update readme

import (
	"math"

	"github.com/razorpay/go-financial/enums/paymentperiod"
	"github.com/shopspring/decimal"
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
func Pmt(rate decimal.Decimal, nper int64, pv decimal.Decimal, fv decimal.Decimal, when paymentperiod.Type) decimal.Decimal {
	one := decimal.NewFromInt(1)
	minusOne := decimal.NewFromInt(-1)
	dWhen := decimal.NewFromInt(when.Value())
	dNper := decimal.NewFromInt(nper)
	dRateWithWhen := rate.Mul(dWhen)

	factor := one.Add(rate).Pow(dNper)
	secondFactor := factor.Sub(one).Mul(one.Add(dRateWithWhen)).Div(rate)

	return pv.Mul(factor).Add(fv).Div(secondFactor).Mul(minusOne)
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
func IPmt(rate decimal.Decimal, per int64, nper int64, pv decimal.Decimal, fv decimal.Decimal, when paymentperiod.Type) decimal.Decimal {
	// TODO: update nan and remove any rounding here.
	totalPmt := Pmt(rate, nper, pv, fv, when)
	one := decimal.NewFromInt(1)
	ipmt := rbl(rate, per, totalPmt, pv, when).Mul(rate)
	if when == paymentperiod.BEGINNING {
		if per < 1 {
			return decimal.NewFromFloat(math.NaN())
		} else if per == 1 {
			return decimal.NewFromFloat(math.NaN())
		} else {
			// paying at the beginning, so discount it.
			return ipmt.Div(one.Add(rate))
		}
	} else {
		if per < 1 {
			return decimal.NewFromFloat(math.NaN())
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
func PPmt(rate decimal.Decimal, per int64, nper int64, pv decimal.Decimal, fv decimal.Decimal, when paymentperiod.Type) decimal.Decimal {
	total := Pmt(rate, nper, pv, fv, when)
	ipmt := IPmt(rate, per, nper, pv, fv, when)
	return total.Sub(ipmt)
}

// Rbl computes remaining balance
func rbl(rate decimal.Decimal, per int64, pmt decimal.Decimal, pv decimal.Decimal, when paymentperiod.Type) decimal.Decimal {
	return Fv(rate, per-1, pmt, pv, when)
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
func Fv(rate decimal.Decimal, nper int64, pmt decimal.Decimal, pv decimal.Decimal, when paymentperiod.Type) decimal.Decimal {
	one := decimal.NewFromInt(1)
	minusOne := decimal.NewFromInt(-1)
	dWhen := decimal.NewFromInt(when.Value())
	dRateWithWhen := rate.Mul(dWhen)
	dNper := decimal.NewFromInt(nper)

	factor := one.Add(rate).Pow(dNper)
	secondFactor := factor.Sub(one).Mul(one.Add(dRateWithWhen)).Div(rate)

	return pv.Mul(factor).Add(pmt.Mul(secondFactor)).Mul(minusOne)
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
func Pv(rate decimal.Decimal, nper int64, pmt decimal.Decimal, fv decimal.Decimal, when paymentperiod.Type) decimal.Decimal {
	one := decimal.NewFromInt(1)
	minusOne := decimal.NewFromInt(-1)
	dWhen := decimal.NewFromInt(when.Value())
	dNper := decimal.NewFromInt(nper)
	dRateWithWhen := rate.Mul(dWhen)

	factor := one.Add(rate).Pow(dNper)
	secondFactor := factor.Sub(one).Mul(one.Add(dRateWithWhen)).Div(rate)

	return fv.Add(pmt.Mul(secondFactor)).Div(factor).Mul(minusOne)
}

/*
Npv computes the Net Present Value of a cash flow series

Params:

 rate	: a discount rate applied once per period
 values	: the value of the cash flow for that time period. Values provided here must be an array of float64

References:
	L. J. Gitman, “Principles of Managerial Finance, Brief,” 3rd ed., Addison-Wesley, 2003, pg. 346.

*/
func Npv(rate decimal.Decimal, values []decimal.Decimal) decimal.Decimal {
	internalNpv := decimal.NewFromFloat(0.0)
	currentRateT := decimal.NewFromFloat(1.0)
	one := decimal.NewFromInt(1)
	for _, currentVal := range values {
		internalNpv = internalNpv.Add(currentVal.Div(currentRateT))
		currentRateT = currentRateT.Mul(one.Add(rate))
	}
	return internalNpv
}
