package ecbank

import (
	"errors"
	"fmt"
	"moneyconverter/money"
	"net/http"
	"net/url"
	"time"
)

const (
	ErrCallingServer      = ecbankError("error calling server")
	ErrClientSide         = ecbankError("client side error when contacting ECB")
	ErrServerSide         = ecbankError("server side error when contacting ECB")
	ErrUnknownStatusCode  = ecbankError("unknown error when contacting ECB")
	ErrUnexpectedFormat   = ecbankError("unexpected response format")
	ErrChangeRateNotFound = ecbankError("couldn't find the exchange rate")
	ErrTimeout            = ecbankError("timed out when waiting for response")
	clientErrorClass      = 4
	serverErrorClass      = 5
)

type Client struct {
	client *http.Client
}

func NewClient(timeout time.Duration) Client{
	return Client{
		client: &http.Client{Timeout: timeout},
	}
}

// retuns class of http status code
func httpStatusClass(statusCode int) int {
	const httpErrorClassSize = 100
	return statusCode / httpErrorClassSize
}

// returns a different error depending on the returned status code.
func checkStatusCode(statusCode int) error {
	switch {
	case statusCode == http.StatusOK:
		return nil

	case httpStatusClass(statusCode) == clientErrorClass:
		return fmt.Errorf("%w: %d", ErrClientSide, statusCode)
	case httpStatusClass(statusCode) == serverErrorClass:
		return fmt.Errorf("%w: %d", ErrServerSide, statusCode)
	default:
		return fmt.Errorf("%w: %d", ErrUnknownStatusCode, statusCode)
	}
}

func (c Client) FetchExchangeRate(source, target money.Currency) (money.ExchangeRate, error) {
	const euroxrefURL = "http://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"
	resp, err := c.client.Get(euroxrefURL)
	if err != nil {
		var urlErr *url.Error
		if ok := errors.As(err, &urlErr); ok && urlErr.Timeout() {
			return money.ExchangeRate{}, fmt.Errorf("%w: %s", ErrTimeout, err.Error())
		}
		return money.ExchangeRate{}, fmt.Errorf("%w: %s", ErrCallingServer, err.Error())
	}

	defer resp.Body.Close()

	if err = checkStatusCode(resp.StatusCode); err != nil {
		return money.ExchangeRate{}, err
	}

	rate, err := readRateFromResponse(source.Code(), target.Code(), resp.Body)
	if err != nil {
		return money.ExchangeRate{}, err
	}

	return rate, nil
}
