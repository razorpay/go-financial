package gofinancial

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"path"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/razorpay/go-financial/enums/interesttype"
)

var writer io.Writer

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
	Payment   float64
	Interest  float64
	Principal float64
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
		if a.Config.Round {
			row.Payment = math.Round(payment)
			row.Principal = math.Round(principalPayment)
			row.Interest = math.Round(interestPayment)
		} else {
			row.Payment = payment
			row.Principal = principalPayment
			row.Interest = interestPayment
		}
		if i == a.Config.periods {
			PerformErrorCorrectionDueToRounding(&row, result, a.Config.AmountBorrowed, a.Config.Round)
		}
		if row.Payment != row.Principal+row.Interest {
			return nil, ErrPayment
		}
		result = append(result, row)

	}
	return result, nil
}

// PerformErrorCorrectionDueToRounding takes care of errors in principal and payment amount due to rounding.
// Only the final row is adjusted for rounding errors.
func PerformErrorCorrectionDueToRounding(finalRow *Row, rows []Row, principal int64, round bool) {
	principalCollected := finalRow.Principal
	for _, row := range rows {
		principalCollected += row.Principal
	}
	if round {
		diff := math.Abs(float64(principal)) - math.Abs(principalCollected)
		if diff > 0 {
			// subtracting diff coz payment, principal and interest are -ve.
			finalRow.Payment = math.Round(finalRow.Payment - diff)
			finalRow.Principal = math.Round(finalRow.Principal - diff)
		} else if diff < 0 {
			finalRow.Payment = math.Round(finalRow.Payment + diff)
			finalRow.Principal = math.Round(finalRow.Principal + diff)
		}
	} else {
		diff := math.Abs(float64(principal)) - math.Abs(principalCollected)
		if diff > 0 {
			finalRow.Payment = finalRow.Payment - diff
			finalRow.Principal = finalRow.Principal - diff
		} else {
			finalRow.Payment = finalRow.Payment + diff
			finalRow.Principal = finalRow.Principal + diff
		}
	}
}

// PrintRows outputs a formatted json for given rows as input.
func PrintRows(rows []Row) {
	bytes, _ := json.MarshalIndent(rows, "", "\t")
	fmt.Printf("%s", bytes)
}

// PlotRows uses the go-echarts package to generate an interactive plot from the Rows array.
func PlotRows(rows []Row, fileName string) error {
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
	for _, row := range rows {
		xAxis = append(xAxis, row.EndDate.Format("2006-01-02"))
		interestArr = append(interestArr, opts.BarData{Name: fmt.Sprintf("%v", -row.Interest), Value: -row.Interest})
		principalArr = append(principalArr, opts.BarData{Value: -row.Principal})
		paymentArr = append(paymentArr, opts.BarData{Value: -row.Payment})
	}
	// Put data into instance
	bar.SetXAxis(xAxis).
		AddSeries("Principal", principalArr).
		AddSeries("Interest", interestArr).
		AddSeries("Payment", paymentArr).SetSeriesOptions(
		charts.WithBarChartOpts(opts.BarChart{
			Stack: "stackA",
		}))
	completePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	filePath := path.Join(completePath, fileName)
	f, err := os.Create(fmt.Sprintf("%s.html", filePath))
	if err != nil {
		return err
	}
	defer func() {
		err = f.Close()
	}()
	if writer == nil {
		writer = f
	}
	if err := bar.Render(writer); err != nil {
		return err
	}
	return err
}
