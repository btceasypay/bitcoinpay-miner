/**
Bitcoinpay
james
*/

package common

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	bitcoinpay "github.com/btceasypay/bitcoinpay/common/hash"
	"github.com/btceasypay/bitcoinpay/core/types/pow"
	"log"
	"math"
	"math/big"
	"math/rand"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
)

func SliceContains(s []uint64, e uint64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
func SliceRemove(s []uint64, e uint64) []uint64 {
	for i, a := range s {
		if a == e {
			return append(s[:i], s[i+1:]...)
		}
	}

	return s
}

func BlockBitsToTarget(bits string, width int) []byte {
	nbits, err := hex.DecodeString(bits[0:2])
	if err != nil {
		fmt.Println("error", err.Error())
	}
	shift := nbits[0] - 3
	value, _ := hex.DecodeString(bits[2:])
	target0 := make([]byte, int(shift))
	tmp := string(value) + string(target0)
	target1 := []byte(tmp)
	if len(target1) < width {
		head := make([]byte, width-len(target1))
		target := string(head) + string(target1)
		return []byte(target)
	}
	return target1
}

func Int2varinthex(x int64) string {
	if x < 0xfd {
		return fmt.Sprintf("%02x", x)
	} else if x < 0xffff {
		return "fd" + Int2lehex(x, 2)
	} else if x < 0xffffffff {
		return "fe" + Int2lehex(x, 4)
	} else {
		return "ff" + Int2lehex(x, 8)
	}
}
func Int2lehex(x int64, width int) string {
	if x <= 0 {
		return fmt.Sprintf("%016x", x)
	}
	bs := make([]byte, width)
	switch width {
	case 2:
		binary.LittleEndian.PutUint16(bs, uint16(x))
	case 4:
		binary.LittleEndian.PutUint32(bs, uint32(x))
	case 8:
		binary.LittleEndian.PutUint64(bs, uint64(x))
	}
	return hex.EncodeToString(bs)
}

// Reverse reverses a byte array.
func Reverse(src []byte) []byte {
	dst := make([]byte, len(src))
	for i := len(src); i > 0; i-- {
		dst[len(src)-i] = src[i-1]
	}
	return dst
}

// FormatHashRate sets the units properly when displaying a hashrate.
func FormatHashRate(h float64, unit string) string {
	if h > 1000000000000 {
		return fmt.Sprintf("%.4fT%s", h/1000000000000, unit)
	} else if h > 1000000000 {
		return fmt.Sprintf("%.4fG%s", h/1000000000, unit)
	} else if h > 1000000 {
		return fmt.Sprintf("%.4fM%s", h/1000000, unit)
	} else if h > 1000 {
		return fmt.Sprintf("%.4fk%s", h/1000, unit)
	} else if h == 0 {
		return fmt.Sprintf("0%s", unit)
	} else if h > 0 {
		return fmt.Sprintf("%.4f%s", h, unit)
	}

	return fmt.Sprintf("%.4f T%s", h, unit)
}

func ReverseByWidth(s []byte, width int) []byte {
	newS := make([]byte, len(s))
	for i := 0; i < (len(s) / width); i += 1 {
		j := i * width
		copy(newS[len(s)-j-width:len(s)-j], s[j:j+width])
	}
	return newS
}

func DiffToTarget(diff float64, powLimit *big.Int, powType pow.PowType) (*big.Int, error) {
	if diff <= 0 {
		return nil, fmt.Errorf("invalid pool difficulty %v (0 or less than "+
			"zero passed)", diff)
	}

	// Round down in the case of a non-integer diff since we only support
	// ints (unless diff < 1 since we don't allow 0)..
	if diff < 1 {
		diff = 1
	} else {
		diff = math.Floor(diff)
	}
	divisor := new(big.Int).SetInt64(int64(diff))
	max := powLimit
	target := new(big.Int)
	if powType == pow.BLAKE2BD {
		target.Div(max, divisor)
	} else {
		target.Div(divisor, max)
	}

	return target, nil
}

// Uint32EndiannessSwap swaps the endianness of a uint32.
func Uint32EndiannessSwap(v uint32) uint32 {
	return (v&0x000000FF)<<24 | (v&0x0000FF00)<<8 |
		(v&0x00FF0000)>>8 | (v&0xFF000000)>>24
}

// RolloverExtraNonce rolls over the extraNonce if it goes over 0x00FFFFFF many
// hashes, since the first byte is reserved for the ID.
func RolloverExtraNonce(v *uint32) {
	if *v&0x00FFFFFF == 0x00FFFFFF {
		*v = *v & 0xFF000000
	} else {
		*v++
	}
}

func ConvertHashToString(hash bitcoinpay.Hash) string {
	newB := make([]byte, 32)
	copy(newB[:], hash[:])
	return hex.EncodeToString(newB)
}

// appDataDir returns an operating system specific directory to be used for
// storing application data for an application.  See AppDataDir for more
// details.  This unexported version takes an operating system argument
// primarily to enable the testing package to properly test the function by
// forcing an operating system that is not the currently one.
func appDataDir(goos, appName string, roaming bool) string {
	if appName == "" || appName == "." {
		return "."
	}

	// The caller really shouldn't prepend the appName with a period, but
	// if they do, handle it gracefully by stripping it.
	appName = strings.TrimPrefix(appName, ".")
	appNameUpper := string(unicode.ToUpper(rune(appName[0]))) + appName[1:]
	appNameLower := string(unicode.ToLower(rune(appName[0]))) + appName[1:]

	// Get the OS specific home directory via the Go standard lib.
	var homeDir string
	usr, err := user.Current()
	if err == nil {
		homeDir = usr.HomeDir
	}

	// Fall back to standard HOME environment variable that works
	// for most POSIX OSes if the directory from the Go standard
	// lib failed.
	if err != nil || homeDir == "" {
		homeDir = os.Getenv("HOME")
	}

	switch goos {
	// Attempt to use the LOCALAPPDATA or APPDATA environment variable on
	// Windows.
	case "windows":
		// Windows XP and before didn't have a LOCALAPPDATA, so fallback
		// to regular APPDATA when LOCALAPPDATA is not set.
		appData := os.Getenv("LOCALAPPDATA")
		if roaming || appData == "" {
			appData = os.Getenv("APPDATA")
		}

		if appData != "" {
			return filepath.Join(appData, appNameUpper)
		}

	case "darwin":
		if homeDir != "" {
			return filepath.Join(homeDir, "Library",
				"Application Support", appNameUpper)
		}

	case "plan9":
		if homeDir != "" {
			return filepath.Join(homeDir, appNameLower)
		}

	default:
		if homeDir != "" {
			return filepath.Join(homeDir, "."+appNameLower)
		}
	}

	// Fall back to the current directory if all else fails.
	return "."
}

// AppDataDir returns an operating system specific directory to be used for
// storing application data for an application.
//
// The appName parameter is the name of the application the data directory is
// being requested for.  This function will prepend a period to the appName for
// POSIX style operating systems since that is standard practice.  An empty
// appName or one with a single dot is treated as requesting the current
// directory so only "." will be returned.  Further, the first character
// of appName will be made lowercase for POSIX style operating systems and
// uppercase for Mac and Windows since that is standard practice.
//
// The roaming parameter only applies to Windows where it specifies the roaming
// application data profile (%APPDATA%) should be used instead of the local one
// (%LOCALAPPDATA%) that is used by default.
//
// Example results:
//  dir := AppDataDir("myapp", false)
//   POSIX (Linux/BSD): ~/.myapp
//   Mac OS: $HOME/Library/Application Support/Myapp
//   Windows: %LOCALAPPDATA%\Myapp
//   Plan 9: $home/myapp
func AppDataDir(appName string, roaming bool) string {
	return appDataDir(runtime.GOOS, appName, roaming)
}

func Target2BlockBits(target string) []byte {
	// 8
	d, _ := hex.DecodeString(target[0:16])
	return Reverse(d)
}

func HexMustDecode(hexStr string) []byte {
	b, err := hex.DecodeString(hexStr)
	if err != nil {

		panic(err)
	}
	return b
}

func GetCurrentDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalln(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func RandUint64() uint64 {
	return rand.Uint64()
}

func RandUint32() (uint32, error) {
	var b [4]byte
	if _, err := rand.Read(b[:]); err != nil {
		return 0, err
	}
	return uint32(binary.LittleEndian.Uint32(b[:])), nil
}

func InArray(val interface{}, arr interface{}) bool {
	switch arr.(type) {
	case []string:
		for _, v := range arr.([]string) {
			if v == val {
				return true
			}
		}
	case []int:
		for _, v := range arr.([]int) {
			if v == val {
				return true
			}
		}
	}

	return false
}

func Timeout(timeout time.Duration, runFunc func()) bool {
	var wg = new(sync.WaitGroup)
	c := make(chan interface{})
	wg.Add(1)
	go func() {
		defer close(c)
		wg.Wait()
	}()
	go func() {
		runFunc()
		c <- nil
		wg.Done()
	}()
	select {
	case <-c:
		return false
	case <-time.After(timeout):
		return true
	}
}

func TimeoutRun(timeout time.Duration, runFunc, afterFun func()) bool {
	var wg = new(sync.WaitGroup)
	c := make(chan interface{})
	wg.Add(1)
	go func() {
		defer close(c)
		wg.Wait()
	}()
	go func() {
		runFunc()
		c <- nil
		wg.Done()
	}()
	select {
	case <-c:
		return false
	case <-time.After(timeout):
		Timeout(1, func() {
			afterFun()
		})
		return true
	}
}
func GetNeedHashTimesByTarget(target string) *big.Int {
	times := big.NewInt(1)
	for i := 0; i < len(target)-1; i++ {
		tmp := target[i : i+1]
		if tmp == "0" {
			times.Lsh(times, 4)
		} else {
			n, _ := strconv.ParseUint(tmp, 16, 32)
			if n <= 1 {
				n = 1
			}
			n1 := int64(16 / n)
			times.Mul(times, big.NewInt(n1))
			break
		}
	}
	return times
}
