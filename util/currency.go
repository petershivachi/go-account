package util

const (
	USD = "USD"
	EUR = "EUR"
	GBP = "GBP"
	JPY = "JPY"
	NZD = "NZD"
	RUB = "RUB"
	CAD = "CAD"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, GBP, JPY, NZD, RUB, CAD:
		return true
	}
	return false
}
