// Copyright (c) 2019 The bitcoinpay developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
package bitcoinpay

import (
	"fmt"
	"github.com/btceasypay/bitcoinpay-miner/cl"
	"github.com/btceasypay/bitcoinpay-miner/common"
	"github.com/btceasypay/bitcoinpay-miner/core"
	"github.com/btceasypay/bitcoinpay-miner/stats_server"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	POW_CUCKROO   = "cuckaroo"
	POW_CUCKROOM  = "cuckaroom"
	POW_CUCKROO29 = "cuckaroo29"
)

type BitcoinpayRobot struct {
	core.MinerRobot
	Work                 BitcoinpayWork
	Devices              []core.BaseDevice
	Stratu               *BitcoinpayStratum
	AllTransactionsCount int64
}

func (this *BitcoinpayRobot) GetPow(i int, device *cl.Device) core.BaseDevice {
	switch this.Cfg.NecessaryConfig.Pow {
	case POW_CUCKROOM:
		if !this.Cfg.OptionConfig.Cuda {
			deviceMiner := &Cuckaroo{}
			deviceMiner.MiningType = "cuckaroom"
			deviceMiner.Init(i, device, this.Pool, this.Quit, this.Cfg)
			this.Devices = append(this.Devices, deviceMiner)
			return deviceMiner
		} else {
			deviceMiner := &CudaCuckaroom{}
			deviceMiner.MiningType = "cuckaroom"
			deviceMiner.Init(i, device, this.Pool, this.Quit, this.Cfg)
			this.Devices = append(this.Devices, deviceMiner)
			return deviceMiner
		}
	case POW_CUCKROO:
		if !this.Cfg.OptionConfig.Cuda {
			deviceMiner := &Cuckaroo{}
			deviceMiner.MiningType = "cuckaroo"
			deviceMiner.Init(i, device, this.Pool, this.Quit, this.Cfg)
			this.Devices = append(this.Devices, deviceMiner)
			return deviceMiner
		} else {
			deviceMiner := &CudaCuckaroo{}
			deviceMiner.MiningType = "cuckaroo"
			deviceMiner.Init(i, device, this.Pool, this.Quit, this.Cfg)
			this.Devices = append(this.Devices, deviceMiner)
			return deviceMiner
		}

	default:
		log.Fatalln(this.Cfg.NecessaryConfig.Pow, " pow has not exist!")
	}
	return nil
}

func (this *BitcoinpayRobot) InitDevice() {
	this.MinerRobot.InitDevice()
	for i, device := range this.ClDevices {
		deviceMiner := this.GetPow(i, device)
		if deviceMiner == nil {
			return
		}
	}
}

// runing
func (this *BitcoinpayRobot) Run() {
	this.Wg = &sync.WaitGroup{}
	this.InitDevice()
	//mining service
	connectName := "solo"
	this.Pool = false
	if this.Cfg.PoolConfig.Pool != "" { //is pool mode
		connectName = "pool"
		this.Stratu = &BitcoinpayStratum{}
		_ = this.Stratu.StratumConn(this.Cfg)
		this.Wg.Add(1)
		go func() {
			defer this.Wg.Done()
			this.Stratu.HandleReply()
		}()
		this.Pool = true
	}
	common.MinerLoger.Info(fmt.Sprintf("%s miner start", connectName))
	this.Work = BitcoinpayWork{}
	this.Work.Cfg = this.Cfg
	this.Work.Rpc = this.Rpc
	this.Work.stra = this.Stratu
	// Device Miner
	for _, dev := range this.Devices {
		dev.SetIsValid(true)
		if len(this.UseDevices) > 0 && !common.InArray(strconv.Itoa(dev.GetMinerId()), this.UseDevices) {
			dev.SetIsValid(false)
			continue
		}
		dev.SetPool(this.Pool)
		dev.InitDevice()
		this.Wg.Add(1)
		go dev.Mine(this.Wg)
		this.Wg.Add(1)
		go dev.Status(this.Wg)
	}
	//refresh work
	this.Wg.Add(1)
	go func() {
		defer this.Wg.Done()
		this.ListenWork()
	}()
	//submit work
	this.Wg.Add(1)
	go func() {
		defer this.Wg.Done()
		this.SubmitWork()
	}()
	//submit status
	this.Wg.Add(1)
	go func() {
		defer this.Wg.Done()
		this.Status()
	}()

	//http server stats
	if this.Cfg.OptionConfig.StatsServer != "" {
		this.Wg.Add(1)
		go func() {
			defer this.Wg.Done()
			stats_server.HandleRouter(this.Cfg, this.Devices)
		}()
	}

	this.Wg.Wait()
}

// ListenWork
func (this *BitcoinpayRobot) ListenWork() {
	common.MinerLoger.Info("listen new work server")
	t := time.NewTicker(time.Second * time.Duration(this.Cfg.OptionConfig.TaskInterval))
	isFirst := true
	defer t.Stop()
	r := false
	for {
		select {
		case <-this.Quit:
			return
		case <-t.C:
			r = false
			if this.Pool {
				r = this.Work.PoolGet() // get new work
			} else {
				r = this.Work.Get() // get new work
			}
			if r {
				common.MinerLoger.Debug("new task started")
				validDeviceCount := 0
				for _, dev := range this.Devices {
					if !dev.GetIsValid() {
						continue
					}
					dev.SetForceUpdate(false)
					validDeviceCount++
					newWork := this.Work.CopyNew()
					dev.SetNewWork(&newWork)
				}
				if validDeviceCount <= 0 {
					common.MinerLoger.Error("There is no valid device to mining,please check your config!")
					os.Exit(1)
				}
				if isFirst {
					isFirst = false
				}
			} else if this.Work.ForceUpdate {
				for _, dev := range this.Devices {
					common.MinerLoger.Debug("task stopped by force")
					dev.SetNewWork(&this.Work)
					dev.SetForceUpdate(true)
				}
			}
		}
	}
}

// ListenWork
func (this *BitcoinpayRobot) SubmitWork() {
	common.MinerLoger.Info("listen submit block server")
	this.Wg.Add(1)
	go func() {
		defer this.Wg.Done()
		str := ""
		var logContent string
		var count int
		var arr []string
		for {
			common.MinerLoger.Debug("===============================Listen Submit=====================")
			select {
			case <-this.Quit:
				return
			case str = <-this.SubmitStr:
				if str == "" {
					atomic.AddUint64(&this.StaleShares, 1)
					continue
				}
				var err error
				var height, txCount, block string
				if this.Pool {
					arr = strings.Split(str, "-")
					block = arr[0]
					err = this.Work.PoolSubmit(str)
				} else {
					//solo miner
					arr = strings.Split(str, "-")
					txCount = arr[1]
					height = arr[2]
					block = arr[0]
					err = this.Work.Submit(block)
				}
				if err != nil {
					if err != ErrSameWork || err == ErrSameWork {
						if err == ErrStratumStaleWork {
							atomic.AddUint64(&this.StaleShares, 1)
						} else {
							atomic.AddUint64(&this.InvalidShares, 1)
						}
					}
				} else {
					atomic.AddUint64(&this.ValidShares, 1)
					if !this.Pool {
						count, _ = strconv.Atoi(txCount)
						this.AllTransactionsCount += int64(count)
						logContent = fmt.Sprintf("receive block, block height = %s,Including %s transactions; Received Total transactions = %d\n",
							height, txCount, this.AllTransactionsCount)
						common.MinerLoger.Info(logContent)
					}
				}
			}
		}

	}()
	for _, dev := range this.Devices {
		go dev.SubmitShare(this.SubmitStr)
	}
}

// stats the submit result
func (this *BitcoinpayRobot) Status() {
	var valid, rejected, staleShares uint64
	for {
		select {
		case <-this.Quit:
			return
		default:
			if this.Work.stra == nil && this.Work.Block == nil {
				common.Usleep(20 * 1000)
				continue
			}
			valid = atomic.LoadUint64(&this.ValidShares)
			rejected = atomic.LoadUint64(&this.InvalidShares)
			staleShares = atomic.LoadUint64(&this.StaleShares)
			if this.Pool {
				valid = atomic.LoadUint64(&this.Stratu.ValidShares)
				rejected = atomic.LoadUint64(&this.Stratu.InvalidShares)
				staleShares = atomic.LoadUint64(&this.Stratu.StaleShares)
			}
			this.Cfg.OptionConfig.Accept = int(valid)
			this.Cfg.OptionConfig.Reject = int(rejected)
			this.Cfg.OptionConfig.Stale = int(staleShares)
			total := valid + rejected + staleShares
			common.MinerLoger.Info(fmt.Sprintf("Global stats: Accepted: %v,Stale: %v, Rejected: %v, Total: %v",
				valid,
				staleShares,
				rejected,
				total,
			))
			common.Usleep(20 * 1000)
		}
	}
}
