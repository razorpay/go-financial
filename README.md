
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
| nper                         |      |    Computes the number of periodic payments|
| pv                           |  ✅   |   Computes the present value of a payment|
| rate                         |      |    Computes the rate of interest per period|
| irr                          |      |    Computes the internal rate of return|
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
func Fv(rate float64, nper int64, pmt float64, pv float64, when paymentperiod.Type) float64  
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
[Run on go-playground](https://play.golang.org/p/l2-5aCHTBmH)


## Pv  

```go  
func Pv(rate float64, nper int64, pmt float64, fv float64, when paymentperiod.Type) float64 
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
[Run on go-playground](https://play.golang.org/p/xe1dXKxDEcY)


## Npv

```go  
func Npv(rate float64, values []float64) float64 
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
)

func main() {
	rate :=  decimal.NewFromFloat(0.281)
	values := []decimal.Decimal{decimal.NewFromInt(-100), decimal.NewFromInt(39), decimal.NewFromInt(59), decimal.NewFromInt(55), decimal.NewFromInt(20)}
	npv := gofinancial.Npv(rate, values)
	fmt.Printf("npv:%v", npv.Round(0))
	// Output:
	// npv: -0.008478591638455768
}
```
[Run on go-playground](https://play.golang.org/p/ma1it8-rFzn)


##  Pmt  
  
```go  
func Pmt(rate float64, nper int64, pv float64, fv float64, when paymentperiod.Type) float64  
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
	"math"
)

func main() {
	rate := 0.18 / 12
	nper := int64(12 * 2)
	pv := float64(100000)
	fv := float64(0)
	when := paymentperiod.ENDING
	pmt := gofinancial.Pmt(rate, nper, pv, fv, when)
	fmt.Printf("payment:%v", math.Round(pmt))
        // Output:
        // payment:-4992
}
```

[Run on go-playground](https://play.golang.org/p/Xa4XsqKi1te)

### Example(Pmt-Investment)

If an investment gives 6% rate of return compounded annually, how much amount should you invest each month to get 10,00,000 amount after 10 years?

```go
package main

import (
	"fmt"
	gofinancial "github.com/razorpay/go-financial"
	"github.com/razorpay/go-financial/enums/paymentperiod"
	"math"
)

func main() {
	rate := 0.06
	nper := int64(10)
	pv := float64(0)
	fv := float64(1000000)
	when := paymentperiod.BEGINNING
	pmt := gofinancial.Pmt(rate, nper, pv, fv, when)
	fmt.Printf("payment each year:%v", math.Round(pmt))
        // Output:
        // payment each year:-71574
}
```
[Run on go-playground](https://play.golang.org/p/iNeg4-1Y7fs)



##  IPmt  
  
```go  
func IPmt(rate float64, per int64, nper int64, pv float64, fv float64, when paymentperiod.Type) float64  
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
	"math"
)

func main() {
	rate := 0.18 / 12
	nper := int64(12 * 2)
	pv := float64(100000)
	fv := float64(0)
	when := paymentperiod.ENDING

	for i := int64(0); i < nper; i++ {
		pmt := gofinancial.IPmt(rate, i+1, nper, pv, fv, when)
		fmt.Printf("period:%d interest:%v\n", i+1, math.Round(pmt))
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
[Run on go-playground](https://play.golang.org/p/CinprXnphL4)

### Example(IPmt-Investment)

If an investment gives 6% rate of return compounded annually, how much interest will you earn each year against your yearly payments(71574) to get 10,00,000 amount after 10 years

```go
package main

import (
	"fmt"
	gofinancial "github.com/razorpay/go-financial"
	"github.com/razorpay/go-financial/enums/paymentperiod"
	"math"
)

func main() {
	rate := 0.06
	nper := int64(10)
	pv := float64(0)
	fv := float64(1000000)
	when := paymentperiod.BEGINNING

	for i := int64(1); i < nper+1; i++ {
		pmt := gofinancial.IPmt(rate, i+1, nper, pv, fv, when)
		fmt.Printf("period:%d interest earned:%v\n", i, math.Round(pmt))
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
[Run on go-playground](https://play.golang.org/p/yqmymAllA7k)

## PPmt  
  
```go  
func PPmt(rate float64, per int64, nper int64, pv float64, fv float64, when paymentperiod.Type, round bool) float64  
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
	"math"
)

func main() {
	rate := 0.18 / 12
	nper := int64(12 * 2)
	pv := float64(100000)
	fv := float64(0)
	when := paymentperiod.ENDING

	for i := int64(0); i < nper; i++ {
		pmt := gofinancial.PPmt(rate, i+1, nper, pv, fv, when, true)
		fmt.Printf("period:%d principal:%v\n", i+1, math.Round(pmt))
	}
	// Output:
	// period:1 principal:-3492
	// period:2 principal:-3544
	// period:3 principal:-3598
	// period:4 principal:-3652
	// period:5 principal:-3706
	// period:6 principal:-3762
	// period:7 principal:-3818
	// period:8 principal:-3876
	// period:9 principal:-3934
	// period:10 principal:-3993
	// period:11 principal:-4053
	// period:12 principal:-4113
	// period:13 principal:-4175
	// period:14 principal:-4238
	// period:15 principal:-4301
	// period:16 principal:-4366
	// period:17 principal:-4431
	// period:18 principal:-4498
	// period:19 principal:-4565
	// period:20 principal:-4634
	// period:21 principal:-4703
	// period:22 principal:-4774
	// period:23 principal:-4846
	// period:24 principal:-4918
}
```
[Run on go-playground](https://play.golang.org/p/bZuHcvmzUSK)
