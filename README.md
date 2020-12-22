# go-financial

This package is a go native port of the numpy-financial package with some additional helper
functions.

The functions in this package are a scalar version of their vectorised counterparts in 
the numpy-financial library.

Currently, only some functions are ported,
which are as follows:


| numpy-financial function 	| go native function ported? 	|
|:------------------------:	|:------------------:	|
| fv                       	|  ✅ 	|
| pmt                      	|  ✅ 	|
| ipmt                     	|  ✅ 	|
| pmt                      	|  ✅ 	|
| ppmt                     	|  ✅ 	|
| nper                     	|  	|
| pv                       	|  	|
| rate                     	|  	|
| irr                      	|  	|
| npv                      	|  	|
| mirr                     	|  	|

# Usage

While the numpy-financial package contains a set of elementary financial functions, this pkg also contains some helper functions on top of it.

For example, To generate the schedule for a loan of 20 lakhs over 15years at 12%p.a., you can do the following:

```go
package main

import (
	"time"

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
	currentDate := time.Now().In(loc)

	config := financial.Config{
                // start date is inclusive
		StartDate:      currentDate,
                // end date is inclusive
		EndDate:        currentDate.AddDate(15, 0, 0).AddDate(0, 0, -1), 
		Frequency:      frequency.ANNUALLY,
                // AmountBorrowed is in paisa
		AmountBorrowed: 200000000,
                // InterestType can be flat or reducing.
		InterestType:   interesttype.REDUCING,
                // interest is in basis points
		Interest:       1200, 
		PaymentPeriod:  paymentperiod.ENDING,
		Round:          true,
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
