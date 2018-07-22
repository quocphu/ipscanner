// Package ipscanner implements methods to scan provided ip & port are open or not
package ipscanner

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// IPScanner struct
type IPScanner struct {
	threadNum        uint32
	Ports            []int
	From             []byte
	To               []byte
	TotalIP          uint32
	ipRangeOfAThread []uint32
	timeOut          time.Duration
}

// NewScanner creates new object of IPScanner
//
// @from: Starting ip
// @to: Ending ip
// @ports: Array of port
// @numberOfThread: Number go route
func NewScanner(from string, to string, ports []int, numberOfThread uint32) *IPScanner {
	s := &IPScanner{
		threadNum: numberOfThread,
		Ports:     ports,
		From:      stringToArray(from),
		To:        stringToArray(to),
	}
	s.separateIpForAThread()
	return s
}

// Scan scan all ip
func (sc *IPScanner) Scan(c chan string) {
	ipRange := sc.ipRangeOfAThread

	for i := 0; i < len(ipRange)-1; i++ {
		from := ipRange[i]
		to := ipRange[i+1]

		go sc.directScan(from, to, to-from, c)
	}
}

func (sc *IPScanner) directScan(from uint32, to uint32, totalIP uint32, c chan string) {
	var count uint32
	for i := from; i < to; i++ {
		for p := 0; p < len(sc.Ports); p++ {
			count++
			ipByte := make([]byte, 4)
			binary.BigEndian.PutUint32(ipByte, uint32(i))
			ip := fmt.Sprintf("%d.%d.%d.%d:%d", ipByte[0], ipByte[1], ipByte[2], ipByte[3], sc.Ports[p])
			sc.tryConnect(ip, c, count == totalIP)
		}
	}

}
func (sc *IPScanner) separateIpForAThread() {
	fromNumber := binary.BigEndian.Uint32(sc.From)
	toNumer := binary.BigEndian.Uint32(sc.To)

	totalIP := toNumer - fromNumber
	if totalIP < uint32(sc.threadNum) {
		panic(fmt.Sprintf("Total IP must be greater than total thread number: total ip = %d, thread = %d", totalIP, sc.threadNum))
	}

	ipPerThread := totalIP / uint32(sc.threadNum)
	modIpNumber := int32(totalIP % sc.threadNum)

	var ipOfThread []uint32
	ipOfThread = append(ipOfThread, fromNumber)

	var i uint32 = 1
	for ; i < sc.threadNum; i++ {
		stop := ipOfThread[i-1] + ipPerThread

		if modIpNumber >= 0 {
			modIpNumber--
			stop++
		}
		ipOfThread = append(ipOfThread, stop)
	}

	ipOfThread = append(ipOfThread, totalIP+fromNumber+1)
	sc.ipRangeOfAThread = ipOfThread
	sc.TotalIP = totalIP
}

func (sc *IPScanner) tryConnect(ip string, c chan string, isLast bool) {
	conn, err := net.DialTimeout("tcp", ip, sc.timeOut*time.Second)
	fmt.Println(ip)
	if err != nil {
		c <- "FAIL" // + ip

	} else {
		c <- ip
	}

	if isLast {
		c <- "END"
	}

	if conn != nil {
		return
	}
}

func stringToArray(ip string) []byte {
	var rs []byte
	arrIp := strings.Split(ip, ".")
	for i := 0; i < len(arrIp); i++ {
		if s, err := strconv.ParseUint(arrIp[i], 10, 32); err == nil {
			rs = append(rs, byte(s))
		}
	}
	return rs
}
