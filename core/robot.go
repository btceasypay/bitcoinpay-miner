// Copyright (c) 2019 The bitcoinpay developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
package core

import (
	"github.com/btceasypay/bitcoinpay-miner/cl"
	"github.com/btceasypay/bitcoinpay-miner/common"
	"os"
	"strings"
	"sync"
)

const (
	SYMBOL_BTP = "BTP"
)

//var devicesTypesForMining = cl.DeviceTypeAll

type Robot interface {
	Run()        // uses device to calulate the nonce
	ListenWork() //listen the solo or pool work
	SubmitWork() //submit the work
}

type MinerRobot struct {
	Cfg              *common.GlobalConfig //config
	ValidShares      uint64
	StaleShares      uint64
	InvalidShares    uint64
	AllDiffOneShares uint64
	Wg               *sync.WaitGroup
	Started          uint32
	Quit             chan os.Signal
	Work             *Work
	ClDevices        []*cl.Device
	Rpc              *common.RpcClient
	Pool             bool
	SubmitStr        chan string
	UseDevices       []string
}

//init GPU device
func (this *MinerRobot) InitDevice() {
	var typ = common.DevicesTypesForGPUMining
	if this.Cfg.OptionConfig.CPUMiner {
		common.MinerLoger.Warn("The parameter CPUMiner is deprecated !")
		cpuDevice := &cl.Device{}
		this.ClDevices = append(this.ClDevices, cpuDevice)
		return
	}
	needPlatform := ""
	if this.Cfg.OptionConfig.Cuda {
		needPlatform = "CUDA"
	}
	this.ClDevices = common.GetDevices(typ, needPlatform)
	if this.ClDevices == nil {
		common.MinerLoger.Info("Some GPU drivers error occurs! please check your GPU drivers.")
		return
	}
	this.UseDevices = []string{}
	if this.Cfg.OptionConfig.UseDevices != "" {
		this.UseDevices = strings.Split(this.Cfg.OptionConfig.UseDevices, ",")
	}
}
