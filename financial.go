package calculator

type Financial interface {
	GetPrincipal(config Config, period int64) float64
	GetInterest(config Config, period int64) float64
	GetPayment(config Config) float64
}
