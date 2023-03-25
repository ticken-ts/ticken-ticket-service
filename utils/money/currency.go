package money

import (
	"fmt"
)

type Currency struct {
	Name   string
	Symbol string
	Type   CurrencyType
}

func ParseCurrency(symbol string) (Currency, error) {
	switch symbol {
	case Ether.Symbol:
		return Ether, nil
	case PesoArg.Symbol:
		return PesoArg, nil
	default:
		return Currency{}, fmt.Errorf("currency %s not supported", symbol)
	}
}

type CurrencyType uint8

const (
	CryptoCurrencyType CurrencyType = 1
	FiatCurrencyType   CurrencyType = 2
)

var (
	Ether   = Currency{Name: "Ether", Symbol: "ETH", Type: CryptoCurrencyType}
	PesoArg = Currency{Name: "Pesos (ARG)", Symbol: "AR$", Type: FiatCurrencyType}
)

type Money struct {
	Currency Currency
	Amount   float64
}

func (money *Money) IsCrypto() bool {
	return money.Currency.Type == CryptoCurrencyType
}

func (money *Money) IsFiat() bool {
	return money.Currency.Type == FiatCurrencyType
}

func (money *Money) String() string {
	return fmt.Sprintf("%.2f %s", money.Amount, money.Currency.Symbol)
}

func BuildFrom(amount float64, currencySymbol string) (*Money, error) {
	curr, err := ParseCurrency(currencySymbol)
	if err != nil {
		return nil, err
	}
	return &Money{Amount: amount, Currency: curr}, nil
}
