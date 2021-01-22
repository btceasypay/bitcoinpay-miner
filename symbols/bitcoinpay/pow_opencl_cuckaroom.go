//+build opencl,!cuda

/**
Bitcoinpay
james
*/
package bitcoinpay

import (
	"github.com/btceasypay/bitcoinpay-miner/common"
	"github.com/btceasypay/bitcoinpay-miner/core"
	"os"
	"sync"
)

type CudaCuckaroom struct {
	core.Device
	ClearBytes    []byte
	Work          *BitcoinpayWork
	header        MinerBlockData
	EdgeBits      int
	Step          int
	WorkGroupSize int
	LocalSize     int
	Nedge         int
	Edgemask      uint64
	Nonces        []uint32
}

func (this *CudaCuckaroom) InitDevice() {
	common.MinerLoger.Error("AMD Not Support CUDA!")
	os.Exit(1)
}

func (this *CudaCuckaroom) Mine(wg *sync.WaitGroup) {
	defer this.Release()
	defer wg.Done()
}

func (this *CudaCuckaroom) SubmitShare(substr chan string) {
	this.Device.SubmitShare(substr)
}
