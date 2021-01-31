# Building from source

## Table of Contents

* [Prequisite](#prequisite)
    * [Common](#common)
    * [Linux](#linux)
        * [Ubuntu](#ubuntu)
        * [Centos](#centos)
    * [macOS](#macos)
    * [Windows](#windows)
* [Build](#build)
    * [Windows-additional step](#windows-additional-step)
    * [Linux-additional step](#Linux-additional step)
    
### Prequisite

### Common

1. [Git](https://git-scm.com/downloads) 
2. [Go](https://golang.org/dl/) version >= 1.12

### Linux


#### Ubuntu

```bash
$ sudo apt-get install beignet-dev nvidia-cuda-dev nvidia-cuda-toolkit
```
        
#### Centos 

```bash
$ sudo yum install opencl-headers
$ sudo yum install ocl-icd
$ sudo ln -s /usr/lib64/libOpenCL.so.1 /usr/lib/libOpenCL.so
```  
### MacOS

### Windows

Install [**Build Tools for Visual Studio**](https://visualstudio.microsoft.com/thank-you-downloading-visual-studio/?sku=BuildTools&rel=16)
    
## Build 

### 1. Get Source code

```bash
$ git clone git@github.com:Bitcoinpay/bitcoinpay-miner.git
```

### 2. Build the cudacuckaroom library 

[Build Step](lib/cuda/cuckaroom/README.md)

### 3. Build bitcoinpay-miner  

```bash
$ cp lib/cuda/cuckaroom/cudacuckoo.dll lib/cuda/
$ cp lib/cuda/cuckaroom/cudacuckoo.so lib/cuda/
//# mac
$ go build --tags opencl
//# linux apt install musl-tools g++ -y
$ CGO_ENABLED=1 CC=musl-gcc CXX=g++ GOOS=linux go build -tags cuda -o linux-miner main.go
//# windows 
$ CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -tags cuda -o win-miner.exe main.go
```

### 4. Verify Build OK

```bash
$ ./bitcoinpay-miner --version
```
