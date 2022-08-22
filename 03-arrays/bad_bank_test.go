package arrays

import "testing"

func TestBadBank(t *testing.T) {
	transactions := []Transaction{
		{
			From: "Rodrigo",
			To:   "Amaya",
			Sum:  100,
		},
		{
			From: "Nayra",
			To:   "Rodrigo",
			Sum:  25,
		},
	}

	AssertEqual(t, BalanceFor(transactions, "Amaya"), 100)
	AssertEqual(t, BalanceFor(transactions, "Rodrigo"), -75)
	AssertEqual(t, BalanceFor(transactions, "Nayra"), -25)
}
