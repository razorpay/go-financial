package gofinancial

// Reducing implements financial methods for facilitating a loan use case, following a reducing rate of interest.
type Reducing struct {
}

// GetPrincipal returns principal amount contribution in a given period towards a loan, depending on config.
func (r *Reducing) GetPrincipal(config Config, period int64) float64 {
	return PPmt(config.getInterestRatePerPeriodInDecimal(), period, config.periods, float64(config.AmountBorrowed), 0, config.PaymentPeriod, config.Round)
}

// GetInterest returns interest amount contribution in a given period towards a loan, depending on config.
func (r *Reducing) GetInterest(config Config, period int64) float64 {
	return IPmt(config.getInterestRatePerPeriodInDecimal(), period, config.periods, float64(config.AmountBorrowed), 0, config.PaymentPeriod)
}

// GetPayment returns the periodic payment to be done for a loan depending on config.
func (r *Reducing) GetPayment(config Config) float64 {
	return Pmt(config.getInterestRatePerPeriodInDecimal(), config.periods, float64(config.AmountBorrowed), 0, config.PaymentPeriod)
}
