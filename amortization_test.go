package gofinancial

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"github.com/go-echarts/go-echarts/v2/charts"

	"github.com/razorpay/go-financial/enums/interesttype"
	"github.com/razorpay/go-financial/enums/paymentperiod"
	"github.com/smartystreets/assertions"

	"github.com/razorpay/go-financial/enums/frequency"
)

const (
	precision = 0.0000001
)

func Test_amortization_GenerateTable(t *testing.T) {
	type fields struct {
		Config *Config
	}

	tests := []struct {
		name    string
		fields  fields
		want    []Row
		wantErr bool
	}{
		{
			name:    "monthly table with rounding, reducing interest",
			fields:  fields{Config: getConfigDto(frequency.MONTHLY, true, interesttype.REDUCING, 1000000, 2400, 0)},
			want:    getRowsWithRounding(t),
			wantErr: false,
		},
		{
			name:    "monthly table without rounding, reducing interest",
			fields:  fields{Config: getConfigDto(frequency.MONTHLY, false, interesttype.REDUCING, 1000000, 2400, 0)},
			want:    getRowsWithoutRounding(t),
			wantErr: false,
		},
		{
			name: "daily table, flat interest, with rounding",
			fields: fields{
				Config: &Config{
					StartDate:      time.Date(2020, 4, 15, 0, 0, 0, 0, time.UTC),
					EndDate:        time.Date(2020, 5, 14, 0, 0, 0, 0, time.UTC),
					Frequency:      frequency.DAILY,
					AmountBorrowed: 1000000,
					InterestType:   interesttype.FLAT,
					Interest:       7300,
					PaymentPeriod:  paymentperiod.ENDING,
					EnableRounding: true,
					RoundingPlaces: 0,
				},
			},
			want:    getRowsFlatWithRounding(t),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := NewAmortization(tt.fields.Config)
			if err != nil {
				t.Errorf("NewAmortization() call failed. error = %v", err)
			}
			got, err := a.GenerateTable()
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateTable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Fatalf("length mismatch of rows generate, want=%v, got=%v", len(tt.want), len(got))
			}
			for idx := range got {
				if err := verifyRow(t, got[idx], tt.want[idx]); err != nil {
					t.Fatal(err)
				}
			}
			if err := principalCheck(t, got, -tt.fields.Config.AmountBorrowed); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func principalCheck(t *testing.T, rows []Row, actualPrincipal int64) error {
	expectedPrincipal := decimal.Zero
	dPrecision := decimal.NewFromFloat(precision)
	for _, row := range rows {
		expectedPrincipal = expectedPrincipal.Add(row.Principal)
	}
	if err := isAlmostEqual(expectedPrincipal, decimal.NewFromInt(actualPrincipal), dPrecision); err != nil {
		return fmt.Errorf("error:%v, principalCheck failed. expected:%v, got:%v", err.Error(), expectedPrincipal, actualPrincipal)
	}
	return nil
}

func verifyRow(t *testing.T, actual Row, expected Row) error {
	dPrecision := decimal.NewFromFloat(precision)
	if err := isAlmostEqual(actual.Principal, expected.Principal, dPrecision); err != nil {
		return fmt.Errorf("error:%v, principal values are not almost equal. expected:%v, got:%v", err.Error(), expected.Principal, actual.Principal)
	}
	if err := isAlmostEqual(actual.Interest, expected.Interest, dPrecision); err != nil {
		return fmt.Errorf("error:%v, interest values are not almost equal. expected:%v, got:%v", err.Error(), expected.Interest, actual.Interest)
	}
	if err := isAlmostEqual(actual.Payment, expected.Payment, dPrecision); err != nil {
		return fmt.Errorf("error:%v, payment values are not equal. expected:%v, got:%v", err.Error(), expected.Payment, actual.Payment)
	}
	if err := isAlmostEqual(actual.Principal.Add(actual.Interest), actual.Payment, dPrecision); err != nil {
		return fmt.Errorf("error:%v, the calculation is not correct. %v(Interest) + %v(Principal) != %v(Payment)", err.Error(), actual.Interest, actual.Principal, actual.Payment)
	}
	if !actual.StartDate.Equal(expected.StartDate) {
		return fmt.Errorf("start date value mismatch. Expected startDate:%v, endDate:%v, got startDate:%v endDate:%v", expected.StartDate, expected.EndDate, actual.StartDate, actual.EndDate)
	}
	if !actual.EndDate.Equal(expected.EndDate) {
		return fmt.Errorf("end date value mismatch. Expected startDate:%v, endDate:%v, got startDate:%v endDate:%v", expected.StartDate, expected.EndDate, actual.StartDate, actual.EndDate)
	}
	return nil
}

func isAlmostEqual(first decimal.Decimal, second decimal.Decimal, tolerance decimal.Decimal) error {
	diff := first.Sub(second)
	if diff.Abs().LessThanOrEqual(tolerance) {
		return nil
	} else {
		return fmt.Errorf("%s is not equal to %s with %s tolerance", first.String(), second.String(), tolerance.String())
	}
}

func getRowsWithRounding(t *testing.T) []Row {
	return []Row{
		{Period: 1, StartDate: timeParseUtil(t, "2020-04-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-05-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-52871), Interest: decimal.NewFromInt(-20000), Principal: decimal.NewFromInt(-32871)},
		{Period: 2, StartDate: timeParseUtil(t, "2020-05-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-06-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-52871), Interest: decimal.NewFromInt(-19342), Principal: decimal.NewFromInt(-33529)},
		{Period: 3, StartDate: timeParseUtil(t, "2020-06-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-07-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-52871), Interest: decimal.NewFromInt(-18672), Principal: decimal.NewFromInt(-34199)},
		{Period: 4, StartDate: timeParseUtil(t, "2020-07-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-08-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-52871), Interest: decimal.NewFromInt(-17988), Principal: decimal.NewFromInt(-34883)},
		{Period: 5, StartDate: timeParseUtil(t, "2020-08-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-09-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-52871), Interest: decimal.NewFromInt(-17290), Principal: decimal.NewFromInt(-35581)},
		{Period: 6, StartDate: timeParseUtil(t, "2020-09-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-10-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-52871), Interest: decimal.NewFromInt(-16579), Principal: decimal.NewFromInt(-36292)},
		{Period: 7, StartDate: timeParseUtil(t, "2020-10-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-11-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-52871), Interest: decimal.NewFromInt(-15853), Principal: decimal.NewFromInt(-37018)},
		{Period: 8, StartDate: timeParseUtil(t, "2020-11-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-12-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-52871), Interest: decimal.NewFromInt(-15112), Principal: decimal.NewFromInt(-37759)},
		{Period: 9, StartDate: timeParseUtil(t, "2020-12-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-01-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-52871), Interest: decimal.NewFromInt(-14357), Principal: decimal.NewFromInt(-38514)},
		{Period: 10, StartDate: timeParseUtil(t, "2021-01-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-02-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-52871), Interest: decimal.NewFromInt(-13587), Principal: decimal.NewFromInt(-39284)},
		{Period: 11, StartDate: timeParseUtil(t, "2021-02-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-03-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-52871), Interest: decimal.NewFromInt(-12801), Principal: decimal.NewFromInt(-40070)},
		{Period: 12, StartDate: timeParseUtil(t, "2021-03-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-04-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-52871), Interest: decimal.NewFromInt(-12000), Principal: decimal.NewFromInt(-40871)},
		{Period: 13, StartDate: timeParseUtil(t, "2021-04-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-05-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-52871), Interest: decimal.NewFromInt(-11183), Principal: decimal.NewFromInt(-41688)},
		{Period: 14, StartDate: timeParseUtil(t, "2021-05-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-06-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-52871), Interest: decimal.NewFromInt(-10349), Principal: decimal.NewFromInt(-42522)},
		{Period: 15, StartDate: timeParseUtil(t, "2021-06-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-07-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-52871), Interest: decimal.NewFromInt(-9498), Principal: decimal.NewFromInt(-43373)},
		{Period: 16, StartDate: timeParseUtil(t, "2021-07-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-08-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-52871), Interest: decimal.NewFromInt(-8631), Principal: decimal.NewFromInt(-44240)},
		{Period: 17, StartDate: timeParseUtil(t, "2021-08-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-09-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-52871), Interest: decimal.NewFromInt(-7746), Principal: decimal.NewFromInt(-45125)},
		{Period: 18, StartDate: timeParseUtil(t, "2021-09-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-10-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-52871), Interest: decimal.NewFromInt(-6844), Principal: decimal.NewFromInt(-46027)},
		{Period: 19, StartDate: timeParseUtil(t, "2021-10-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-11-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-52871), Interest: decimal.NewFromInt(-5923), Principal: decimal.NewFromInt(-46948)},
		{Period: 20, StartDate: timeParseUtil(t, "2021-11-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-12-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-52871), Interest: decimal.NewFromInt(-4984), Principal: decimal.NewFromInt(-47887)},
		{Period: 21, StartDate: timeParseUtil(t, "2021-12-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2022-01-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-52871), Interest: decimal.NewFromInt(-4026), Principal: decimal.NewFromInt(-48845)},
		{Period: 22, StartDate: timeParseUtil(t, "2022-01-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2022-02-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-52871), Interest: decimal.NewFromInt(-3049), Principal: decimal.NewFromInt(-49822)},
		{Period: 23, StartDate: timeParseUtil(t, "2022-02-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2022-03-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-52871), Interest: decimal.NewFromInt(-2053), Principal: decimal.NewFromInt(-50818)},
		{Period: 24, StartDate: timeParseUtil(t, "2022-03-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2022-04-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-52871), Interest: decimal.NewFromInt(-1037), Principal: decimal.NewFromInt(-51834)},
	}
}

func getRowsWithoutRounding(t *testing.T) []Row {
	return []Row{
		{Period: 1, StartDate: timeParseUtil(t, "2020-04-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-05-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromFloat(-52871.0972532498902312), Interest: decimal.NewFromFloat(-20000), Principal: decimal.NewFromFloat(-32871.0972532498902312)},
		{Period: 2, StartDate: timeParseUtil(t, "2020-05-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-06-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromFloat(-52871.0972532498902312), Interest: decimal.NewFromFloat(-19342.58), Principal: decimal.NewFromFloat(-33528.5172532498902312)},
		{Period: 3, StartDate: timeParseUtil(t, "2020-06-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-07-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromFloat(-52871.0972532498902312), Interest: decimal.NewFromFloat(-18672.0116), Principal: decimal.NewFromFloat(-34199.0856532498902312)},
		{Period: 4, StartDate: timeParseUtil(t, "2020-07-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-08-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromFloat(-52871.0972532498902312), Interest: decimal.NewFromFloat(-17988.031832), Principal: decimal.NewFromFloat(-34883.0654212498902312)},
		{Period: 5, StartDate: timeParseUtil(t, "2020-08-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-09-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromFloat(-52871.0972532498902312), Interest: decimal.NewFromFloat(-17290.37246864), Principal: decimal.NewFromFloat(-35580.7247846098902312)},
		{Period: 6, StartDate: timeParseUtil(t, "2020-09-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-10-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromFloat(-52871.0972532498902312), Interest: decimal.NewFromFloat(-16578.7599180128), Principal: decimal.NewFromFloat(-36292.3373352370902312)},
		{Period: 7, StartDate: timeParseUtil(t, "2020-10-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-11-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromFloat(-52871.0972532498902312), Interest: decimal.NewFromFloat(-15852.915116373056), Principal: decimal.NewFromFloat(-37018.1821368768342312)},
		{Period: 8, StartDate: timeParseUtil(t, "2020-11-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-12-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromFloat(-52871.0972532498902312), Interest: decimal.NewFromFloat(-15112.55341870051712), Principal: decimal.NewFromFloat(-37758.5438345493731112)},
		{Period: 9, StartDate: timeParseUtil(t, "2020-12-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-01-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromFloat(-52871.0972532498902312), Interest: decimal.NewFromFloat(-14357.3844870745274624), Principal: decimal.NewFromFloat(-38513.7127661753627688)},
		{Period: 10, StartDate: timeParseUtil(t, "2021-01-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-02-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromFloat(-52871.0972532498902312), Interest: decimal.NewFromFloat(-13587.112176816018011648), Principal: decimal.NewFromFloat(-39283.985076433872219552)},
		{Period: 11, StartDate: timeParseUtil(t, "2021-02-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-03-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromFloat(-52871.0972532498902312), Interest: decimal.NewFromFloat(-12801.43442035233838457), Principal: decimal.NewFromFloat(-40069.66283289755184663)},
		{Period: 12, StartDate: timeParseUtil(t, "2021-03-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-04-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromFloat(-52871.0972532498902312), Interest: decimal.NewFromFloat(-12000.043108759385183984), Principal: decimal.NewFromFloat(-40871.054144490505047216)},
		{Period: 13, StartDate: timeParseUtil(t, "2021-04-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-05-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromFloat(-52871.0972532498902312), Interest: decimal.NewFromFloat(-11182.62397093457285805592), Principal: decimal.NewFromFloat(-41688.47328231531737314408)},
		{Period: 14, StartDate: timeParseUtil(t, "2021-05-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-06-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromFloat(-52871.0972532498902312), Interest: decimal.NewFromFloat(-10348.8564503532643342505984), Principal: decimal.NewFromFloat(-42522.2408028966258969494016)},
		{Period: 15, StartDate: timeParseUtil(t, "2021-06-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-07-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromFloat(-52871.0972532498902312), Interest: decimal.NewFromFloat(-9498.413579360329540571690368), Principal: decimal.NewFromFloat(-43372.683673889560690628309632)},
		{Period: 16, StartDate: timeParseUtil(t, "2021-07-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-08-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromFloat(-52871.0972532498902312), Interest: decimal.NewFromFloat(-8630.96185094753618425412417536), Principal: decimal.NewFromFloat(-44240.13540230235404694587582464)},
		{Period: 17, StartDate: timeParseUtil(t, "2021-08-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-09-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromFloat(-52871.0972532498902312), Interest: decimal.NewFromFloat(-7746.1610879664868994798466588672), Principal: decimal.NewFromFloat(-45124.9361652834033317201533411328)},
		{Period: 18, StartDate: timeParseUtil(t, "2021-09-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-10-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromFloat(-52871.0972532498902312), Interest: decimal.NewFromFloat(-6843.664309725816620550723592044544), Principal: decimal.NewFromFloat(-46027.432943524073610649276407955456)},
		{Period: 19, StartDate: timeParseUtil(t, "2021-10-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-11-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromFloat(-52871.0972532498902312), Interest: decimal.NewFromFloat(-5923.11759592033298256949806388543488), Principal: decimal.NewFromFloat(-46947.97965732955724863050193611456512)},
		{Period: 20, StartDate: timeParseUtil(t, "2021-11-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-12-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromFloat(-52871.0972532498902312), Interest: decimal.NewFromFloat(-4984.1599478387396591396080251631435776), Principal: decimal.NewFromFloat(-47886.9373054111505720603919748368564224)},
		{Period: 21, StartDate: timeParseUtil(t, "2021-12-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2022-01-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromFloat(-52871.0972532498902312), Interest: decimal.NewFromFloat(-4026.423146795514424829480185666406449152), Principal: decimal.NewFromFloat(-48844.674106454375806370519814333593550848)},
		{Period: 22, StartDate: timeParseUtil(t, "2022-01-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2022-02-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromFloat(-52871.0972532498902312), Interest: decimal.NewFromFloat(-3049.53160973142475773770978937973457813504), Principal: decimal.NewFromFloat(-49821.56564351846547346229021062026542186496)},
		{Period: 23, StartDate: timeParseUtil(t, "2022-02-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2022-03-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromFloat(-52871.0972532498902312), Interest: decimal.NewFromFloat(-2053.1022419260531767582239851673292696977408), Principal: decimal.NewFromFloat(-50817.9950113238370544417760148326707303022592)},
		{Period: 24, StartDate: timeParseUtil(t, "2022-03-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2022-04-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromFloat(-52871.721800245812584823171708404468609676476416), Interest: decimal.NewFromFloat(-1036.744286764574301623748464870675855091695616), Principal: decimal.NewFromFloat(-51834.9775134812382831994232435337927545847808)},
	}
}

func getRowsFlatWithRounding(t *testing.T) []Row {
	return []Row{
		{Period: 1, StartDate: timeParseUtil(t, "2020-04-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-04-15 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-35333), Interest: decimal.NewFromInt(-2000), Principal: decimal.NewFromInt(-33333)},
		{Period: 2, StartDate: timeParseUtil(t, "2020-04-16 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-04-16 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-35333), Interest: decimal.NewFromInt(-2000), Principal: decimal.NewFromInt(-33333)},
		{Period: 3, StartDate: timeParseUtil(t, "2020-04-17 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-04-17 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-35333), Interest: decimal.NewFromInt(-2000), Principal: decimal.NewFromInt(-33333)},
		{Period: 4, StartDate: timeParseUtil(t, "2020-04-18 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-04-18 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-35333), Interest: decimal.NewFromInt(-2000), Principal: decimal.NewFromInt(-33333)},
		{Period: 5, StartDate: timeParseUtil(t, "2020-04-19 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-04-19 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-35333), Interest: decimal.NewFromInt(-2000), Principal: decimal.NewFromInt(-33333)},
		{Period: 6, StartDate: timeParseUtil(t, "2020-04-20 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-04-20 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-35333), Interest: decimal.NewFromInt(-2000), Principal: decimal.NewFromInt(-33333)},
		{Period: 7, StartDate: timeParseUtil(t, "2020-04-21 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-04-21 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-35333), Interest: decimal.NewFromInt(-2000), Principal: decimal.NewFromInt(-33333)},
		{Period: 8, StartDate: timeParseUtil(t, "2020-04-22 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-04-22 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-35333), Interest: decimal.NewFromInt(-2000), Principal: decimal.NewFromInt(-33333)},
		{Period: 9, StartDate: timeParseUtil(t, "2020-04-23 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-04-23 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-35333), Interest: decimal.NewFromInt(-2000), Principal: decimal.NewFromInt(-33333)},
		{Period: 10, StartDate: timeParseUtil(t, "2020-04-24 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-04-24 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-35333), Interest: decimal.NewFromInt(-2000), Principal: decimal.NewFromInt(-33333)},
		{Period: 11, StartDate: timeParseUtil(t, "2020-04-25 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-04-25 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-35333), Interest: decimal.NewFromInt(-2000), Principal: decimal.NewFromInt(-33333)},
		{Period: 12, StartDate: timeParseUtil(t, "2020-04-26 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-04-26 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-35333), Interest: decimal.NewFromInt(-2000), Principal: decimal.NewFromInt(-33333)},
		{Period: 13, StartDate: timeParseUtil(t, "2020-04-27 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-04-27 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-35333), Interest: decimal.NewFromInt(-2000), Principal: decimal.NewFromInt(-33333)},
		{Period: 14, StartDate: timeParseUtil(t, "2020-04-28 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-04-28 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-35333), Interest: decimal.NewFromInt(-2000), Principal: decimal.NewFromInt(-33333)},
		{Period: 15, StartDate: timeParseUtil(t, "2020-04-29 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-04-29 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-35333), Interest: decimal.NewFromInt(-2000), Principal: decimal.NewFromInt(-33333)},
		{Period: 16, StartDate: timeParseUtil(t, "2020-04-30 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-04-30 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-35333), Interest: decimal.NewFromInt(-2000), Principal: decimal.NewFromInt(-33333)},
		{Period: 17, StartDate: timeParseUtil(t, "2020-05-01 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-05-01 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-35333), Interest: decimal.NewFromInt(-2000), Principal: decimal.NewFromInt(-33333)},
		{Period: 18, StartDate: timeParseUtil(t, "2020-05-02 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-05-02 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-35333), Interest: decimal.NewFromInt(-2000), Principal: decimal.NewFromInt(-33333)},
		{Period: 19, StartDate: timeParseUtil(t, "2020-05-03 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-05-03 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-35333), Interest: decimal.NewFromInt(-2000), Principal: decimal.NewFromInt(-33333)},
		{Period: 20, StartDate: timeParseUtil(t, "2020-05-04 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-05-04 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-35333), Interest: decimal.NewFromInt(-2000), Principal: decimal.NewFromInt(-33333)},
		{Period: 21, StartDate: timeParseUtil(t, "2020-05-05 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-05-05 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-35333), Interest: decimal.NewFromInt(-2000), Principal: decimal.NewFromInt(-33333)},
		{Period: 22, StartDate: timeParseUtil(t, "2020-05-06 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-05-06 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-35333), Interest: decimal.NewFromInt(-2000), Principal: decimal.NewFromInt(-33333)},
		{Period: 23, StartDate: timeParseUtil(t, "2020-05-07 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-05-07 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-35333), Interest: decimal.NewFromInt(-2000), Principal: decimal.NewFromInt(-33333)},
		{Period: 24, StartDate: timeParseUtil(t, "2020-05-08 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-05-08 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-35333), Interest: decimal.NewFromInt(-2000), Principal: decimal.NewFromInt(-33333)},
		{Period: 25, StartDate: timeParseUtil(t, "2020-05-09 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-05-09 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-35333), Interest: decimal.NewFromInt(-2000), Principal: decimal.NewFromInt(-33333)},
		{Period: 26, StartDate: timeParseUtil(t, "2020-05-10 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-05-10 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-35333), Interest: decimal.NewFromInt(-2000), Principal: decimal.NewFromInt(-33333)},
		{Period: 27, StartDate: timeParseUtil(t, "2020-05-11 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-05-11 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-35333), Interest: decimal.NewFromInt(-2000), Principal: decimal.NewFromInt(-33333)},
		{Period: 28, StartDate: timeParseUtil(t, "2020-05-12 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-05-12 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-35333), Interest: decimal.NewFromInt(-2000), Principal: decimal.NewFromInt(-33333)},
		{Period: 29, StartDate: timeParseUtil(t, "2020-05-13 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-05-13 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-35333), Interest: decimal.NewFromInt(-2000), Principal: decimal.NewFromInt(-33333)},
		{Period: 30, StartDate: timeParseUtil(t, "2020-05-14 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-05-14 23:59:59 +0000 UTC"), Payment: decimal.NewFromInt(-35343), Interest: decimal.NewFromInt(-2000), Principal: decimal.NewFromInt(-33343)},
	}
}

func timeParseUtil(t *testing.T, input string) time.Time {
	resultTime, err := time.Parse("2006-01-02 15:04:05 -0700 MST", input)
	if err != nil {
		t.Fatalf("invalid time format, %v", err)
	}
	return resultTime
}

func getConfigDto(frequency frequency.Type, round bool, interestType interesttype.Type, amount int64, interest int64, places int32) *Config {
	return &Config{
		StartDate:      time.Date(2020, 4, 15, 0, 0, 0, 0, time.UTC),
		EndDate:        time.Date(2022, 4, 14, 0, 0, 0, 0, time.UTC),
		Frequency:      frequency,
		AmountBorrowed: amount,
		InterestType:   interestType,
		Interest:       interest,
		EnableRounding: round,
		RoundingPlaces: places,
	}
}

func TestPlotRows(t *testing.T) {
	type args struct {
		rows     []Row
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"plot for loan schedule",
			args{
				rows:     getRowsWithRounding(t),
				fileName: "loan-schedule",
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if err = PlotRows(tt.args.rows, tt.args.fileName); (err != nil) != tt.wantErr {
				t.Errorf("PlotRows() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRenderer(t *testing.T) {
	type args struct {
		bar *charts.Bar
	}
	rows := getRowsWithRounding(t)
	tests := []struct {
		name         string
		args         args
		writer       *bytes.Buffer
		stringWanted string
		wantErr      bool
		err          error
		errWriter    io.Writer
	}{
		{
			"success",
			args{
				bar: getStackedBarPlot(rows),
			},
			&bytes.Buffer{},
			getExpectedHtmlString(),
			false,
			nil,
			nil,
		},
		{
			"error while writing",
			args{
				bar: getStackedBarPlot(rows),
			},
			nil,
			getExpectedHtmlString(),
			true,
			errors.New("error writer"),
			&errorWriter{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.wantErr {
				err = renderer(tt.args.bar, tt.errWriter)
			} else {
				err = renderer(tt.args.bar, tt.writer)
				result := assertions.ShouldEqual(getHtmlWithoutUniqueId(tt.stringWanted), getHtmlWithoutUniqueId(tt.writer.String()))
				if result != "" {
					t.Errorf("Rendere() expected != actual. diff:%v", result)
				}
			}
			if err != nil && err.Error() != tt.err.Error() {
				t.Errorf("Renderer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

type errorWriter struct{}

func (er *errorWriter) Write(p []byte) (n int, err error) {
	return 1, errors.New("error writer")
}

func getHtmlWithoutUniqueId(input string) string {
	lines := strings.Split(input, "\n")
	var result []string
	for i := range lines {
		// skipping the unique id lines
		if i >= 11 && i <= 18 {
			continue
		}
		result = append(result, lines[i])
	}
	return strings.Join(result, "\n")
}

func getExpectedHtmlString() string {
	return `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Awesome go-echarts</title>
    <script src="https://go-echarts.github.io/go-echarts-assets/assets/echarts.min.js"></script>
</head>

<body>
<div class="container">
    <div class="item" id="upuUxUZbTXgQ" style="width:1200px;height:600px;"></div>
</div>

<script type="text/javascript">
    "use strict";
    let goecharts_upuUxUZbTXgQ = echarts.init(document.getElementById('upuUxUZbTXgQ'), "white");
    let option_upuUxUZbTXgQ = {"color":["#c23531","#2f4554","#61a0a8","#d48265","#91c7ae","#749f83","#ca8622","#bda29a","#6e7074","#546570"],"dataZoom":[{"type":"inside","end":50},{"type":"slider","end":50}],"legend":{"show":true},"series":[{"name":"Principal","type":"bar","stack":"stackA","waveAnimation":false,"data":[{"value":32871},{"value":33528},{"value":34199},{"value":34883},{"value":35581},{"value":36292},{"value":37018},{"value":37758},{"value":38514},{"value":39284},{"value":40070},{"value":40871},{"value":41688},{"value":42522},{"value":43373},{"value":44240},{"value":45125},{"value":46027},{"value":46948},{"value":47887},{"value":48845},{"value":49822},{"value":50818},{"value":51836}]},{"name":"Interest","type":"bar","stack":"stackA","waveAnimation":false,"data":[{"name":"20000","value":20000},{"name":"19343","value":19343},{"name":"18672","value":18672},{"name":"17988","value":17988},{"name":"17290","value":17290},{"name":"16579","value":16579},{"name":"15853","value":15853},{"name":"15113","value":15113},{"name":"14357","value":14357},{"name":"13587","value":13587},{"name":"12801","value":12801},{"name":"12000","value":12000},{"name":"11183","value":11183},{"name":"10349","value":10349},{"name":"9498","value":9498},{"name":"8631","value":8631},{"name":"7746","value":7746},{"name":"6844","value":6844},{"name":"5923","value":5923},{"name":"4984","value":4984},{"name":"4026","value":4026},{"name":"3049","value":3049},{"name":"2053","value":2053},{"name":"1037","value":1037}]},{"name":"Payment","type":"bar","stack":"stackA","waveAnimation":false,"data":[{"value":52871},{"value":52871},{"value":52871},{"value":52871},{"value":52871},{"value":52871},{"value":52871},{"value":52871},{"value":52871},{"value":52871},{"value":52871},{"value":52871},{"value":52871},{"value":52871},{"value":52871},{"value":52871},{"value":52871},{"value":52871},{"value":52871},{"value":52871},{"value":52871},{"value":52871},{"value":52871},{"value":52873}]}],"title":{"text":"Loan repayment schedule"},"toolbox":{"show":true},"tooltip":{"show":false},"xAxis":[{"data":["2020-05-14","2020-06-14","2020-07-14","2020-08-14","2020-09-14","2020-10-14","2020-11-14","2020-12-14","2021-01-14","2021-02-14","2021-03-14","2021-04-14","2021-05-14","2021-06-14","2021-07-14","2021-08-14","2021-09-14","2021-10-14","2021-11-14","2021-12-14","2022-01-14","2022-02-14","2022-03-14","2022-04-14"]}],"yAxis":[{}]};
    goecharts_upuUxUZbTXgQ.setOption(option_upuUxUZbTXgQ);
</script>

<style>
    .container {margin-top:30px; display: flex;justify-content: center;align-items: center;}
    .item {margin: auto;}
</style>
</body>
</html>
`
}
