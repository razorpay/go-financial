package gofinancial

import "github.com/shopspring/decimal"

// Flat implements financial methods for facilitating a loan use case, following a flat rate of interest.
type Flat struct {
}

// GetPrincipal returns principal amount contribution in a given period towards a loan, depending on config.
func (f *Flat) GetPrincipal(config Config, _ int64) decimal.Decimal {
	dPeriod := decimal.NewFromInt(config.periods)
	dAmount := decimal.NewFromInt(config.AmountBorrowed)
	minusOne := decimal.NewFromInt(-1)
	return dAmount.Div(dPeriod).Mul(minusOne)
}

// GetInterest returns interest amount contribution in a given period towards a loan, depending on config.
func (f *Flat) GetInterest(config Config, _ int64) decimal.Decimal {
	dAmount := decimal.NewFromInt(config.AmountBorrowed)
	minusOne := decimal.NewFromInt(-1)
	return config.getInterestRatePerPeriodInDecimal().Mul(dAmount).Mul(minusOne)
}

// GetPayment returns the periodic payment to be done for a loan depending on config.
func (f *Flat) GetPayment(config Config) decimal.Decimal {
	dPeriod := decimal.NewFromInt(config.periods)
	dAmount := decimal.NewFromInt(config.AmountBorrowed)
	minusOne := decimal.NewFromInt(-1)
	totalInterest := config.getInterestRatePerPeriodInDecimal().Mul(dPeriod).Mul(dAmount)
	Payment := totalInterest.Add(dAmount).Mul(minusOne).Div(dPeriod)
	return Payment
}
