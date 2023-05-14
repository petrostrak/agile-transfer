package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"
)

func (app *application) currencyConvertion(from, to string, amount decimal.Decimal) (decimal.Decimal, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.freecurrencyapi.com/v1/latest?apikey=QAVbfQcb3HY3YDFtDWdIm7yXzGUMymbipsxXYOj6&currencies=%s&base_currency=%s", to, from))
	if err != nil {
		return decimal.Decimal{}, nil
	}

	var convert struct {
		Data map[string]float64
	}

	err = json.NewDecoder(resp.Body).Decode(&convert)
	if err != nil {
		return decimal.Decimal{}, err
	}

	unit, ok := convert.Data[to]
	if !ok {
		return decimal.Decimal{}, errors.New("could not convert currency")
	}

	multiplier := decimal.NewFromFloat(unit)
	if err != nil {
		return decimal.Decimal{}, nil
	}

	return amount.Mul(multiplier), nil
}
