package arrays

type Transaction struct {
	From string
	To   string
	Sum  float64
}

func BalanceFor(transactions []Transaction, name string) float64 {
	var balance float64
	for _, transaction := range transactions {
		if transaction.From == name {
			balance -= transaction.Sum
		}

		if transaction.To == name {
			balance += transaction.Sum
		}
	}

	return balance
}
