package coinbase

import (
	"github.com/btceasypay/bitcoinpay/core/types"
)

type Coinbase085 struct {
	CoinbaseBase
}

//in bitcoinpay 0.8.5 the getblocktemplate result about coinbasevalue is contain tx fee
//so if we want get subsidy may need minus tx fee
//the miner coinbase tx need contain tx fee
func (this *Coinbase085) GetCoinbaseTx() *types.Tx {
	subsidy := this.CoinbaseValue
	subsidy -= this.TotalFee
	coinbaseTx := this.CalcCoinbaseTx(subsidy)
	if coinbaseTx == nil {
		return nil
	}
	coinbaseTx.Tx.TxOut[0].Amount += uint64(this.PackTxTotalFee)
	return coinbaseTx
}
