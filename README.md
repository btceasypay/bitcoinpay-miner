# Bitcoinpay Miner

[![Build Status](https://travis-ci.com/Bitcoinpay/Bitcoinpay-miner.svg?token=n9AoZUDqAJmhesf4MYUd&branch=master)](https://travis-ci.com/Bitcoinpay/Bitcoinpay-miner)

> The official GPU miner of the Bitcoinpay network  

**Bitcoinpay-miner** is an GPU miner for the Bitcoinpay netowrk. It's the official reference implement maintained by the Bitcoinpay team.
Currently it support 3 Bitcoinpay POW algorithms including Cuckaroo, Cuckatoo and Blake2bd.

## Table of Contents
* [Install](#install)
* [Usage](#usage)
   - [Run with config file](#run-with-config-file)
   - [Run by Command line options](#command-line-usage)
* [Build](#build)
   - [Building from source](#building-from-source)
* [Tutorial](#tutorial)    
* [FAQ](#faq)


## Install

[![Releases](https://img.shields.io/github/downloads/Bitcoinpay/Bitcoinpay-miner/total.svg)][Releases]

Standalone installation archive for *Linux*, *macOS* and *Windows* are provided in
the [Releases] section. 
Please download an archive for your operating system and unpack the content to a place
accessible from command line. 

| Builds | Release | Date |
| ------ | ------- | ---- |
| Last   | [![GitHub release](https://img.shields.io/github/release/Bitcoinpay/Bitcoinpay-miner/all.svg)][Releases] | [![GitHub Release Date](https://img.shields.io/github/release-date-pre/Bitcoinpay/Bitcoinpay-miner.svg)][Releases] |
| Stable | [![GitHub release](https://img.shields.io/github/release/Bitcoinpay/Bitcoinpay-miner.svg)][latest] | [![GitHub Release Date](https://img.shields.io/github/release-date/Bitcoinpay/Bitcoinpay-miner.svg)][latest] |

## Usage

### Run with config file 
1. go to your 
2. create a new config file by copying from the example config file. 
```bash
$ cp example.solo.conf solo.conf
```
3. edit the config file which your create, you might need to change the `mineraddress`. 
you need to create a Bitcoinpay address if you don't have it. Please see [FAQ](#FAQ)  
4. run miner with the config file

```bash
$ ./bitcoinpay-miner -C solo.conf
```

### Command line usage

The Bitcoinpay-miner is a command line program. This means you can also launch it by provided valid command line options. For a full list of available command optinos, please run:

```bash
$ ./bitcoinpay-miner --help 
Usage:
  bitcoinpay-miner [OPTIONS]

Debug Command:
  -l, --listdevices             List number of devices.
  -v, --version                 show the version of miner

The Config File Options:
  -C, --configfile=             Path to configuration file
      --minerlog=               Write miner log file

The Necessary Config Options:
  -P, --pow=                    blake2bd|cuckaroo|cuckatoo (cuckaroo)
  -S, --symbol=                 Symbol (PMEER)
  -N, --network=                network privnet|testnet|mainnet (testnet)

The Solo Config Option:
  -M, --mineraddress=           Miner Address
  -s, --rpcserver=              RPC server to connect to (127.0.0.1)
  -u, --rpcuser=                RPC username
  -p, --rpcpass=                RPC password
      --randstr=                Rand String,Your Unique Marking. (Come from Bitcoinpay!)
      --notls                   Do not verify tls certificates (true)
      --rpccert=                RPC server certificate chain for validation

The pool Config Option:
  -o, --pool=                   Pool to connect to (e.g.stratum+tcp://pool:port)
  -m, --pooluser=               Pool username
  -n, --poolpass=               Pool password

The Optional Config Option:
      --cpuminer                CPUMiner (false)
      --proxy=                  Connect via SOCKS5 proxy (eg. 127.0.0.1:9050)
      --proxyuser=              Username for proxy server
      --proxypass=              Password for proxy server
      --trimmerTimes=           the cuckaroo trimmer times (15)
      --intensity=              Intensities (the work size is 2^intensity) per device. Single global value or a comma separated
                                list. (24)
      --worksize=               The explicitly declared sizes of the work to do per device (overrides intensity). Single global
                                value or a comma separated list. (256)
      --timeout=                rpc timeout. (60)
      --use_devices=            all gpu devices,you can use ./bitcoinpay-miner -l to see. examples:0,1 use the #0 device and #1
                                device
      --max_tx_count=           max pack tx count (1000)
      --max_sig_count=          max sign tx count (4000)
      --log_level=              info|debug|error|warn|trace (info)
      --stats_server=           stats web server (127.0.0.1:1235)
      --edge_bits=              edge bits (24)
      --local_size=             local size (4096)
      --group_size=             work group size (256)
      --cuda                    is cuda (false)
      --task_interval=          get blocktemplate interval (2)
      --task_force_stop         force stop the current task when miner fail to get blocktemplate from the bitcoinpay full node.
                                (true)
      --mining_sync_mode        force stop the current task when new task come. (true)
      --force_solo              force solo mining (false)
      --big_graph_start_height= big graph start main height, how many days later,the r29 will be the main pow (45)
      --expand=                 expand enum 0,1,2 (0)
      --ntrims=                 trim times  (50)
      --genablocks=             genablocks (4096)
      --genatpb=                genatpb (256)
      --genbtpb=                genbtpb (256)
      --trimtpb=                genbtpb (64)
      --tailtpb=                tailtpb (1024)
      --recoverblocks=          recoverblocks (1024)
      --recovertpb=             recovertpb (1024)

Help Options:
  -h, --help                    Show this help message
 
```
Please see [Bitcoinpay-Miner User References](https://Bitcoinpay.github.io/docs/en/reference/Bitcoinpay-miner/) for more details

## Build
### Building from source
See [BUILD.md](BUILD.md) for build/compilation details.

## Tutorial

## FAQ

### How to create Bitcoinpay adderss
There are several ways to create a Bitcoinpay address. you can use [bx][Bx] command , [Bitcoinpay-wallet][Bitcoinpay-wallet], etc.
The most easy way to download the [Bitcoinpay wallet][Bitcoinpay.io], which provide a more user friendly GUI to create your address/wallet step by step. 

### Which POW algorithm I should choose to mine ?
Bitcoinpay test network support mixing minning, which means your can choice from `Cuckaroo`, `Cuckatoo` and `Blake2bd` anyone you like. 
But the start difficulty targets are quite different. For the most case you might use `Cuckaroo` as a safe choice at the beginning. 

### Where I can find more documentation ? 
Please find more documentation from the [Bitcoinpay doc site at https://bitcoinpay.github.io](https://bitcoinpay.github.io/docs/en/reference/bitcoinpay-miner/)

[Releases]: https://github.com/btceasypay/bitcoinpay-miner/releases
[Latest]: https://github.com/btceasypay/bitcoinpay-miner/releases/latest
[Bx]: https://bitcoinpay.github.io/docs/en/reference/bxtools/
[Bitcoinpay-wallet]: https://github.com/btceasypay/bitcoinpay-wallet