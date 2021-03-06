package account

import (
	"chain_proto/config"
	"encoding/json"
	"errors"
	"log"

	"github.com/shopspring/decimal"
)

// Account is a model of user in the blockchain.
type Account struct {
	Addr    string
	Balance decimal.Decimal
}

// New initialises a new Account struct
func New(addr string) *Account {
	return &Account{
		Addr:    addr,
		Balance: decimal.New(0, config.MaxDecimalDigit),
	}
}

func (a *Account) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Addr    string `json:"addr"`
		Balance string `json:"balance"`
	}{
		Addr:    a.Addr,
		Balance: a.BalanceString(),
	})
}

func (a *Account) BalanceString() string {
	return a.Balance.StringFixed(config.MaxDecimalDigit)
}

// Send sends a specified amount of coins from an account to another
func (a *Account) Send(amount decimal.Decimal, recipient *Account) error {
	if a.Balance.LessThan(amount) {
		return errors.New("error: balance not enough")
	}
	a.Balance.Sub(amount)
	recipient.Receive(amount)
	return nil
}

// Receive sums a specified amount of coins to the current balance
func (a *Account) Receive(amount decimal.Decimal) {
	log.Printf("info: account(%s) received %s coins\n", a.Addr, amount.String())
	a.Balance = a.Balance.Add(amount)
}
