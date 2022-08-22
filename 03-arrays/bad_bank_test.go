package arrays

import "testing"

func TestBadBank(t *testing.T) {
	amaya := Account{Name: "Amaya", Balance: 100}
	rodrigo := Account{Name: "Rodrigo", Balance: 75}
	nayra := Account{Name: "Nayra", Balance: 200}
	transactions := []Transaction{
		NewTransaction(rodrigo, amaya, 100),
		NewTransaction(nayra, rodrigo, 25),
	}

	newBalanceFor := func(account Account) float64 {
		return NewBalanceFor(account, transactions).Balance
	}

	AssertEqual(t, newBalanceFor(amaya), 200)
	AssertEqual(t, newBalanceFor(rodrigo), 0)
	AssertEqual(t, newBalanceFor(nayra), 175)
}
