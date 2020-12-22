package calculator

import (
	"fmt"
	"testing"
	"time"

	"github.com/razorpay/go-financial/enums/interesttype"
	"github.com/razorpay/go-financial/enums/paymentperiod"
	"github.com/smartystreets/assertions"

	"github.com/razorpay/go-financial/enums/frequency"
)

const (
	PRECISION = 0.0001
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
			fields:  fields{Config: getConfigDto(frequency.MONTHLY, true, interesttype.REDUCING, 1000000, 2400)},
			want:    getRowsWithRounding(t),
			wantErr: false,
		},
		{
			name:    "monthly table without rounding, reducing interest",
			fields:  fields{Config: getConfigDto(frequency.MONTHLY, false, interesttype.REDUCING, 1000000, 2400)},
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
					Round:          true,
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
	expectedPrincipal := 0.0
	for _, row := range rows {
		expectedPrincipal += row.Principal
	}
	if err := assertions.ShouldAlmostEqual(expectedPrincipal, actualPrincipal, PRECISION); err != "" {
		return fmt.Errorf("principalCheck failed. expected:%v, got:%v", expectedPrincipal, actualPrincipal)
	}
	return nil
}

func verifyRow(t *testing.T, actual Row, expected Row) error {
	if err := assertions.ShouldAlmostEqual(actual.Principal, expected.Principal, PRECISION); err != "" {
		return fmt.Errorf("principal values are not almost equal. expected:%v, got:%v", expected.Principal, actual.Principal)
	}
	if err := assertions.ShouldAlmostEqual(actual.Interest, expected.Interest, PRECISION); err != "" {
		return fmt.Errorf("interest values are not almost equal. expected:%v, got:%v", expected.Interest, actual.Interest)
	}
	if err := assertions.ShouldAlmostEqual(actual.Payment, expected.Payment, PRECISION); err != "" {
		return fmt.Errorf("payment values are not equal. expected:%v, got:%v", expected.Payment, actual.Payment)
	}
	if err := assertions.ShouldAlmostEqual(actual.Principal+actual.Interest, actual.Payment, PRECISION); err != "" {
		return fmt.Errorf("the calculation is not correct. %v(Interest) + %v(Principal) != %v(Payment)", actual.Interest, actual.Principal, actual.Payment)
	}
	if !actual.StartDate.Equal(expected.StartDate) {
		return fmt.Errorf("start date value mismatch. Expected startDate:%v, endDate:%v, got startDate:%v endDate:%v", expected.StartDate, expected.EndDate, actual.StartDate, actual.EndDate)
	}
	if !actual.EndDate.Equal(expected.EndDate) {
		return fmt.Errorf("end date value mismatch. Expected startDate:%v, endDate:%v, got startDate:%v endDate:%v", expected.StartDate, expected.EndDate, actual.StartDate, actual.EndDate)
	}
	return nil
}

func getRowsWithRounding(t *testing.T) []Row {
	return []Row{
		{Period: 1, StartDate: timeParseUtil(t, "2020-04-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-05-14 23:59:59 +0000 UTC"), Payment: -52871, Interest: -20000, Principal: -32871},
		{Period: 2, StartDate: timeParseUtil(t, "2020-05-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-06-14 23:59:59 +0000 UTC"), Payment: -52871, Interest: -19343, Principal: -33528},
		{Period: 3, StartDate: timeParseUtil(t, "2020-06-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-07-14 23:59:59 +0000 UTC"), Payment: -52871, Interest: -18672, Principal: -34199},
		{Period: 4, StartDate: timeParseUtil(t, "2020-07-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-08-14 23:59:59 +0000 UTC"), Payment: -52871, Interest: -17988, Principal: -34883},
		{Period: 5, StartDate: timeParseUtil(t, "2020-08-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-09-14 23:59:59 +0000 UTC"), Payment: -52871, Interest: -17290, Principal: -35581},
		{Period: 6, StartDate: timeParseUtil(t, "2020-09-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-10-14 23:59:59 +0000 UTC"), Payment: -52871, Interest: -16579, Principal: -36292},
		{Period: 7, StartDate: timeParseUtil(t, "2020-10-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-11-14 23:59:59 +0000 UTC"), Payment: -52871, Interest: -15853, Principal: -37018},
		{Period: 8, StartDate: timeParseUtil(t, "2020-11-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-12-14 23:59:59 +0000 UTC"), Payment: -52871, Interest: -15113, Principal: -37758},
		{Period: 9, StartDate: timeParseUtil(t, "2020-12-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-01-14 23:59:59 +0000 UTC"), Payment: -52871, Interest: -14357, Principal: -38514},
		{Period: 10, StartDate: timeParseUtil(t, "2021-01-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-02-14 23:59:59 +0000 UTC"), Payment: -52871, Interest: -13587, Principal: -39284},
		{Period: 11, StartDate: timeParseUtil(t, "2021-02-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-03-14 23:59:59 +0000 UTC"), Payment: -52871, Interest: -12801, Principal: -40070},
		{Period: 12, StartDate: timeParseUtil(t, "2021-03-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-04-14 23:59:59 +0000 UTC"), Payment: -52871, Interest: -12000, Principal: -40871},
		{Period: 13, StartDate: timeParseUtil(t, "2021-04-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-05-14 23:59:59 +0000 UTC"), Payment: -52871, Interest: -11183, Principal: -41688},
		{Period: 14, StartDate: timeParseUtil(t, "2021-05-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-06-14 23:59:59 +0000 UTC"), Payment: -52871, Interest: -10349, Principal: -42522},
		{Period: 15, StartDate: timeParseUtil(t, "2021-06-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-07-14 23:59:59 +0000 UTC"), Payment: -52871, Interest: -9498, Principal: -43373},
		{Period: 16, StartDate: timeParseUtil(t, "2021-07-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-08-14 23:59:59 +0000 UTC"), Payment: -52871, Interest: -8631, Principal: -44240},
		{Period: 17, StartDate: timeParseUtil(t, "2021-08-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-09-14 23:59:59 +0000 UTC"), Payment: -52871, Interest: -7746, Principal: -45125},
		{Period: 18, StartDate: timeParseUtil(t, "2021-09-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-10-14 23:59:59 +0000 UTC"), Payment: -52871, Interest: -6844, Principal: -46027},
		{Period: 19, StartDate: timeParseUtil(t, "2021-10-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-11-14 23:59:59 +0000 UTC"), Payment: -52871, Interest: -5923, Principal: -46948},
		{Period: 20, StartDate: timeParseUtil(t, "2021-11-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-12-14 23:59:59 +0000 UTC"), Payment: -52871, Interest: -4984, Principal: -47887},
		{Period: 21, StartDate: timeParseUtil(t, "2021-12-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2022-01-14 23:59:59 +0000 UTC"), Payment: -52871, Interest: -4026, Principal: -48845},
		{Period: 22, StartDate: timeParseUtil(t, "2022-01-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2022-02-14 23:59:59 +0000 UTC"), Payment: -52871, Interest: -3049, Principal: -49822},
		{Period: 23, StartDate: timeParseUtil(t, "2022-02-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2022-03-14 23:59:59 +0000 UTC"), Payment: -52871, Interest: -2053, Principal: -50818},
		{Period: 24, StartDate: timeParseUtil(t, "2022-03-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2022-04-14 23:59:59 +0000 UTC"), Payment: -52873, Interest: -1037, Principal: -51836},
	}
}

func getRowsWithoutRounding(t *testing.T) []Row {
	return []Row{
		{Period: 1, StartDate: timeParseUtil(t, "2020-04-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-05-14 23:59:59 +0000 UTC"), Payment: -52871.097253249915, Interest: -20000, Principal: -32871.097253249915},
		{Period: 2, StartDate: timeParseUtil(t, "2020-05-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-06-14 23:59:59 +0000 UTC"), Payment: -52871.097253249915, Interest: -19342.578054935002, Principal: -33528.51919831491},
		{Period: 3, StartDate: timeParseUtil(t, "2020-06-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-07-14 23:59:59 +0000 UTC"), Payment: -52871.097253249915, Interest: -18672.007670968705, Principal: -34199.08958228121},
		{Period: 4, StartDate: timeParseUtil(t, "2020-07-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-08-14 23:59:59 +0000 UTC"), Payment: -52871.097253249915, Interest: -17988.025879323082, Principal: -34883.07137392683},
		{Period: 5, StartDate: timeParseUtil(t, "2020-08-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-09-14 23:59:59 +0000 UTC"), Payment: -52871.097253249915, Interest: -17290.364451844544, Principal: -35580.73280140537},
		{Period: 6, StartDate: timeParseUtil(t, "2020-09-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-10-14 23:59:59 +0000 UTC"), Payment: -52871.097253249915, Interest: -16578.749795816435, Principal: -36292.34745743348},
		{Period: 7, StartDate: timeParseUtil(t, "2020-10-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-11-14 23:59:59 +0000 UTC"), Payment: -52871.097253249915, Interest: -15852.902846667766, Principal: -37018.194406582144},
		{Period: 8, StartDate: timeParseUtil(t, "2020-11-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-12-14 23:59:59 +0000 UTC"), Payment: -52871.097253249915, Interest: -15112.538958536128, Principal: -37758.55829471379},
		{Period: 9, StartDate: timeParseUtil(t, "2020-12-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-01-14 23:59:59 +0000 UTC"), Payment: -52871.097253249915, Interest: -14357.36779264185, Principal: -38513.72946060807},
		{Period: 10, StartDate: timeParseUtil(t, "2021-01-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-02-14 23:59:59 +0000 UTC"), Payment: -52871.097253249915, Interest: -13587.09320342969, Principal: -39284.00404982023},
		{Period: 11, StartDate: timeParseUtil(t, "2021-02-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-03-14 23:59:59 +0000 UTC"), Payment: -52871.097253249915, Interest: -12801.413122433281, Principal: -40069.68413081663},
		{Period: 12, StartDate: timeParseUtil(t, "2021-03-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-04-14 23:59:59 +0000 UTC"), Payment: -52871.097253249915, Interest: -12000.019439816957, Principal: -40871.07781343296},
		{Period: 13, StartDate: timeParseUtil(t, "2021-04-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-05-14 23:59:59 +0000 UTC"), Payment: -52871.097253249915, Interest: -11182.597883548293, Principal: -41688.49936970162},
		{Period: 14, StartDate: timeParseUtil(t, "2021-05-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-06-14 23:59:59 +0000 UTC"), Payment: -52871.097253249915, Interest: -10348.82789615426, Principal: -42522.269357095654},
		{Period: 15, StartDate: timeParseUtil(t, "2021-06-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-07-14 23:59:59 +0000 UTC"), Payment: -52871.097253249915, Interest: -9498.382509012346, Principal: -43372.71474423757},
		{Period: 16, StartDate: timeParseUtil(t, "2021-07-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-08-14 23:59:59 +0000 UTC"), Payment: -52871.097253249915, Interest: -8630.928214127603, Principal: -44240.16903912231},
		{Period: 17, StartDate: timeParseUtil(t, "2021-08-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-09-14 23:59:59 +0000 UTC"), Payment: -52871.097253249915, Interest: -7746.124833345148, Principal: -45124.97241990476},
		{Period: 18, StartDate: timeParseUtil(t, "2021-09-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-10-14 23:59:59 +0000 UTC"), Payment: -52871.097253249915, Interest: -6843.625384947052, Principal: -46027.47186830286},
		{Period: 19, StartDate: timeParseUtil(t, "2021-10-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-11-14 23:59:59 +0000 UTC"), Payment: -52871.097253249915, Interest: -5923.075947581004, Principal: -46948.02130566891},
		{Period: 20, StartDate: timeParseUtil(t, "2021-11-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2021-12-14 23:59:59 +0000 UTC"), Payment: -52871.097253249915, Interest: -4984.115521467622, Principal: -47886.98173178229},
		{Period: 21, StartDate: timeParseUtil(t, "2021-12-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2022-01-14 23:59:59 +0000 UTC"), Payment: -52871.097253249915, Interest: -4026.375886831973, Principal: -48844.72136641794},
		{Period: 22, StartDate: timeParseUtil(t, "2022-01-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2022-02-14 23:59:59 +0000 UTC"), Payment: -52871.097253249915, Interest: -3049.4814595036164, Principal: -49821.6157937463},
		{Period: 23, StartDate: timeParseUtil(t, "2022-02-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2022-03-14 23:59:59 +0000 UTC"), Payment: -52871.097253249915, Interest: -2053.049143628683, Principal: -50818.04810962123},
		{Period: 24, StartDate: timeParseUtil(t, "2022-03-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2022-04-14 23:59:59 +0000 UTC"), Payment: -52871.097253249915, Interest: -1036.6881814362714, Principal: -51834.40907181364},
	}
}

func getRowsFlatWithRounding(t *testing.T) []Row {
	return []Row{
		{Period: 1, StartDate: timeParseUtil(t, "2020-04-15 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-04-15 23:59:59 +0000 UTC"), Payment: -35333, Interest: -2000, Principal: -33333},
		{Period: 2, StartDate: timeParseUtil(t, "2020-04-16 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-04-16 23:59:59 +0000 UTC"), Payment: -35333, Interest: -2000, Principal: -33333},
		{Period: 3, StartDate: timeParseUtil(t, "2020-04-17 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-04-17 23:59:59 +0000 UTC"), Payment: -35333, Interest: -2000, Principal: -33333},
		{Period: 4, StartDate: timeParseUtil(t, "2020-04-18 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-04-18 23:59:59 +0000 UTC"), Payment: -35333, Interest: -2000, Principal: -33333},
		{Period: 5, StartDate: timeParseUtil(t, "2020-04-19 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-04-19 23:59:59 +0000 UTC"), Payment: -35333, Interest: -2000, Principal: -33333},
		{Period: 6, StartDate: timeParseUtil(t, "2020-04-20 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-04-20 23:59:59 +0000 UTC"), Payment: -35333, Interest: -2000, Principal: -33333},
		{Period: 7, StartDate: timeParseUtil(t, "2020-04-21 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-04-21 23:59:59 +0000 UTC"), Payment: -35333, Interest: -2000, Principal: -33333},
		{Period: 8, StartDate: timeParseUtil(t, "2020-04-22 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-04-22 23:59:59 +0000 UTC"), Payment: -35333, Interest: -2000, Principal: -33333},
		{Period: 9, StartDate: timeParseUtil(t, "2020-04-23 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-04-23 23:59:59 +0000 UTC"), Payment: -35333, Interest: -2000, Principal: -33333},
		{Period: 10, StartDate: timeParseUtil(t, "2020-04-24 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-04-24 23:59:59 +0000 UTC"), Payment: -35333, Interest: -2000, Principal: -33333},
		{Period: 11, StartDate: timeParseUtil(t, "2020-04-25 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-04-25 23:59:59 +0000 UTC"), Payment: -35333, Interest: -2000, Principal: -33333},
		{Period: 12, StartDate: timeParseUtil(t, "2020-04-26 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-04-26 23:59:59 +0000 UTC"), Payment: -35333, Interest: -2000, Principal: -33333},
		{Period: 13, StartDate: timeParseUtil(t, "2020-04-27 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-04-27 23:59:59 +0000 UTC"), Payment: -35333, Interest: -2000, Principal: -33333},
		{Period: 14, StartDate: timeParseUtil(t, "2020-04-28 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-04-28 23:59:59 +0000 UTC"), Payment: -35333, Interest: -2000, Principal: -33333},
		{Period: 15, StartDate: timeParseUtil(t, "2020-04-29 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-04-29 23:59:59 +0000 UTC"), Payment: -35333, Interest: -2000, Principal: -33333},
		{Period: 16, StartDate: timeParseUtil(t, "2020-04-30 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-04-30 23:59:59 +0000 UTC"), Payment: -35333, Interest: -2000, Principal: -33333},
		{Period: 17, StartDate: timeParseUtil(t, "2020-05-01 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-05-01 23:59:59 +0000 UTC"), Payment: -35333, Interest: -2000, Principal: -33333},
		{Period: 18, StartDate: timeParseUtil(t, "2020-05-02 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-05-02 23:59:59 +0000 UTC"), Payment: -35333, Interest: -2000, Principal: -33333},
		{Period: 19, StartDate: timeParseUtil(t, "2020-05-03 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-05-03 23:59:59 +0000 UTC"), Payment: -35333, Interest: -2000, Principal: -33333},
		{Period: 20, StartDate: timeParseUtil(t, "2020-05-04 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-05-04 23:59:59 +0000 UTC"), Payment: -35333, Interest: -2000, Principal: -33333},
		{Period: 21, StartDate: timeParseUtil(t, "2020-05-05 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-05-05 23:59:59 +0000 UTC"), Payment: -35333, Interest: -2000, Principal: -33333},
		{Period: 22, StartDate: timeParseUtil(t, "2020-05-06 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-05-06 23:59:59 +0000 UTC"), Payment: -35333, Interest: -2000, Principal: -33333},
		{Period: 23, StartDate: timeParseUtil(t, "2020-05-07 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-05-07 23:59:59 +0000 UTC"), Payment: -35333, Interest: -2000, Principal: -33333},
		{Period: 24, StartDate: timeParseUtil(t, "2020-05-08 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-05-08 23:59:59 +0000 UTC"), Payment: -35333, Interest: -2000, Principal: -33333},
		{Period: 25, StartDate: timeParseUtil(t, "2020-05-09 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-05-09 23:59:59 +0000 UTC"), Payment: -35333, Interest: -2000, Principal: -33333},
		{Period: 26, StartDate: timeParseUtil(t, "2020-05-10 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-05-10 23:59:59 +0000 UTC"), Payment: -35333, Interest: -2000, Principal: -33333},
		{Period: 27, StartDate: timeParseUtil(t, "2020-05-11 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-05-11 23:59:59 +0000 UTC"), Payment: -35333, Interest: -2000, Principal: -33333},
		{Period: 28, StartDate: timeParseUtil(t, "2020-05-12 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-05-12 23:59:59 +0000 UTC"), Payment: -35333, Interest: -2000, Principal: -33333},
		{Period: 29, StartDate: timeParseUtil(t, "2020-05-13 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-05-13 23:59:59 +0000 UTC"), Payment: -35333, Interest: -2000, Principal: -33333},
		{Period: 30, StartDate: timeParseUtil(t, "2020-05-14 00:00:00 +0000 UTC"), EndDate: timeParseUtil(t, "2020-05-14 23:59:59 +0000 UTC"), Payment: -35343, Interest: -2000, Principal: -33343},
	}
}

func timeParseUtil(t *testing.T, input string) time.Time {
	resultTime, err := time.Parse("2006-01-02 15:04:05 -0700 MST", input)
	if err != nil {
		t.Fatalf("invalid time format, %v", err)
	}
	return resultTime
}

func getConfigDto(frequency frequency.Type, round bool, interestType interesttype.Type, amount int64, interest int64) *Config {
	return &Config{
		StartDate:      time.Date(2020, 4, 15, 0, 0, 0, 0, time.UTC),
		EndDate:        time.Date(2022, 4, 14, 0, 0, 0, 0, time.UTC),
		Frequency:      frequency,
		AmountBorrowed: amount,
		InterestType:   interestType,
		Interest:       interest,
		Round:          round,
	}
}
