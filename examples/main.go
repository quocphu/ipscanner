package main

import (
	"fmt"
	"time"

	"github.com/quocphu/ipscanner"
)

func main() {
	var threadNum uint32 = 20
	cOpen := make(chan string, threadNum)
	start := time.Now()
	// seed := [4]int{45, 79, 1, 0}
	// from := []byte{85, 125, 34, 0}
	// to := []byte{85, 125, 34, 255}

	// port := []int{8545}
	from := "85.125.34.205"
	to := "85.125.34.255"
	port := []int{8545}

	p := ipscanner.NewScanner(from, to, port, threadNum)
	du := time.Now().Sub(start)
	fmt.Println("Time ", du)
	fmt.Println("Total IP ", p.TotalIP)

	p.Scan(cOpen)
	var total uint32
	var totalFinishThread uint32
	start = time.Now()
	func() {
		for {
			select {
			case msg := <-cOpen:
				if msg == "END" {
					totalFinishThread++
					fmt.Println("Total end ", totalFinishThread)
				} else {
					if msg != "FAIL" {
						fmt.Println(msg)
					}
					total++
				}
				//fmt.Println("Total ", total)
				// if total == p.SumIP {
				// 	fmt.Println("Sum ip stop")
				// 	return
				// }
				if totalFinishThread == threadNum {
					return
				}
			default:
			}
		}
	}()
	fmt.Println("Done ", total)
	du = time.Now().Sub(start)
	fmt.Println("Time ", du)

}
