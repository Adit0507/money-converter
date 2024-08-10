package money

type Amount struct {
	quantity Decimal
	currency Currency
}

const ErrTooPrecise = Error("quantity is too precise")

func NewAmount(quantity Decimal, currency Currency) (Amount, error) {
	if quantity.precision > currency.precision {
		return Amount{}, ErrTooPrecise
	}

	return Amount{quantity: quantity, currency: currency}, nil
}
