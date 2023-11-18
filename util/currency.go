package util

// Map for all supported currencies
type mapEntry struct {
	Key string
	Value string
}

const (
	USD = "USD"
	SGD = "SGD"
	EUR = "EUR"
	RMB = "RMB"
)

// IsSupportedCurrency returns true if currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, SGD, EUR, RMB:
		return true
	}
	return false
}