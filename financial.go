package gofinancial

// Financial interface defines the methods to be over ridden for different financial use cases.
type Financial interface {
	GetPrincipal(config Config, period int64) float64
	GetInterest(config Config, period int64) float64
	GetPayment(config Config) float64
}
