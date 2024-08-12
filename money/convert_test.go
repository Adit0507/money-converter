package money_test

import (
	"moneyconverter/money"
	"reflect"
	"testing"
)

func mustParseCurrency(t *testing.T, code string) money.Currency {
	t.Helper()
	currency, err := money.ParseCurrency(code)
	if err != nil {
		t.Fatalf("cannot parse currency %s code", code)
	}
	return currency
}

func mustParseAmount(t *testing.T, value string, code string) money.Amount {
	t.Helper()
	n, err := money.ParseDecimal(value)
	if err != nil {
		t.Fatalf("invalid number: %s", value)
	}
	currency, err := money.ParseCurrency(code)
	if err != nil {
		t.Fatalf("invalid currency code: %s", code)
	}
	amount, err := money.NewAmount(n, currency)
	if err != nil {
		t.Fatalf("cannot create amount with value %v and currency code %s", n, code)
	}
	return amount
}

type stubRate struct {
	rate string
	err  error
}

func (m stubRate) FetchExchangeRate(_, _ money.Currency) (money.ExchangeRate, error) {
	rate, _ := money.ParseDecimal(m.rate)
	return money.ExchangeRate(rate), m.err
}

func TestConvert(t *testing.T) {
	tt := map[string]struct {
		amount   money.Amount
		to       money.Currency
		stub     stubRate
		validate func(t *testing.T, got money.Amount, err error)
	}{
		"34.98 USD to EUR": {
			amount: mustParseAmount(t, "34.98", "USD"),
			to:     mustParseCurrency(t, "EUR"),
			validate: func(t *testing.T, got money.Amount, err error) {
				if err != nil {
					t.Errorf("expected no error, got %s", err.Error())
				}
				expected := money.Amount{}
				if !reflect.DeepEqual(got, expected) {
					t.Errorf("expected %v, got %v", expected, got)
				}
			},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got, err := money.Convert(tc.amount, tc.to, tc.stub)
			tc.validate(t, got, err)
		})
	}
}
