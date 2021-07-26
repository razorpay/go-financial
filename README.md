
# go-financial  
  
This package is a go native port of the numpy-financial package with some additional helper  
functions.  
  
The functions in this package are a scalar version of their vectorised counterparts in   
the [numpy-financial](https://github.com/numpy/numpy-financial) library.  
  
[![unit-tests status](https://github.com/razorpay/go-financial/workflows/unit-tests/badge.svg?branch=master "unit-tests")]("https://github.com/razorpay/go-financial/workflows/unit-tests/badge.svg?branch=master")  [![Go Report Card](https://goreportcard.com/badge/github.com/razorpay/go-financial)](https://goreportcard.com/report/github.com/razorpay/go-financial)  [![codecov](https://codecov.io/gh/razorpay/go-financial/branch/master/graph/badge.svg)](https://codecov.io/gh/razorpay/go-financial)  [![GoDoc](https://godoc.org/github.com/razorpay/go-financial?status.svg)](https://godoc.org/github.com/razorpay/go-financial) [![Release](https://img.shields.io/github/release/razorpay/go-financial.svg?style=flat-square)](https://github.com/razorpay/go-financial/releases)  [![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)
  
Currently, only some functions are ported,  
which are as follows:  
  
  
| numpy-financial function     | go native function ported?   | info|
|:------------------------:    |:------------------:  | :------------------|
| fv                           |  ✅   |   Computes the  future value|
| ipmt                         |  ✅   |   Computes interest payment for a loan|
| pmt                          |  ✅   |   Computes the fixed periodic payment(principal + interest) made against a loan amount|
| ppmt                         |  ✅   |   Computes principal payment for a loan|
| nper                         |  ✅   |    Computes the number of periodic payments|
| pv                           |  ✅   |   Computes the present value of a payment|
| rate                         |  ✅   |    Computes the rate of interest per period|
| irr                          |  ✅   |    Computes the internal rate of return|
| npv                          |  ✅   |   Computes the net present value of a series of cash flow|
| mirr                         |      |    Computes the modified internal rate of return|
  
# Index  
While the numpy-financial package contains a set of elementary financial functions, this pkg also contains some helper functions on top of it. Their usage and description can be found below:  

  * [Amortisation(Generate Table)](#amortisation-generate-table-)
    + [Generated plot](#generated-plot)
  * [Fv(Future value)](#fv)
    + [Example(Fv)](#examplefv)
  * [Pv(Present value)](#pv)
	+ [Example(Pv)](#examplepv)
  * [Npv(Net present value)](#npv)
	+ [Example(Npv)](#examplenpv)
  * [Pmt(Payment)](#pmt)
    + [Example(Pmt-Loan)](#examplepmt-loan)
    + [Example(Pmt-Investment)](#examplepmt-investment)
  * [IPmt(Interest Payment)](#ipmt)
    + [Example(IPmt-Loan)](#exampleipmt-loan)
    + [Example(IPmt-Investment)](#exampleipmt-investment)
  * [PPmt(Principal Payment)](#ppmt)
    + [Example(PPmt-Loan)](#exampleppmt-loan)
  * [Nper(Number of payments)](#nper)
    + [Example(Nper-Loan)](#examplenper-loan)
  * [Rate(Interest Rate)](#rate)
	+ [Example(Rate-Investment)](#examplerate-investment)
  * [Irr(Internal Rate of Return)](#irr)
	+ [Example(Irr)](#exampleirr)
 
 Detailed documentation is available at [godoc](https://godoc.org/github.com/razorpay/go-financial).
## Amortisation(Generate Table)  
  
  
To generate the schedule for a loan of 20 lakhs over 15years at 12%p.a., you can do the following:  
  
```go  
package main

import (
	"time"

	"github.com/shopspring/decimal"

	financial "github.com/razorpay/go-financial"
	"github.com/razorpay/go-financial/enums/frequency"
	"github.com/razorpay/go-financial/enums/interesttype"
	"github.com/razorpay/go-financial/enums/paymentperiod"
)

func main() {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		panic("location loading error")
	}
	currentDate := time.Date(2009, 11, 11, 4, 30, 0, 0, loc)
	config := financial.Config{

		// start date is inclusive
		StartDate: currentDate,

		// end date is inclusive.
		EndDate:   currentDate.AddDate(15, 0, 0).AddDate(0, 0, -1),
		Frequency: frequency.ANNUALLY,

		// AmountBorrowed is in paisa
		AmountBorrowed: decimal.NewFromInt(200000000),

		// InterestType can be flat or reducing
		InterestType: interesttype.REDUCING,

		// interest is in basis points
		Interest: decimal.NewFromInt(1200),

		// amount is paid at the end of the period
		PaymentPeriod: paymentperiod.ENDING,

		// all values will be rounded
		EnableRounding: true,

		// it will be rounded to nearest int
		RoundingPlaces: 0,

		// no error is tolerated
		RoundingErrorTolerance: decimal.Zero,
	}
	amortization, err := financial.NewAmortization(&config)
	if err != nil {
		panic(err)
	}

	rows, err := amortization.GenerateTable()
	if err != nil {
		panic(err)
	}
	// Generates json output of the data
	financial.PrintRows(rows)
	// Generates a html file with plots of the given data.
	financial.PlotRows(rows, "20lakh-loan-repayment-schedule")
} 
  
```  
  
### Generated plot  
<img src="https://media1.giphy.com/media/G714Y7CoFKoA56fNXL/giphy.gif" width="100%">  
  
## Fv  
  
```go  
func Fv(rate decimal.Decimal, nper int64, pmt decimal.Decimal, pv decimal.Decimal, when paymentperiod.Type) decimal.Decimal 
```  
Params:  
```text
 pv   : a present value 
rate  : an interest rate compounded once per period 
nper  : total number of periods 
pmt   : a (fixed) payment, paid either at the beginning (when =  1)
        or the end (when = 0) of each period 
when  : specification of whether payment is made at the beginning (when = 1)
        or the end (when = 0) of each period  
```  

Fv computes future value at the end of some periods(nper).

### Example(Fv)

If an investment has a 6% p.a. rate of return, compounded annually, and you are investing ₹ 10,000 at the end of each year with initial investment of ₹ 10,000, how much amount will you get at the end of 10 years ?

```go
package main

import (
	"fmt"

	gofinancial "github.com/razorpay/go-financial"
	"github.com/razorpay/go-financial/enums/paymentperiod"
	"github.com/shopspring/decimal"
)

func main() {
	rate := decimal.NewFromFloat(0.06)
	nper := int64(10)
	payment := decimal.NewFromInt(-10000)
	pv := decimal.NewFromInt(-10000)
	when := paymentperiod.ENDING

	fv := gofinancial.Fv(rate, nper, payment, pv, when)
	fmt.Printf("fv:%v", fv.Round(0))
	// Output:
	// fv:149716
}
```
[Run on go-playground](https://play.golang.org/p/B08Fs4TlYuF)


## Pv  

```go  
func Pv(rate decimal.Decimal, nper int64, pmt decimal.Decimal, fv decimal.Decimal, when paymentperiod.Type) decimal.Decimal 
```  
Params:
```text
 fv	: a future value
 rate	: an interest rate compounded once per period
 nper	: total number of periods
 pmt	: a (fixed) payment, paid either
	  at the beginning (when =  1) or the end (when = 0) of each period
 when	: specification of whether payment is made
	  at the beginning (when = 1) or the end
	  (when = 0) of each period
```

Pv computes present value some periods(nper) before the future value.

### Example(Pv)

If an investment has a 6% p.a. rate of return, compounded annually, and you wish to possess ₹ 1,49,716 at the end of 10 peroids while providing ₹ 10,000 per period, how much should you put as your initial deposit ?

```go
package main

import (
	"fmt"

	gofinancial "github.com/razorpay/go-financial"
	"github.com/razorpay/go-financial/enums/paymentperiod"
	"github.com/shopspring/decimal"
)

func main() {
	rate := decimal.NewFromFloat(0.06)
	nper := int64(10)
	payment := decimal.NewFromInt(-10000)
	fv := decimal.NewFromInt(149716)
	when := paymentperiod.ENDING

	pv := gofinancial.Pv(rate, nper, payment, fv, when)
	fmt.Printf("pv:%v", pv.Round(0))
	// Output:
	// pv:-10000	
}
```
[Run on go-playground](https://play.golang.org/p/WW6cXQZa2_k)


## Npv

```go  
func Npv(rate decimal.Decimal, values []decimal.Decimal) decimal.Decimal 
```  
Params:
```text
 rate	: a discount rate compounded once per period
 values	: the value of the cash flow for that time period. Values provided here must be an array of float64
```

Npv computes net present value based on the discount rate and the values of cash flow over the course of the cash flow period

### Example(Npv)

Given a rate of 0.281 per period and initial deposit of 100 followed by withdrawls of 39, 59, 55, 20. What is the net present value of the cash flow ?

```go
package main

import (
	"fmt"
	gofinancial "github.com/razorpay/go-financial"	
	"github.com/shopspring/decimal"
)

func main() {
	rate :=  decimal.NewFromFloat(0.281)
	values := []decimal.Decimal{decimal.NewFromInt(-100), decimal.NewFromInt(39), decimal.NewFromInt(59), decimal.NewFromInt(55), decimal.NewFromInt(20)}
	npv := gofinancial.Npv(rate, values)
	fmt.Printf("npv:%v", npv)
	// Output:
	// npv: -0.008478591638426
}
```
[Run on go-playground](https://play.golang.org/p/4nzo1FOR3U0)


##  Pmt  
  
```go  
func Pmt(rate decimal.Decimal, nper int64, pv decimal.Decimal, fv decimal.Decimal, when paymentperiod.Type) decimal.Decimal
```  
Params:  
 ```text
rate  : rate of interest compounded once per period 
nper  : total number of periods to be compounded for 
pv    : present value (e.g., an amount borrowed) 
fv    : future value (e.g., 0) 
when  : specification of whether payment is made at the
         beginning (when = 1) or the end (when = 0) of each period  
``` 

Pmt compute the fixed payment(principal + interest) against a loan amount ( fv =  0).  
It can also be used to calculate the recurring payments needed to achieve a  certain future value given an initial deposit,
a fixed periodically compounded interest rate, and the total number of periods.  

### Example(Pmt-Loan)
If you have a loan of 1,00,000 to be paid after 2 years, with 18% p.a. compounded annually, how much total payment will you have to do each month? This example generates the total monthly payment(principal plus interest) needed for a loan of 1,00,000 over 2 years with 18% rate of interest compounded monthly

```go
package main

import (
	"fmt"
	gofinancial "github.com/razorpay/go-financial"
	"github.com/razorpay/go-financial/enums/paymentperiod"
	"github.com/shopspring/decimal"
)

func main() {
	rate := decimal.NewFromFloat(0.18 / 12)
	nper := int64(12 * 2)
	pv := decimal.NewFromInt(100000)
	fv := decimal.NewFromInt(0)
	when := paymentperiod.ENDING
	pmt := gofinancial.Pmt(rate, nper, pv, fv, when)
	fmt.Printf("payment:%v", pmt.Round(0))
        // Output:
        // payment:-4992
}
```

[Run on go-playground](https://play.golang.org/p/kFGoL9g4VM4)

### Example(Pmt-Investment)

If an investment gives 6% rate of return compounded annually, how much amount should you invest each month to get 10,00,000 amount after 10 years?

```go
package main

import (
	"fmt"
	gofinancial "github.com/razorpay/go-financial"
	"github.com/razorpay/go-financial/enums/paymentperiod"
	"github.com/shopspring/decimal"
)

func main() {
	rate := decimal.NewFromFloat(0.06)
	nper := int64(10)
	pv := decimal.NewFromInt(0)
	fv := decimal.NewFromInt(1000000)
	when := paymentperiod.BEGINNING
	pmt := gofinancial.Pmt(rate, nper, pv, fv, when)
	fmt.Printf("payment each year:%v", pmt.Round(0))
        // Output:
        // payment each year:-71574
}
```
[Run on go-playground](https://play.golang.org/p/aXC3vt-3UwS)



##  IPmt  
  
```go  
func IPmt(rate decimal.Decimal, per int64, nper int64, pv decimal.Decimal, fv decimal.Decimal, when paymentperiod.Type) decimal.Decimal   
```  
IPmt computes interest payment for a loan under a given period.  
  
Params:  
  ```text
rate  : rate of interest compounded once per period 
per   : period under consideration 
nper  : total number of periods to be compounded for 
pv    : present value (e.g., an amount borrowed) 
fv    : future value (e.g., 0) 
when  : specification of whether payment is made at the
          beginning (when = 1) or the end (when = 0) of each period  
```
 
### Example(IPmt-Loan)
If you have a loan of 1,00,000 to be paid after 2 years, with 18% p.a. compounded annually, how much of the total payment done each month will be interest ?

```go
package main

import (
	"fmt"
	gofinancial "github.com/razorpay/go-financial"
	"github.com/razorpay/go-financial/enums/paymentperiod"
	"github.com/shopspring/decimal"
)

func main() {
	rate := decimal.NewFromFloat(0.18 / 12)
	nper := int64(12 * 2)
	pv := decimal.NewFromInt(100000)
	fv := decimal.NewFromInt(0)
	when := paymentperiod.ENDING

	for i := int64(0); i < nper; i++ {
		ipmt := gofinancial.IPmt(rate, i+1, nper, pv, fv, when)
		fmt.Printf("period:%d interest:%v\n", i+1, ipmt.Round(0))
	}
	// Output:
	// period:1 interest:-1500
	// period:2 interest:-1448
	// period:3 interest:-1394
	// period:4 interest:-1340
	// period:5 interest:-1286
	// period:6 interest:-1230
	// period:7 interest:-1174
	// period:8 interest:-1116
	// period:9 interest:-1058
	// period:10 interest:-999
	// period:11 interest:-939
	// period:12 interest:-879
	// period:13 interest:-817
	// period:14 interest:-754
	// period:15 interest:-691
	// period:16 interest:-626
	// period:17 interest:-561
	// period:18 interest:-494
	// period:19 interest:-427
	// period:20 interest:-358
	// period:21 interest:-289
	// period:22 interest:-218
	// period:23 interest:-146
	// period:24 interest:-74
}
```
[Run on go-playground](https://play.golang.org/p/um60s2QZL_8)

### Example(IPmt-Investment)

If an investment gives 6% rate of return compounded annually, how much interest will you earn each year against your yearly payments(71574) to get 10,00,000 amount after 10 years

```go
package main

import (
	"fmt"
	gofinancial "github.com/razorpay/go-financial"
	"github.com/razorpay/go-financial/enums/paymentperiod"
	"github.com/shopspring/decimal"
)

func main() {
	rate := decimal.NewFromFloat(0.06)
	nper := int64(10)
	pv := decimal.NewFromInt(0)
	fv := decimal.NewFromInt(1000000)
	when := paymentperiod.BEGINNING

	for i := int64(1); i < nper+1; i++ {
		ipmt := gofinancial.IPmt(rate, i+1, nper, pv, fv, when)
		fmt.Printf("period:%d interest earned:%v\n", i, ipmt.Round(0))
	}
	// Output:
	// period:1 interest earned:4294
	// period:2 interest earned:8846
	// period:3 interest earned:13672
	// period:4 interest earned:18786
	// period:5 interest earned:24208
	// period:6 interest earned:29955
	// period:7 interest earned:36047
	// period:8 interest earned:42504
	// period:9 interest earned:49348
	// period:10 interest earned:56604
}
```
[Run on go-playground](https://play.golang.org/p/_NjNWuFulNp)

## PPmt  
  
```go  
func PPmt(rate decimal.Decimal, per int64, nper int64, pv decimal.Decimal, fv decimal.Decimal, when paymentperiod.Type) decimal.Decimal   
```  
PPmt computes principal payment for a loan under a given period.  
  
Params:  
```text
rate  : rate of interest compounded once per period 
per   : period under consideration 
nper  : total number of periods to be compounded for 
pv    : present value (e.g., an amount borrowed) 
fv    : future value (e.g., 0) 
when  : specification of whether payment is made at 
        the beginning (when = 1) or the end (when = 0) of each period  
```

### Example(PPmt-Loan)
If you have a loan of 1,00,000 to be paid after 2 years, with 18% p.a. compounded annually, how much total payment done each month will be principal ?
```go
package main

import (
	"fmt"
	gofinancial "github.com/razorpay/go-financial"
	"github.com/razorpay/go-financial/enums/paymentperiod"
	"github.com/shopspring/decimal"
)

func main() {
	rate := decimal.NewFromFloat(0.18 / 12)
	nper := int64(12 * 2)
	pv := decimal.NewFromInt(100000)
	fv := decimal.NewFromInt(0)
	when := paymentperiod.ENDING

	for i := int64(0); i < nper; i++ {
		ppmt := gofinancial.PPmt(rate, i+1, nper, pv, fv, when)
		fmt.Printf("period:%d principal:%v\n", i+1, ppmt.Round(0))
	}
	// Output:
	// period:1 principal:-3492
	// period:2 principal:-3545
	// period:3 principal:-3598
	// period:4 principal:-3652
	// period:5 principal:-3707
	// period:6 principal:-3762
	// period:7 principal:-3819
	// period:8 principal:-3876
	// period:9 principal:-3934
	// period:10 principal:-3993
	// period:11 principal:-4053
	// period:12 principal:-4114
	// period:13 principal:-4176
	// period:14 principal:-4238
	// period:15 principal:-4302
	// period:16 principal:-4366
	// period:17 principal:-4432
	// period:18 principal:-4498
	// period:19 principal:-4566
	// period:20 principal:-4634
	// period:21 principal:-4704
	// period:22 principal:-4774
	// period:23 principal:-4846
	// period:24 principal:-4919
}
```
[Run on go-playground](https://play.golang.org/p/s5nkkIeEj3x)




## Nper  
  
```go  
func Nper(rate decimal.Decimal, pmt decimal.Decimal, pv decimal.Decimal, fv decimal.Decimal, when paymentperiod.Type) (result decimal.Decimal, err error)  
```  
Params:  
```text
rate   : an interest rate compounded once per period
pmt    : a (fixed) payment, paid either at the beginning (when =  1)
         or the end (when = 0) of each period
pv     : a present value
fv     : a future value
when   : specification of whether payment is made at the beginning (when = 1)
         or the end (when = 0) of each period  
``` 

Nper computes the number of periodic payments.

### Example(Nper-Loan)

If a loan has a 6% annual interest, compounded monthly, and you only have \$200/month to pay towards the loan, how long would it take to pay-off the loan of \$5,000?

```go
package main

import (
	"fmt"

	gofinancial "github.com/razorpay/go-financial"
	"github.com/razorpay/go-financial/enums/paymentperiod"
	"github.com/shopspring/decimal"
)


func main() {
	rate := decimal.NewFromFloat(0.06 / 12)
	fv := decimal.NewFromFloat(0)
	payment := decimal.NewFromFloat(-200)
	pv := decimal.NewFromFloat(5000)
	when := paymentperiod.ENDING

	nper,err := gofinancial.Nper(rate, payment, pv, fv, when)
	if err != nil{
		fmt.Printf("error:%v\n", err)
	}
	fmt.Printf("nper:%v",nper.Ceil())
	// Output:
	// nper:27
}
```
[Run on go-playground](https://play.golang.org/p/hm77MTPBGYg)

## Rate  

```go  
func Rate(pv, fv, pmt decimal.Decimal, nper int64, when paymentperiod.Type, maxIter int64, tolerance, initialGuess decimal.Decimal) (decimal.Decimal, error) 
```  
Params:  
```text
pv     : a present value
fv     : a future value
pmt    : a (fixed) payment, paid either at the beginning (when =  1)
         or the end (when = 0) of each period
nper   : total number of periods to be compounded for
when   : specification of whether payment is made at the beginning (when = 1)
         or the end (when = 0) of each period
maxIter	: total number of iterations for which function should run
tolerance : tolerance threshold for acceptable result
initialGuess : an initial guess amount to start from
``` 

Returns:
```text
rate    : a value for the corresponding rate
error   : returns nil if rate difference is less than the threshold (returns an error conversely)
```

Rate computes the interest rate to ensure a balanced cashflow equation

### Example(Rate-Investment)

If an investment of $2000 is done and an amount of $100 is added at the start of each period, for what periodic interest rate would the invester be able to withdraw $3000 after the end of 4 periods ? (assuming 100 iterations, 1e-6 threshold and 0.1 as initial guessing point)

```go
package main

import (
	"fmt"
	gofinancial "github.com/razorpay/go-financial"
	"github.com/razorpay/go-financial/enums/paymentperiod"
	"github.com/shopspring/decimal"
)

func main() {
	fv := decimal.NewFromFloat(-3000)
	pmt := decimal.NewFromFloat(100)
	pv := decimal.NewFromFloat(2000)
	when := paymentperiod.BEGINNING
	nper := decimal.NewFromInt(4)
	maxIter := 100
	tolerance := decimal.NewFromFloat(1e-6)
	initialGuess := decimal.NewFromFloat(0.1),

	rate, err := gofinancial.Rate(pv, fv, pmt, nper, when, maxIter, tolerance, initialGuess)
	if err != nil {
		fmt.Printf(err)
	} else {
		fmt.Printf("rate: %v ", rate)
	}
	// Output:
	// rate: 0.06106257989825202
}
```
[Run on go-playground](https://play.golang.org/p/H2uybe1dbRj)

## Irr  

```go
func Irr(values []decimal.Decimal, maxIter int64, tolerance, prev_point, next_point decimal.Decimal) (decimal.Decimal, error)
```

Params:  
```text
values		: the value of the cash flow for that time period. Values provided here must be an array of float64
maxIter 	: total number of iterations for which the function should run
tolerance 	: accept result only if the difference in iteration values is less than the tolerance provided
prev_point	: an initial point to start approximating from
next_point	: next point to use for secant
``` 

Returns:
```text
rate    : a rate for the corresponding values
error   : returns nil if NPV is close to zero (returns an error conversely)
```

Irr computes the rate to ensure a net zero cashflow

### Example(Irr)

If an initial inflow of $123400 is done followed by successive outflows of $36200, $54800 and $48100; for what value of rate does the npv evalute to zero ? (assuming 100 iterations, 1e-6 threshold and 0.1 and 0.2 as initial guessing points)

```go
package main

import (
	"fmt"
	gofinancial "github.com/razorpay/go-financial"
	"github.com/razorpay/go-financial/enums/paymentperiod"
	"github.com/shopspring/decimal"
)

func main() {
	cashflow := []decimal.Decimal{decimal.NewFromFloat(-123400), decimal.NewFromFloat(36200), decimal.NewFromFloat(54800), decimal.NewFromFloat(48100)}
	maxIter := 100
	tolerance := decimal.NewFromFloat(1e-6)
	prev_point := decimal.NewFromFloat(0.1)
	next_point := decimal.NewFromFloat(0.2)

	irr, err := gofinancial.Irr(cashflow, maxIter, tolerance, prev_point, next_point)
	if err != nil {
		fmt.Printf(err)
	} else {
		fmt.Printf("irr: %v ", irr)
	}
	// Output:
	// irr: 0.05961637856732953787613704103503
}
```
[Run on go-playground](https://play.golang.org/p/H2uybe1dbRj)