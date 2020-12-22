package calculator

import (
	"fmt"
	"testing"
	"time"

	Frequency "github.com/razorpay/go-financial/enums/frequency"
)

type dateGroup struct {
	startDate time.Time
	endDate   time.Time
}

func TestConfig_SetPeriodsAndDates(t *testing.T) {
	type fields struct {
		StartDate time.Time
		EndDate   time.Time
		Frequency Frequency.Type
	}
	// start inclusive.
	// end not inclusive

	tests := []struct {
		name        string
		fields      fields
		wantErr     bool
		wantPeriods int64
		wantDates   []dateGroup
	}{
		{
			name: "daily same year", fields: fields{getDate(2020, 1, 1), getDate(2020, 1, 31), Frequency.DAILY}, wantErr: false,
			wantPeriods: 31, wantDates: []dateGroup{
				{timeParseUtil(t, "2020-01-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-01 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-02 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-02 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-03 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-03 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-04 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-04 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-05 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-05 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-06 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-06 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-07 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-07 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-08 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-08 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-09 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-09 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-10 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-10 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-11 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-11 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-12 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-12 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-13 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-13 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-14 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-14 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-15 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-15 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-16 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-16 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-17 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-17 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-18 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-18 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-19 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-19 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-20 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-20 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-21 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-21 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-22 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-22 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-23 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-23 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-24 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-24 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-25 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-25 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-26 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-26 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-27 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-27 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-28 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-28 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-29 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-29 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-30 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-30 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-31 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-31 23:59:59 +0000 UTC")},
			},
		},
		{
			name: "weekly same year", fields: fields{getDate(2020, 1, 1), getDate(2020, 4, 14), Frequency.WEEKLY}, wantErr: false,
			wantPeriods: 15, wantDates: []dateGroup{
				{timeParseUtil(t, "2020-01-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-07 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-08 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-14 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-15 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-21 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-22 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-28 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-29 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-02-04 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-02-05 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-02-11 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-02-12 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-02-18 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-02-19 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-02-25 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-02-26 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-03-03 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-03-04 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-03-10 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-03-11 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-03-17 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-03-18 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-03-24 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-03-25 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-03-31 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-04-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-04-07 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-04-08 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-04-14 23:59:59 +0000 UTC")},
			},
		},
		{
			name: "weekly different year", fields: fields{getDate(2020, 1, 1), getDate(2021, 2, 23), Frequency.WEEKLY}, wantErr: false,
			wantPeriods: 60, wantDates: []dateGroup{
				{timeParseUtil(t, "2020-01-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-07 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-08 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-14 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-15 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-21 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-22 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-28 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-01-29 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-02-04 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-02-05 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-02-11 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-02-12 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-02-18 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-02-19 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-02-25 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-02-26 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-03-03 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-03-04 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-03-10 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-03-11 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-03-17 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-03-18 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-03-24 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-03-25 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-03-31 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-04-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-04-07 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-04-08 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-04-14 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-04-15 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-04-21 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-04-22 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-04-28 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-04-29 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-05-05 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-05-06 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-05-12 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-05-13 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-05-19 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-05-20 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-05-26 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-05-27 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-06-02 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-06-03 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-06-09 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-06-10 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-06-16 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-06-17 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-06-23 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-06-24 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-06-30 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-07-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-07-07 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-07-08 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-07-14 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-07-15 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-07-21 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-07-22 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-07-28 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-07-29 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-08-04 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-08-05 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-08-11 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-08-12 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-08-18 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-08-19 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-08-25 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-08-26 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-09-01 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-09-02 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-09-08 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-09-09 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-09-15 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-09-16 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-09-22 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-09-23 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-09-29 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-09-30 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-10-06 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-10-07 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-10-13 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-10-14 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-10-20 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-10-21 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-10-27 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-10-28 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-11-03 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-11-04 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-11-10 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-11-11 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-11-17 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-11-18 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-11-24 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-11-25 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-12-01 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-12-02 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-12-08 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-12-09 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-12-15 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-12-16 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-12-22 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-12-23 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-12-29 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-12-30 00:00:00 +0000 UTC"), timeParseUtil(t, "2021-01-05 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2021-01-06 00:00:00 +0000 UTC"), timeParseUtil(t, "2021-01-12 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2021-01-13 00:00:00 +0000 UTC"), timeParseUtil(t, "2021-01-19 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2021-01-20 00:00:00 +0000 UTC"), timeParseUtil(t, "2021-01-26 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2021-01-27 00:00:00 +0000 UTC"), timeParseUtil(t, "2021-02-02 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2021-02-03 00:00:00 +0000 UTC"), timeParseUtil(t, "2021-02-09 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2021-02-10 00:00:00 +0000 UTC"), timeParseUtil(t, "2021-02-16 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2021-02-17 00:00:00 +0000 UTC"), timeParseUtil(t, "2021-02-23 23:59:59 +0000 UTC")},
			},
		},

		{
			name: "monthly same year", fields: fields{getDate(2020, 1, 1), getDate(2020, 9, 30), Frequency.MONTHLY}, wantErr: false,
			wantPeriods: 9, wantDates: []dateGroup{
				{timeParseUtil(t, "2020-01-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-31 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-02-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-02-29 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-03-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-03-31 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-04-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-04-30 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-05-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-05-31 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-06-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-06-30 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-07-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-07-31 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-08-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-08-31 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-09-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-09-30 23:59:59 +0000 UTC")},
			},
		},
		{
			name: "monthly different year", fields: fields{getDate(2020, 1, 1), getDate(2021, 10, 31), Frequency.MONTHLY}, wantErr: false,
			wantPeriods: 22, wantDates: []dateGroup{
				{timeParseUtil(t, "2020-01-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-01-31 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-02-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-02-29 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-03-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-03-31 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-04-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-04-30 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-05-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-05-31 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-06-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-06-30 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-07-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-07-31 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-08-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-08-31 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-09-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-09-30 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-10-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-10-31 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-11-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-11-30 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2020-12-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2020-12-31 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2021-01-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2021-01-31 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2021-02-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2021-02-28 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2021-03-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2021-03-31 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2021-04-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2021-04-30 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2021-05-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2021-05-31 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2021-06-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2021-06-30 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2021-07-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2021-07-31 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2021-08-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2021-08-31 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2021-09-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2021-09-30 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2021-10-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2021-10-31 23:59:59 +0000 UTC")},
			},
		},
		{
			name: "annually", fields: fields{getDate(2020, 5, 1), getDate(2022, 4, 30), Frequency.ANNUALLY}, wantErr: false,
			wantPeriods: 2, wantDates: []dateGroup{
				{timeParseUtil(t, "2020-05-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2021-04-30 23:59:59 +0000 UTC")},
				{timeParseUtil(t, "2021-05-01 00:00:00 +0000 UTC"), timeParseUtil(t, "2022-04-30 23:59:59 +0000 UTC")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				StartDate: tt.fields.StartDate,
				EndDate:   tt.fields.EndDate,
				Frequency: tt.fields.Frequency,
			}
			if err := c.SetPeriodsAndDates(); (err != nil) != tt.wantErr {
				t.Fatalf("SetPeriodsAndDates() error = %v, wantErr %v", err, tt.wantErr)
			}
			if c.periods != tt.wantPeriods {
				t.Fatalf("want periods: %v, got periods:%v", tt.wantPeriods, c.periods)
			}
			if err := areDatesEqual(c.startDates, c.endDates, tt.wantDates); err != nil {
				t.Fatalf("dates are not equal. error:%v", err)
			}
		})
	}
}

func areDatesEqual(actualStartDates []time.Time, actualEndDates []time.Time, expected []dateGroup) error {
	for idx := range expected {
		if !actualStartDates[idx].Equal(expected[idx].startDate) || !actualEndDates[idx].Equal(expected[idx].endDate) {
			return fmt.Errorf("expected startDate:%v, endDate:%v, got startDate:%v endDate:%v", expected[idx].startDate, expected[idx].endDate, actualStartDates[idx], actualEndDates[idx])
		}
	}
	return nil
}

func getDate(y int, m int, d int) time.Time {
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
}
