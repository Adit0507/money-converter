package main

import (
	"flag"
	"fmt"
	"moneyconverter/ecbank"
	"moneyconverter/money"
	"os"
	"time"
)

func main() {
	from := flag.String("from", "", "source currency, required")
	to := flag.String("to", "EUR", "target currency")
	flag.Parse()

	targetCurrency, err := money.ParseCurrency(*to)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to parse target currency %q: %s.\n", *to, err.Error())
		os.Exit(1)
	}
	
	amount := parseAmount(from)
	rates :=ecbank.NewClient(30*time.Second)

	convertedAmount, err := money.Convert(amount, targetCurrency, rates)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to convert %s to %s: %s.\n", amount, targetCurrency, err.Error())
		os.Exit(1)
	}

	fmt.Printf("%s = %s \n", amount, convertedAmount)
}

func parseAmount(from *string) money.Amount {
	fromCurrency, err := money.ParseCurrency(*from)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to parse source currency %q: %s.\n", *from, err.Error())
		os.Exit(1)
	}
	
	value := flag.Arg(0)
	if value == "" {
		_, _ = fmt.Fprintln(os.Stderr, "missing amount to convert")
		flag.Usage()
		os.Exit(1)
	}

	// parsing quantitty to convert
	quantity, err := money.ParseDecimal(value)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to parse value %q: %s.\n", value, err.Error())
		os.Exit(1)
	}

	amount, err := money.NewAmount(quantity, fromCurrency)
	if err!= nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	return amount
}