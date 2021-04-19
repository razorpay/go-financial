/*
This package is a go native port of the numpy-financial package with some additional helper
functions.

The functions in this package are a scalar version of their vectorised counterparts in
the numpy-financial(https://github.com/numpy/numpy-financial) library.

Currently, only some functions are ported, the remaining will be ported soon.
*/
package gofinancial

import (
	"fmt"
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
	var secondFactor decimal.Decimal
	if rate.Equal(decimal.Zero) {
		secondFactor = dNper
	} else {
		secondFactor = factor.Sub(one).Mul(one.Add(dRateWithWhen)).Div(rate)
	}
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
	totalPmt := Pmt(rate, nper, pv, fv, when)
	one := decimal.NewFromInt(1)
	ipmt := rbl(rate, per, totalPmt, pv, when).Mul(rate)
	if when == paymentperiod.BEGINNING {
		if per == 1 {
			return decimal.Zero
		} else {
			// paying at the beginning, so discount it.
			return ipmt.Div(one.Add(rate))
		}
	} else {
		return ipmt
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
func Nper(rate decimal.Decimal, pmt decimal.Decimal, pv decimal.Decimal, fv decimal.Decimal, when paymentperiod.Type) (result decimal.Decimal, err error) {
	defer func() {
		if r := recover(); r != nil {
			result = decimal.Zero
			err = fmt.Errorf("%w: %v", ErrOutOfBounds, r)
		}
	}()
	one := decimal.NewFromInt(1)
	minusOne := decimal.NewFromInt(-1)
	dWhen := decimal.NewFromInt(when.Value())
	dRateWithWhen := rate.Mul(dWhen)
	z := pmt.Mul(one.Add(dRateWithWhen)).Div(rate)
	numerator := minusOne.Mul(fv).Add(z).Div(pv.Add(z))
	denominator := one.Add(rate)
	floatNumerator, _ := numerator.BigFloat().Float64()
	floatDenominator, _ := denominator.BigFloat().Float64()
	logNumerator := math.Log(floatNumerator)
	logDenominator := math.Log(floatDenominator)
	dlogDenominator := decimal.NewFromFloat(logDenominator)
	result = decimal.NewFromFloat(logNumerator).Div(dlogDenominator)
	return result, nil
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

/*
This function computes the ratio that is used to find a single value that sets the non-liner equation to zero

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
func getRateRatio(pv, fv, pmt, curRate decimal.Decimal, nper int64, when paymentperiod.Type) decimal.Decimal {
	oneInDecimal := decimal.NewFromInt(1)
	whenInDecimal := decimal.NewFromInt(when.Value())
	nperInDecimal := decimal.NewFromInt(nper)

	f0 := curRate.Add(oneInDecimal).Pow(decimal.NewFromInt(nper)) // f0 := math.Pow((1 + curRate), float64(nper))
	f1 := f0.Div(curRate.Add(oneInDecimal))                       // f1 := f0 / (1 + curRate)

	yP0 := pv.Mul(f0)
	yP1 := pmt.Mul(oneInDecimal.Add(curRate.Mul(whenInDecimal))).Mul(f0.Sub(oneInDecimal)).Div(curRate)
	y := fv.Add(yP0).Add(yP1) // y := fv + pv*f0 + pmt*(1.0+curRate*when.Value())*(f0-1)/curRate

	derivativeP0 := nperInDecimal.Mul(f1).Mul(pv)
	derivativeP1 := pmt.Mul(whenInDecimal).Mul(f0.Sub(oneInDecimal)).Div(curRate)
	derivativeP2s0 := oneInDecimal.Add(curRate.Mul(whenInDecimal))
	derivativeP2s1 := ((curRate.Mul((nperInDecimal)).Mul(f1)).Sub(f0).Add(oneInDecimal)).Div(curRate.Mul(curRate))
	derivativeP2 := derivativeP2s0.Mul(derivativeP2s1)
	derivative := derivativeP0.Add(derivativeP1).Add(derivativeP2)
	// derivative := (float64(nper) * f1 * pv) + (pmt * ((when.Value() * (f0 - 1) / curRate) + ((1.0 + curRate*when.Value()) * ((curRate*float64(nper)*f1 - f0 + 1) / (curRate * curRate)))))

	return y.Div(derivative)
}

/*
Rate computes the Interest rate per period by running Newton Rapson to find an approximate value for:
 y = fv + pv*(1+rate)**nper + pmt*(1+rate*when)/rate*((1+rate)**nper-1)*(0 - y_previous) /(rate - rate_previous) = dy/drate {derivative of y w.r.t. rate}

Params:
 nper	: number of compounding periods
 pmt	: a (fixed) payment, paid either
	  at the beginning (when = 1) or the end (when = 0) of each period
 pv	: a present value
 fv	: a future value
 when 	: specification of whether payment is made
	  at the beginning (when = 1) or the end (when = 0) of each period
 maxIter 	: total number of iterations to perform calculation
 tolerance 	: accept result only if the difference in iteration values is less than the tolerance provided
 initialGuess 	: an initial point to start approximating from

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
func Rate(pv, fv, pmt decimal.Decimal, nper int64, when paymentperiod.Type, maxIter int64, tolerance, initialGuess decimal.Decimal) (decimal.Decimal, error) {
	var nextIterRate, currentIterRate decimal.Decimal = initialGuess, initialGuess

	for iter := int64(0); iter < maxIter; iter++ {
		currentIterRate = nextIterRate
		nextIterRate = currentIterRate.Sub(getRateRatio(pv, fv, pmt, currentIterRate, nper, when))
		//skip further loops if |nextIterRate-currentIterRate| < tolerance
		if nextIterRate.Sub(currentIterRate).Abs().LessThan(tolerance) {
			break
		}
	}

	if nextIterRate.Sub(currentIterRate).Abs().GreaterThanOrEqual(tolerance) {
		return decimal.Zero, ErrTolerence
	}
	return nextIterRate, nil
}
