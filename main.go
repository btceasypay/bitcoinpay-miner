// Copyright (c) 2019 The bitcoinpay developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
package main

import (
	"github.com/btceasypay/bitcoinpay-miner/common"
	"github.com/btceasypay/bitcoinpay-miner/core"
	"github.com/btceasypay/bitcoinpay-miner/symbols/bitcoinpay"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"time"
)

var robotminer core.Robot

//init the config file
func init() {
	cfg, _, err := common.LoadConfig()
	if err != nil {
		log.Fatal("[error] config error,please check it.[", err, "]")
		return
	}
	//init miner robot
	robotminer = GetRobot(cfg)
}

func main() {
	// Use all processor cores.
	runtime.GOMAXPROCS(runtime.NumCPU())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		common.MinerLoger.Info("Got Control+C, exiting...")
		os.Exit(0)
	}()
	if robotminer == nil {
		common.MinerLoger.Error("[error] Please check the coin in the README.md! if this coin is supported, use -S to set")
		return
	}
	robotminer.Run()
}

//get current coin miner
func GetRobot(cfg *common.GlobalConfig) core.Robot {
	switch strings.ToUpper(cfg.NecessaryConfig.Symbol) {
	case core.SYMBOL_PMEER:
		r := &bitcoinpay.BitcoinpayRobot{}
		r.Cfg = cfg
		r.Started = uint32(time.Now().Unix())
		r.Rpc = &common.RpcClient{Cfg: cfg}
		r.SubmitStr = make(chan string)
		return r
	default:

	}
	return nil
}
