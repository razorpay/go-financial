package calculator

type Flat struct {
}

func (f *Flat) GetPrincipal(config Config, _ int64) float64 {
	return -float64(config.AmountBorrowed) / float64(config.periods)
}

func (f *Flat) GetInterest(config Config, _ int64) float64 {
	return -config.getInterestRatePerPeriodInDecimal() * float64(config.AmountBorrowed)
}

func (f *Flat) GetPayment(config Config) float64 {
	return -(config.getInterestRatePerPeriodInDecimal()*float64(config.periods)*float64(config.AmountBorrowed) +
		float64(config.AmountBorrowed)) / float64(config.periods)
}
