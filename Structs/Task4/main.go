package main

import (
	"errors"
	"fmt"
)

type BankAccount struct {
	Owner   string
	balance float64
}

func (b *BankAccount) Deposit(amount float64) {
	b.balance += amount
}

func (b *BankAccount) Withdraw(amount float64) error {
	if amount > b.balance {
		return errors.New("insufficient funds")
	}
	b.balance -= amount
	return nil
}

func (b *BankAccount) Balance() float64 {
	return b.balance
}

func main() {
	accounts := []*BankAccount{
		{"Ali", 100},
		{"Assem", 200},
		{"Askar", 300},
	}

	for _, acc := range accounts {
		acc.Deposit(acc.Balance() * 0.05)
	}

	for _, acc := range accounts {
		fmt.Printf("%s: %.2f\n", acc.Owner, acc.Balance())
	}
}
