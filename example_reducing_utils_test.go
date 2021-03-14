package gofinancial_test

import (
	"fmt"

	gofinancial "github.com/razorpay/go-financial"
	"github.com/razorpay/go-financial/enums/paymentperiod"
	"github.com/shopspring/decimal"
)

// If you have a loan of 1,00,000 to be paid after 2 years, with 18% p.a. compounded annually, how much total payment will you have to do each month?
// This example generates the total monthly payment(principal plus interest) needed for a loan of 1,00,000 over 2 years with 18% rate of interest compounded monthly
func ExamplePmt_loan() {
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

// If an investment gives 6% rate of return compounded annually, how much amount should you invest each month to get 10,00,000 amount after 10 years?
func ExamplePmt_investment() {
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

// If an investment has a 6% p.a. rate of return, compounded annually, and you are investing ₹ 10,000 at the end of each year with initial investment of ₹ 10,000, how
// much amount will you get at the end of 10 years ?
func ExampleFv() {
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

// If you have a loan of 1,00,000 to be paid after 2 years, with 18% p.a. compounded annually, how much of the total payment done each month will be interest ?
func ExampleIPmt_loan() {
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

// If an investment gives 6% rate of return compounded annually, how much interest will you earn each year against your
// yearly payments(71574) to get 10,00,000 amount after 10 years
func ExampleIPmt_investment() {
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

// If you have a loan of 1,00,000 to be paid after 2 years, with 18% p.a. compounded annually, how much total payment done each month will be principal ?
func ExamplePPmt_loan() {
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

// If an investment has a 6% p.a. rate of return, compounded annually, and you wish to possess ₹ 1,49,716 at the end of 10 peroids while providing ₹ 10,000 per period,
// how much should you put as your initial deposit ?
func ExamplePv() {
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

// Given a discount rate of 8% per period and initial deposit of 40000 followed by withdrawls of 5000, 8000, 12000 and 30000.
// What is the net present value of the cash flow ?
func ExampleNpv() {
	rate := decimal.NewFromFloat(0.08)
	values := []decimal.Decimal{decimal.NewFromInt(-40000), decimal.NewFromInt(5000), decimal.NewFromInt(8000), decimal.NewFromInt(12000), decimal.NewFromInt(30000)}
	npv := gofinancial.Npv(rate, values)
	fmt.Printf("npv:%v", npv.Round(0))
	// Output:
	// npv:3065
}

// If a loan has a 6% annual interest, compounded monthly, and you only have $200/month to pay towards the loan,
// how long would it take to pay-off the loan of $5,000?
func ExampleNper() {
	rate := decimal.NewFromFloat(0.06 / 12)
	fv := decimal.Zero
	payment := decimal.NewFromInt(-200)
	pv := decimal.NewFromInt(5000)
	when := paymentperiod.ENDING

	nper, _ := gofinancial.Nper(rate, payment, pv, fv, when)
	fmt.Printf("nper:%v", nper.Ceil())
	// Output:
	// nper:27
}
