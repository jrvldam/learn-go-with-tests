package arrays

type Transaction struct {
	From string
	To   string
	Sum  float64
}

func BalanceFor(transactions []Transaction, name string) float64 {
	adjustBalance := func(currentBalance float64, transaction Transaction) float64 {
		if transaction.From == name {
			return currentBalance - transaction.Sum
		}

		if transaction.To == name {
			return currentBalance + transaction.Sum
		}

		return currentBalance
	}

	return Reduce(transactions, adjustBalance, 0.0)
}
