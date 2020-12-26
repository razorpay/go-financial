package gofinancial

// Flat implements financial methods for facilitating a loan use case, following a flat rate of interest.
type Flat struct {
}

// GetPrincipal returns principal amount contribution in a given period towards a loan, depending on config.
func (f *Flat) GetPrincipal(config Config, _ int64) float64 {
	return -float64(config.AmountBorrowed) / float64(config.periods)
}

// GetInterest returns interest amount contribution in a given period towards a loan, depending on config.
func (f *Flat) GetInterest(config Config, _ int64) float64 {
	return -config.getInterestRatePerPeriodInDecimal() * float64(config.AmountBorrowed)
}

// GetPayment returns the periodic payment to be done for a loan depending on config.
func (f *Flat) GetPayment(config Config) float64 {
	return -(config.getInterestRatePerPeriodInDecimal()*float64(config.periods)*float64(config.AmountBorrowed) +
		float64(config.AmountBorrowed)) / float64(config.periods)
}
