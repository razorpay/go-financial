package gofinancial

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/shopspring/decimal"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/razorpay/go-financial/enums/interesttype"
)

// Amortization struct holds the configuration and financial details.
type Amortization struct {
	Config    *Config
	Financial Financial
}

// NewAmortization return a new amortisation object with config and financial fields initialised.
func NewAmortization(c *Config) (*Amortization, error) {
	a := Amortization{Config: c}
	if err := a.Config.setPeriodsAndDates(); err != nil {
		return nil, err
	}
	a.Config.setTolerance()
	switch a.Config.InterestType {
	case interesttype.REDUCING:
		a.Financial = &Reducing{}
	case interesttype.FLAT:
		a.Financial = &Flat{}
	}
	return &a, nil
}

// Row represents a single row in an amortization schedule.
type Row struct {
	Period    int64
	StartDate time.Time
	EndDate   time.Time
	Payment   decimal.Decimal
	Interest  decimal.Decimal
	Principal decimal.Decimal
}

// GenerateTable constructs the amortization table based on the configuration.
func (a Amortization) GenerateTable() ([]Row, error) {
	var result []Row
	for i := int64(1); i <= a.Config.periods; i++ {
		var row Row
		row.Period = i
		row.StartDate = a.Config.startDates[i-1]
		row.EndDate = a.Config.endDates[i-1]

		payment := a.Financial.GetPayment(*a.Config)
		principalPayment := a.Financial.GetPrincipal(*a.Config, i)
		interestPayment := a.Financial.GetInterest(*a.Config, i)
		if a.Config.EnableRounding {
			row.Payment = payment.Round(a.Config.RoundingPlaces)
			row.Principal = principalPayment.Round(a.Config.RoundingPlaces)
			// to avoid rounding errors.
			row.Interest = row.Payment.Sub(row.Principal)
		} else {
			row.Payment = payment
			row.Principal = principalPayment
			row.Interest = interestPayment
		}
		if i == a.Config.periods {
			PerformErrorCorrectionDueToRounding(&row, result, a.Config.AmountBorrowed, a.Config.EnableRounding, a.Config.RoundingPlaces)
		}
		if err := sanityCheckUpdate(&row, a.Config.RoundingErrorTolerance); err != nil {
			return nil, err
		}
		result = append(result, row)

	}
	return result, nil
}

// PerformErrorCorrectionDueToRounding takes care of errors in principal and payment amount due to rounding.
// Only the final row is adjusted for rounding errors.
func PerformErrorCorrectionDueToRounding(finalRow *Row, rows []Row, principal int64, round bool, places int32) {
	principalCollected := finalRow.Principal
	dPrincipal := decimal.NewFromInt(principal)
	for _, row := range rows {
		principalCollected = principalCollected.Add(row.Principal)
	}
	diff := dPrincipal.Abs().Sub(principalCollected.Abs())
	if round {
		if diff.GreaterThan(decimal.Zero) {
			// subtracting diff coz payment, principal and interest are -ve.
			finalRow.Payment = finalRow.Payment.Sub(diff).Round(places)
			finalRow.Principal = finalRow.Principal.Sub(diff).Round(places)
		} else if diff.LessThan(decimal.Zero) {
			finalRow.Payment = finalRow.Payment.Add(diff).Round(places)
			finalRow.Principal = finalRow.Principal.Add(diff).Round(places)
		}
	} else {
		if diff.GreaterThan(decimal.Zero) {
			finalRow.Payment = finalRow.Payment.Sub(diff)
			finalRow.Principal = finalRow.Principal.Sub(diff)
		} else {
			finalRow.Payment = finalRow.Payment.Add(diff)
			finalRow.Principal = finalRow.Principal.Add(diff)
		}
	}
}

// sanityCheckUpdate verifies the equation,
// payment = principal + interest for every row.
// If there is a mismatch due to rounding error and it is withing the tolerance,
// the difference is adjusted against the interest.
func sanityCheckUpdate(row *Row, tolerance int64) error {
	if !row.Payment.Equal(row.Principal.Add(row.Interest)) {

		diff := row.Payment.Abs().Sub(row.Principal.Add(row.Interest).Abs())
		if diff.LessThanOrEqual(decimal.NewFromInt(tolerance)) {
			row.Interest = row.Interest.Sub(diff)
		} else {
			return ErrPayment
		}
	}
	return nil
}

// PrintRows outputs a formatted json for given rows as input.
func PrintRows(rows []Row) {
	bytes, _ := json.MarshalIndent(rows, "", "\t")
	fmt.Printf("%s", bytes)
}

// PlotRows uses the go-echarts package to generate an interactive plot from the Rows array.
func PlotRows(rows []Row, fileName string) (err error) {
	bar := getStackedBarPlot(rows)
	completePath, err := os.Getwd()
	if err != nil {
		return err
	}
	filePath := path.Join(completePath, fileName)
	f, err := os.Create(fmt.Sprintf("%s.html", filePath))
	if err != nil {
		return err
	}
	defer func() {
		// setting named err
		ferr := f.Close()
		if err == nil {
			err = ferr
		}
	}()
	return renderer(bar, f)
}

// getStackedBarPlot returns an instance for stacked bar plot.
func getStackedBarPlot(rows []Row) *charts.Bar {
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title: "Loan repayment schedule",
	},
	),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "1200px",
			Height: "600px",
		}),
		charts.WithToolboxOpts(opts.Toolbox{Show: true}),
		charts.WithLegendOpts(opts.Legend{Show: true}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:  "inside",
			Start: 0,
			End:   50,
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:  "slider",
			Start: 0,
			End:   50,
		}),
	)
	var xAxis []string
	var interestArr []opts.BarData
	var principalArr []opts.BarData
	var paymentArr []opts.BarData
	minusOne := decimal.NewFromInt(-1)
	for _, row := range rows {
		xAxis = append(xAxis, row.EndDate.Format("2006-01-02"))
		interestArr = append(interestArr, opts.BarData{Value: row.Interest.Mul(minusOne).String()})
		principalArr = append(principalArr, opts.BarData{Value: row.Principal.Mul(minusOne).String()})
		paymentArr = append(paymentArr, opts.BarData{Value: row.Payment.Mul(minusOne).String()})
	}
	// Put data into instance
	bar.SetXAxis(xAxis).
		AddSeries("Principal", principalArr).
		AddSeries("Interest", interestArr).
		AddSeries("Payment", paymentArr).SetSeriesOptions(
		charts.WithBarChartOpts(opts.BarChart{
			Stack: "stackA",
		}))
	return bar
}

// renderer renders the bar into the writer interface
func renderer(bar *charts.Bar, writer io.Writer) error {
	return bar.Render(writer)
}
