package money

import "fmt"

type ratesFetcher interface {
	FetchExchangeRate(source, target Currency) (ExchangeRate, error)
}

// applies the change rate to convert an amt. to a target currency
func Convert(amount Amount, to Currency, rates ratesFetcher) (Amount, error) {
	r, err := rates.FetchExchangeRate(amount.currency, to)
	if err != nil {
		return Amount{}, fmt.Errorf("cannt get change rate: %w", err)
	}

	convertedValue, err := applyExchangeRate(amount, to, r)
	if err != nil {
		return Amount{}, err
	}

	return convertedValue, nil
}

func multiply(d Decimal, r ExchangeRate) Decimal {
	dec := Decimal{
		subunits:  d.subunits * r.subunits,
		precision: d.precision + r.precision,
	}

	dec.simplify()
	return dec
}

type ExchangeRate Decimal

// returns a new amount representung the input
func applyExchangeRate(a Amount, target Currency, rate ExchangeRate) (Amount, error) {
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
	}, nil
}