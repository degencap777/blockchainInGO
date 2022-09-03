package transaction

type TXInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
}
