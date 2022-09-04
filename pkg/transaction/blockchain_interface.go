package transaction

type blockchain interface {
	FindSpendableOutputs(address string, amount int) (int, map[string][]int)
}
