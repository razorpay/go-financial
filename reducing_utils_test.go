package gofinancial

import (
	"errors"
	"testing"

	"github.com/shopspring/decimal"

	"github.com/razorpay/go-financial/enums/paymentperiod"
)

func Test_Pmt(t *testing.T) {
	type args struct {
		rate decimal.Decimal
		nper int64
		pv   decimal.Decimal
		fv   decimal.Decimal
		when paymentperiod.Type
	}
	tests := []struct {
		name string
		args args
		want decimal.Decimal
	}{
		{
			"7.5% p.a., monthly basis, 15 yrs", args{decimal.NewFromFloat(0.075 / 12), 12 * 15, decimal.NewFromInt(200000), decimal.NewFromInt(0), paymentperiod.ENDING}, decimal.NewFromFloat(-1854.0247200054675),
		},
		{
			"bigDecimal case. would give nan if it were float64.", args{decimal.NewFromFloat(12), 400, decimal.NewFromInt(10000), decimal.NewFromInt(5000), paymentperiod.BEGINNING}, decimal.NewFromFloat(-9230.7692307692307692),
		},
		{
			"24% p.a., monthly basis, 2 yrs", args{decimal.NewFromFloat(0.24 / 12), 12 * 2, decimal.NewFromInt(1000000), decimal.NewFromInt(0), 0}, decimal.NewFromFloat(-52871.097253249915),
		},
		{
			"8% p.a., monthly basis, 5 yrs", args{decimal.NewFromFloat(0.08 / 12), 12 * 5, decimal.NewFromInt(15000), decimal.NewFromInt(0), 0}, decimal.NewFromFloat(-304.1459143262052370338701494),
		},
		{
			"0%p.a. , monthly basis, 15 yrs", args{decimal.Zero, 15 * 12, decimal.NewFromInt(200000), decimal.Zero, paymentperiod.ENDING}, decimal.NewFromFloat(-1111.111111111111),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Pmt(tt.args.rate, tt.args.nper, tt.args.pv, tt.args.fv, tt.args.when)
			if err := isAlmostEqual(got, tt.want, decimal.NewFromFloat(precision)); err != nil {
				t.Errorf("error: %v, pmt() = %v, want %v", err.Error(), got.BigInt().String(), tt.want)
			}
		})
	}
}

func Test_Fv(t *testing.T) {
	type args struct {
		rate decimal.Decimal
		nper int64
		pmt  decimal.Decimal
		pv   decimal.Decimal
		when paymentperiod.Type
	}
	tests := []struct {
		name string
		args args
		want decimal.Decimal
	}{
		{
			name: "success", args: args{
				rate: decimal.NewFromFloat(0.05 / 12),
				nper: 10 * 12,
				pmt:  decimal.NewFromInt(-100),
				pv:   decimal.NewFromInt(-100),
				when: paymentperiod.ENDING,
			},
			want: decimal.NewFromFloat(15692.928894335893),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Fv(tt.args.rate, tt.args.nper, tt.args.pmt, tt.args.pv, tt.args.when)
			if err := isAlmostEqual(got, tt.want, decimal.NewFromFloat(precision)); err != nil {
				t.Errorf("error = %v, fv() = %v, want %v", err.Error(), got, tt.want)
			}
		})
	}
}

func Test_IPmt(t *testing.T) {
	type args struct {
		rate decimal.Decimal
		per  int64
		nper int64
		pv   decimal.Decimal
		fv   decimal.Decimal
		when paymentperiod.Type
	}
	tests := []struct {
		name string
		args args
		want decimal.Decimal
	}{
		{
			name: "1period",
			args: args{
				rate: decimal.NewFromFloat(0.0824 / 12),
				per:  1,
				nper: 1 * 12,
				pv:   decimal.NewFromInt(2500),
				fv:   decimal.NewFromInt(0),
				when: paymentperiod.ENDING,
			},
			want: decimal.NewFromFloat(-17.166666666666668),
		},
		{
			name: "2period",
			args: args{
				rate: decimal.NewFromFloat(0.0824 / 12),
				per:  2,
				nper: 1 * 12,
				pv:   decimal.NewFromInt(2500),
				fv:   decimal.NewFromInt(0),
				when: paymentperiod.ENDING,
			},
			want: decimal.NewFromFloat(-15.7893374573507768960793587710732749),
		},
		{
			name: "3period",
			args: args{
				rate: decimal.NewFromFloat(0.0824 / 12),
				per:  3,
				nper: 1 * 12,
				pv:   decimal.NewFromInt(2500),
				fv:   decimal.NewFromInt(0),
				when: paymentperiod.ENDING,
			},
			want: decimal.NewFromFloat(-14.4025505874642504602108554951782324459586257024424875),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IPmt(tt.args.rate, tt.args.per, tt.args.nper, tt.args.pv, tt.args.fv, tt.args.when)
			if err := isAlmostEqual(got, tt.want, decimal.NewFromFloat(precision)); err != nil {
				t.Errorf("error: %v, ipmt() = %v, want %v", err.Error(), got, tt.want)
			}
		})
	}
}

//
func Test_PPmt(t *testing.T) {
	type args struct {
		rate decimal.Decimal
		per  int64
		nper int64
		pv   decimal.Decimal
		fv   decimal.Decimal
		when paymentperiod.Type
	}

	tests := []struct {
		name string
		args args
		want decimal.Decimal
	}{
		{
			name: "period1",
			args: args{
				rate: decimal.NewFromFloat(0.0824 / 12),
				per:  1,
				nper: 1 * 12,
				pv:   decimal.NewFromInt(2500),
				fv:   decimal.NewFromInt(0),
				when: paymentperiod.ENDING,
			},
			want: decimal.NewFromFloat(-200.5819236867801753),
		},
		{
			name: "period2",
			args: args{
				rate: decimal.NewFromFloat(0.0824 / 12),
				per:  2,
				nper: 1 * 12,
				pv:   decimal.NewFromInt(2500),
				fv:   decimal.NewFromInt(0),
				when: paymentperiod.ENDING,
			},
			want: decimal.NewFromFloat(-201.9592528960960659039206412289267251),
		},
		{
			name: "period3",
			args: args{
				rate: decimal.NewFromFloat(0.0824 / 12),
				per:  3,
				nper: 1 * 12,
				pv:   decimal.NewFromInt(2500),
				fv:   decimal.NewFromInt(0),
				when: paymentperiod.ENDING,
			},
			want: decimal.NewFromFloat(-203.3460397659825923397891445048217675540413742975575125),
		},
		{
			name: "period4",
			args: args{
				rate: decimal.NewFromFloat(0.0824 / 12),
				per:  4,
				nper: 1 * 12,
				pv:   decimal.NewFromInt(2500),
				fv:   decimal.NewFromInt(0),
				when: paymentperiod.ENDING,
			},
			want: decimal.NewFromFloat(-204.7423492390423394738416027282628017795870286168635449048641975308641975),
		},
		{
			name: "period5",
			args: args{
				rate: decimal.NewFromFloat(0.0824 / 12),
				per:  5,
				nper: 1 * 12,
				pv:   decimal.NewFromInt(2500),
				fv:   decimal.NewFromInt(0),
				when: paymentperiod.ENDING,
			},
			want: decimal.NewFromFloat(-206.1482467038170969439558163159297469375176170621658196811415267489711932911213991769547325),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PPmt(tt.args.rate, tt.args.per, tt.args.nper, tt.args.pv, tt.args.fv, tt.args.when)
			if err := isAlmostEqual(got, tt.want, decimal.NewFromFloat(precision)); err != nil {
				t.Errorf("error:%v, ppmt() = %v, want %v", err.Error(), got, tt.want)
			}
		})
	}
}

func Test_Pv(t *testing.T) {
	type args struct {
		rate decimal.Decimal
		nper int64
		pmt  decimal.Decimal
		fv   decimal.Decimal
		when paymentperiod.Type
	}
	tests := []struct {
		name string
		args args
		want decimal.Decimal
	}{
		{
			name: "success", args: args{
				rate: decimal.NewFromFloat(0.24 / 12),
				nper: 1 * 12,
				pmt:  decimal.NewFromInt(-300),
				fv:   decimal.NewFromInt(1000),
				when: paymentperiod.BEGINNING,
			},
			want: decimal.NewFromFloat(2447.561238019001),
		}, {
			name: "success", args: args{
				rate: decimal.NewFromFloat(0.24 / 12),
				nper: 1 * 12,
				pmt:  decimal.NewFromInt(-300),
				fv:   decimal.NewFromInt(1000),
				when: paymentperiod.ENDING,
			},
			want: decimal.NewFromFloat(2384.1091906934976),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Pv(tt.args.rate, tt.args.nper, tt.args.pmt, tt.args.fv, tt.args.when)
			if err := isAlmostEqual(got, tt.want, decimal.NewFromFloat(precision)); err != nil {
				t.Errorf("error:%v, pv() = %v, want %v", err, got, tt.want)
			}
		})
	}
}

func Test_Npv(t *testing.T) {
	type args struct {
		rate   decimal.Decimal
		values []decimal.Decimal
	}
	tests := []struct {
		name string
		args args
		want decimal.Decimal
	}{
		{
			name: "success", args: args{
				rate:   decimal.NewFromFloat(0.2),
				values: []decimal.Decimal{decimal.NewFromInt(-1000), decimal.NewFromInt(100), decimal.NewFromInt(100), decimal.NewFromInt(100)},
			},
			want: decimal.NewFromFloat(-789.3518518518518),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Npv(tt.args.rate, tt.args.values)
			if err := isAlmostEqual(got, tt.want, decimal.NewFromFloat(precision)); err != nil {
				t.Errorf("npv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Nper(t *testing.T) {
	type args struct {
		rate decimal.Decimal
		fv   decimal.Decimal
		pmt  decimal.Decimal
		pv   decimal.Decimal
		when paymentperiod.Type
	}
	tests := []struct {
		name    string
		args    args
		want    decimal.Decimal
		wantErr bool
		err     error
	}{
		{
			name: "success", args: args{
				rate: decimal.NewFromFloat(0.07 / 12),
				fv:   decimal.NewFromInt(0),
				pmt:  decimal.NewFromInt(-150),
				pv:   decimal.NewFromInt(8000),
				when: paymentperiod.ENDING,
			},
			want:    decimal.NewFromFloat(64.0733487706618586),
			wantErr: false,
			err:     nil,
		},
		{
			name: "failure", args: args{
				rate: decimal.NewFromFloat(1e100),
				fv:   decimal.NewFromInt(0),
				pmt:  decimal.NewFromInt(-150),
				pv:   decimal.NewFromInt(8000),
				when: paymentperiod.ENDING,
			},
			want:    decimal.NewFromFloat(64.0733487706618586),
			wantErr: true,
			err:     ErrOutOfBounds,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := Nper(tt.args.rate, tt.args.pmt, tt.args.pv, tt.args.fv, tt.args.when); err != nil {
				if !tt.wantErr && errors.Is(err, tt.err) {
					t.Errorf("error is not equal, want=%v, got=%v", tt.err, err)
				}
			} else {
				if err := isAlmostEqual(got, tt.want, decimal.NewFromFloat(precision)); err != nil {
					t.Errorf("fv() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
