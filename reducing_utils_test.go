package calculator

import (
	"fmt"
	"math"
	"testing"

	"github.com/smartystreets/assertions"

	"github.com/razorpay/go-financial/enums/paymentperiod"
)

func Test_pmt(t *testing.T) {
	type args struct {
		rate float64
		nper int64
		pv   float64
		fv   float64
		when paymentperiod.Type
	}
	tests := []struct {
		name string
		args args
		want float64
	}{

		{
			"7.5% p.a., monthly basis, 15 yrs", args{0.075 / 12, 12 * 15, 200000, 0, paymentperiod.ENDING}, -1854.0247200054675,
		},
		{
			"nan case", args{12, 400, 10000, 5000, paymentperiod.BEGINNING}, math.NaN(),
		},
		{
			"24% p.a., monthly basis, 2 yrs", args{0.24 / 12, 12 * 2, 1000000, 0, 0}, -52871.097253249915,
		},
		{
			"8% p.a., monthly basis, 5 yrs", args{0.08 / 12, 12 * 5, 15000, 0, 0}, -304.1459143262052370338701494,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pmt(tt.args.rate, tt.args.nper, tt.args.pv, tt.args.fv, tt.args.when); assertions.ShouldAlmostEqual(got, tt.want) != "" {
				if math.IsNaN(tt.want) && math.IsNaN(got) {
					return
				}
				fmt.Printf("%0.f", got)
				t.Errorf("pmt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fv(t *testing.T) {
	type args struct {
		rate float64
		nper int64
		pmt  float64
		pv   float64
		when paymentperiod.Type
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "success", args: args{
				rate: 0.05 / 12,
				nper: 10 * 12,
				pmt:  -100,
				pv:   -100,
				when: paymentperiod.ENDING,
			},
			want: 15692.928894335893,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fv(tt.args.rate, tt.args.nper, tt.args.pmt, tt.args.pv, tt.args.when); got != tt.want {
				t.Errorf("fv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ipmt(t *testing.T) {
	type args struct {
		rate float64
		per  int64
		nper int64
		pv   float64
		fv   float64
		when paymentperiod.Type
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "1period",
			args: args{
				rate: 0.0824 / 12,
				per:  1,
				nper: 1 * 12,
				pv:   2500,
				fv:   0,
				when: paymentperiod.ENDING,
			},
			want: -17.166666666666668,
		},
		{
			name: "2period",
			args: args{
				rate: 0.0824 / 12,
				per:  2,
				nper: 1 * 12,
				pv:   2500,
				fv:   0,
				when: paymentperiod.ENDING,
			},
			want: -15.78933745735078,
		},
		{
			name: "3period",
			args: args{
				rate: 0.0824 / 12,
				per:  3,
				nper: 1 * 12,
				pv:   2500,
				fv:   0,
				when: paymentperiod.ENDING,
			},
			want: -14.402550587464265,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ipmt(tt.args.rate, tt.args.per, tt.args.nper, tt.args.pv, tt.args.fv, tt.args.when); assertions.ShouldAlmostEqual(got, tt.want) != "" {
				t.Errorf("ipmt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ppmt(t *testing.T) {
	type args struct {
		rate  float64
		per   int64
		nper  int64
		pv    float64
		fv    float64
		when  paymentperiod.Type
		round bool
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "period1",
			args: args{
				rate:  0.0824 / 12,
				per:   1,
				nper:  1 * 12,
				pv:    2500,
				fv:    0,
				when:  paymentperiod.ENDING,
				round: true,
			},
			want: -201,
		},
		{
			name: "period2",
			args: args{
				rate:  0.0824 / 12,
				per:   2,
				nper:  1 * 12,
				pv:    2500,
				fv:    0,
				when:  paymentperiod.ENDING,
				round: true,
			},
			want: -202,
		},
		{
			name: "period3",
			args: args{
				rate:  0.0824 / 12,
				per:   3,
				nper:  1 * 12,
				pv:    2500,
				fv:    0,
				when:  paymentperiod.ENDING,
				round: true,
			},
			want: -204,
		},
		{
			name: "period4",
			args: args{
				rate:  0.0824 / 12,
				per:   4,
				nper:  1 * 12,
				pv:    2500,
				fv:    0,
				when:  paymentperiod.ENDING,
				round: true,
			},
			want: -205,
		},
		{
			name: "period5",
			args: args{
				rate:  0.0824 / 12,
				per:   5,
				nper:  1 * 12,
				pv:    2500,
				fv:    0,
				when:  paymentperiod.ENDING,
				round: true,
			},
			want: -206,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ppmt(tt.args.rate, tt.args.per, tt.args.nper, tt.args.pv, tt.args.fv, tt.args.when, tt.args.round); assertions.ShouldAlmostEqual(got, tt.want) != "" {
				t.Errorf("ppmt() = %v, want %v", got, tt.want)
			}
		})
	}
}
