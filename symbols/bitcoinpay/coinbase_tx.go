// Copyright (c) 2019 The bitcoinpay developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
package bitcoinpay

import (
	"fmt"
	"github.com/btceasypay/bitcoinpay-miner/common"
	"github.com/btceasypay/bitcoinpay-miner/symbols/bitcoinpay/coinbase"
	"github.com/btceasypay/bitcoinpay/common/hash"
	"github.com/btceasypay/bitcoinpay/core/types"
	"sort"
)

//calc coinbase
func (h *BlockHeader) CalcCoinBase(cfg *common.GlobalConfig, coinbaseStr string, extraNonce uint64, payAddressS string) (*hash.Hash, []Transactions) {
	transactions := make(Transactionses, 0)
	if !h.HasCoinbasePack {
		h.TotalFee = 0
		for i := 0; i < len(h.Transactions); i++ {
			transactions = append(transactions, h.Transactions[i])
		}
		sort.Sort(transactions)
		for i := 0; i < len(transactions); i++ {
			h.TotalFee += transactions[i].Fee
		}
	}
	transactions = make(Transactionses, 0)
	totalTxFee := uint64(0)
	if !h.HasCoinbasePack {
		h.transactions = make([]*types.Tx, 0)
		tmpTrx := make(Transactionses, 0)
		for i := 0; i < len(h.Transactions); i++ {
			tmpTrx = append(tmpTrx, h.Transactions[i])
		}
		sort.Sort(tmpTrx)
		allSigCount := 0
		//every time pack max 1000 transactions and max 5000 sign scripts
		txCount := len(tmpTrx)
		if txCount > (cfg.OptionConfig.MaxTxCount - 1) {
			txCount = cfg.OptionConfig.MaxTxCount - 1
		}
		for i := 0; i < txCount; i++ {
			if allSigCount > (cfg.OptionConfig.MaxSigCount - 1) {
				break
			}
			transactions = append(transactions, tmpTrx[i])
			allSigCount += tmpTrx[i].GetSigCount()
			h.transactions = append(h.transactions, tmpTrx[i].EncodeTx())
		}
		for i := 0; i < len(transactions); i++ {
			totalTxFee += transactions[i].Fee
		}
	} else {
		for i := 1; i < len(h.Transactions); i++ {
			totalTxFee += h.Transactions[i].Fee
		}
	}
	instance := coinbase.GetNewCoinbaseInstance(int(h.Version), cfg.NecessaryConfig.Param, payAddressS, coinbaseStr, extraNonce, h.Height, h.TotalFee, uint64(h.Coinbasevalue), totalTxFee)
	// miner get tx tax
	coinbaseTx := instance.GetCoinbaseTx()
	if coinbaseTx == nil {
		return nil, []Transactions{}
	}
	h.AddCoinbaseTx(coinbaseTx)
	coinbaseTx = fillWitnessToCoinBase(h.transactions)
	txBuf, err := coinbaseTx.Tx.Serialize()
	if err != nil {
		context := "Failed to serialize transaction"
		common.MinerLoger.Error(context)
		return nil, []Transactions{}
	}
	coinbaseData := fmt.Sprintf("%x", txBuf)
	if !h.HasCoinbasePack {
		newtransactions := make(Transactionses, 0)
		newtransactions = append(newtransactions, Transactions{coinbaseTx.Tx.TxHash(), coinbaseData, 0})
		newtransactions = append(newtransactions, transactions...)
		h.Transactions = newtransactions
		h.HasCoinbasePack = true
	} else {
		h.Transactions[0] = Transactions{coinbaseTx.Tx.TxHash(), coinbaseData, 0}
	}
	ha := h.BuildMerkleTreeStore(0)
	return &ha, h.Transactions
}

func (h *BlockHeader) AddCoinbaseTx(coinbaseTx *types.Tx) {
	if h.HasCoinbasePack {
		h.transactions[0] = coinbaseTx
	} else {
		txs := make([]*types.Tx, 0)
		txs = append(txs, coinbaseTx)
		txs = append(txs, h.transactions...)
		h.transactions = txs
	}
}

func fillWitnessToCoinBase(blockTxns []*types.Tx) *types.Tx {
	merkles := BuildMerkleTreeStoreWithness(blockTxns, true)
	txWitnessRoot := merkles[len(merkles)-1]
	witnessPreimage := append(txWitnessRoot.Bytes(), blockTxns[0].Tx.TxIn[0].SignScript...)
	witnessCommitment := hash.DoubleHashH(witnessPreimage[:])
	blockTxns[0].Tx.TxIn[0].PreviousOut.Hash = witnessCommitment
	blockTxns[0].RefreshHash()
	return blockTxns[0]
}
