package bitcoinpay

import (
	"bytes"
	"encoding/hex"
	"github.com/btceasypay/bitcoinpay/common/hash"
	"github.com/btceasypay/bitcoinpay/core/message"
	"github.com/btceasypay/bitcoinpay/core/protocol"
	"github.com/btceasypay/bitcoinpay/core/types"
)

type ParentItems struct {
	Hash hash.Hash `json:"hash"`
	Data string    `json:"data"`
}

type Transactions struct {
	Hash hash.Hash `json:"hash"`
	Data string    `json:"data"`
	Fee  uint64    `json:"fee"`
}

type Transactionses []Transactions

func (p Transactionses) Len() int { return len(p) }

// fee sort desc
func (p Transactionses) Less(i, j int) bool {
	return p[i].Fee > p[j].Fee
}
func (p Transactionses) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (this *Transactions) GetSigCount() int {
	txBytes, _ := hex.DecodeString(this.Data)
	var mtx = new(message.MsgTx)
	_ = mtx.Decode(bytes.NewReader(txBytes), protocol.ProtocolVersion)
	return len(mtx.Tx.TxOut)
}

func (this *Transactions) EncodeTx() *types.Tx {
	txBytes, _ := hex.DecodeString(this.Data)
	var mtx = new(message.MsgTx)
	_ = mtx.Decode(bytes.NewReader(txBytes), protocol.ProtocolVersion)
	tx := types.NewTx(mtx.Tx)
	return tx
}
