package calculator

type Reducing struct {
}

func (r *Reducing) GetPrincipal(config Config, period int64) float64 {
	return ppmt(config.getInterestRatePerPeriodInDecimal(), period, config.periods, float64(config.AmountBorrowed), 0, config.PaymentPeriod, config.Round)
}

func (r *Reducing) GetInterest(config Config, period int64) float64 {
	return ipmt(config.getInterestRatePerPeriodInDecimal(), period, config.periods, float64(config.AmountBorrowed), 0, config.PaymentPeriod)
}

func (r *Reducing) GetPayment(config Config) float64 {
	return pmt(config.getInterestRatePerPeriodInDecimal(), config.periods, float64(config.AmountBorrowed), 0, config.PaymentPeriod)
}
