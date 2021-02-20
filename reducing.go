package gofinancial

import "github.com/shopspring/decimal"

// Reducing implements financial methods for facilitating a loan use case, following a reducing rate of interest.
type Reducing struct {
}

// GetPrincipal returns principal amount contribution in a given period towards a loan, depending on config.
func (r *Reducing) GetPrincipal(config Config, period int64) decimal.Decimal {
	return PPmt(config.getInterestRatePerPeriodInDecimal(), period, config.periods, config.AmountBorrowed, decimal.Zero, config.PaymentPeriod)
}

// GetInterest returns interest amount contribution in a given period towards a loan, depending on config.
func (r *Reducing) GetInterest(config Config, period int64) decimal.Decimal {
	return *IPmt(config.getInterestRatePerPeriodInDecimal(), period, config.periods, config.AmountBorrowed, decimal.Zero, config.PaymentPeriod)
}

// GetPayment returns the periodic payment to be done for a loan depending on config.
func (r *Reducing) GetPayment(config Config) decimal.Decimal {
	return Pmt(config.getInterestRatePerPeriodInDecimal(), config.periods, config.AmountBorrowed, decimal.Zero, config.PaymentPeriod)
}
