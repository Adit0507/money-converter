package money

// applies the change rate to convert an amt. to a target currency
func Convert(amount Amount, to Currency) (Amount, error) {
	convertedValue := applyExchangeRate(amount, to, ExchangeRate{subunits: 2, precision: 0})

	if err := convertedValue.validate(); err != nil {
		return Amount{}, err
	}

	return convertedValue, nil
}

func multiply(d Decimal, r ExchangeRate) Decimal {
	return Decimal{
		subunits:  d.subunits * r.subunits,
		precision: d.precision + r.precision,
	}
}

type ExchangeRate Decimal

// returns a new amount representung the input
func applyExchangeRate(a Amount, target Currency, rate ExchangeRate) Amount {
	converted := multiply(a.quantity, rate)

	switch {
	case converted.precision > target.precision:
		converted.subunits = converted.subunits / pow10(converted.precision-target.precision)

	case converted.precision < target.precision:
		converted.subunits = converted.subunits * pow10(target.precision-converted.precision)
	}

	converted.precision = target.precision

	return Amount{
		currency: target,
		quantity: converted,
	}
}