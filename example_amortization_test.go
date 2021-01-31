package gofinancial_test

import (
	"time"

	gofinancial "github.com/razorpay/go-financial"
	"github.com/razorpay/go-financial/enums/frequency"
	"github.com/razorpay/go-financial/enums/interesttype"
	"github.com/razorpay/go-financial/enums/paymentperiod"
)

// This example generates amortization table for a loan of 20 lakhs over 15years at 12% per annum.
func ExampleAmortization_GenerateTable() {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		panic("location loading error")
	}
	currentDate := time.Date(2009, 11, 11, 0o4, 30, 0o0, 0, loc)
	config := gofinancial.Config{

		// start date is inclusive
		StartDate: currentDate,

		// end date is inclusive.
		EndDate:   currentDate.AddDate(15, 0, 0).AddDate(0, 0, -1),
		Frequency: frequency.ANNUALLY,

		// AmountBorrowed is in paisa
		AmountBorrowed: 200000000,

		// InterestType can be flat or reducing
		InterestType: interesttype.REDUCING,

		// interest is in basis points
		Interest: 1200,

		// amount is paid at the end of the period
		PaymentPeriod: paymentperiod.ENDING,

		// all values will be rounded
		EnableRounding: true,
	}
	amortization, err := gofinancial.NewAmortization(&config)
	if err != nil {
		panic(err)
	}

	rows, err := amortization.GenerateTable()
	if err != nil {
		panic(err)
	}
	gofinancial.PrintRows(rows)
	// Output:
	// [
	// 	{
	// 		"Period": 1,
	// 		"StartDate": "2009-11-11T04:30:00+05:30",
	// 		"EndDate": "2010-11-10T23:59:59+05:30",
	// 		"Payment": "-29364848",
	// 		"Interest": "-24000000",
	// 		"Principal": "-5364848"
	// 	},
	// 	{
	// 		"Period": 2,
	// 		"StartDate": "2010-11-11T00:00:00+05:30",
	// 		"EndDate": "2011-11-10T23:59:59+05:30",
	// 		"Payment": "-29364848",
	// 		"Interest": "-23356218",
	// 		"Principal": "-6008630"
	// 	},
	// 	{
	// 		"Period": 3,
	// 		"StartDate": "2011-11-11T00:00:00+05:30",
	// 		"EndDate": "2012-11-10T23:59:59+05:30",
	// 		"Payment": "-29364848",
	// 		"Interest": "-22635183",
	// 		"Principal": "-6729665"
	// 	},
	// 	{
	// 		"Period": 4,
	// 		"StartDate": "2012-11-11T00:00:00+05:30",
	// 		"EndDate": "2013-11-10T23:59:59+05:30",
	// 		"Payment": "-29364848",
	// 		"Interest": "-21827623",
	// 		"Principal": "-7537225"
	// 	},
	// 	{
	// 		"Period": 5,
	// 		"StartDate": "2013-11-11T00:00:00+05:30",
	// 		"EndDate": "2014-11-10T23:59:59+05:30",
	// 		"Payment": "-29364848",
	// 		"Interest": "-20923156",
	// 		"Principal": "-8441692"
	// 	},
	// 	{
	// 		"Period": 6,
	// 		"StartDate": "2014-11-11T00:00:00+05:30",
	// 		"EndDate": "2015-11-10T23:59:59+05:30",
	// 		"Payment": "-29364848",
	// 		"Interest": "-19910154",
	// 		"Principal": "-9454694"
	// 	},
	// 	{
	// 		"Period": 7,
	// 		"StartDate": "2015-11-11T00:00:00+05:30",
	// 		"EndDate": "2016-11-10T23:59:59+05:30",
	// 		"Payment": "-29364848",
	// 		"Interest": "-18775590",
	// 		"Principal": "-10589258"
	// 	},
	// 	{
	// 		"Period": 8,
	// 		"StartDate": "2016-11-11T00:00:00+05:30",
	// 		"EndDate": "2017-11-10T23:59:59+05:30",
	// 		"Payment": "-29364848",
	// 		"Interest": "-17504880",
	// 		"Principal": "-11859968"
	// 	},
	// 	{
	// 		"Period": 9,
	// 		"StartDate": "2017-11-11T00:00:00+05:30",
	// 		"EndDate": "2018-11-10T23:59:59+05:30",
	// 		"Payment": "-29364848",
	// 		"Interest": "-16081683",
	// 		"Principal": "-13283165"
	// 	},
	// 	{
	// 		"Period": 10,
	// 		"StartDate": "2018-11-11T00:00:00+05:30",
	// 		"EndDate": "2019-11-10T23:59:59+05:30",
	// 		"Payment": "-29364848",
	// 		"Interest": "-14487704",
	// 		"Principal": "-14877144"
	// 	},
	// 	{
	// 		"Period": 11,
	// 		"StartDate": "2019-11-11T00:00:00+05:30",
	// 		"EndDate": "2020-11-10T23:59:59+05:30",
	// 		"Payment": "-29364848",
	// 		"Interest": "-12702447",
	// 		"Principal": "-16662401"
	// 	},
	// 	{
	// 		"Period": 12,
	// 		"StartDate": "2020-11-11T00:00:00+05:30",
	// 		"EndDate": "2021-11-10T23:59:59+05:30",
	// 		"Payment": "-29364848",
	// 		"Interest": "-10702959",
	// 		"Principal": "-18661889"
	// 	},
	// 	{
	// 		"Period": 13,
	// 		"StartDate": "2021-11-11T00:00:00+05:30",
	// 		"EndDate": "2022-11-10T23:59:59+05:30",
	// 		"Payment": "-29364848",
	// 		"Interest": "-8463532",
	// 		"Principal": "-20901316"
	// 	},
	// 	{
	// 		"Period": 14,
	// 		"StartDate": "2022-11-11T00:00:00+05:30",
	// 		"EndDate": "2023-11-10T23:59:59+05:30",
	// 		"Payment": "-29364848",
	// 		"Interest": "-5955374",
	// 		"Principal": "-23409474"
	// 	},
	// 	{
	// 		"Period": 15,
	// 		"StartDate": "2023-11-11T00:00:00+05:30",
	// 		"EndDate": "2024-11-10T23:59:59+05:30",
	// 		"Payment": "-29364868",
	// 		"Interest": "-3146237",
	// 		"Principal": "-26218631"
	// 	}
	// ]
}
