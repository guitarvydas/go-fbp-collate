// Paul Morrison's "Collate" example on page 91 of "Flow-Based Programming, 2nd Edition"
//
//  +---+
//  | A |---+       +----------+
//  +---+   |       |          |       +--------+
//          +------>|in[0]     |       |        |
//                  |    C  out|------>|    D   |
//          +------>|in[1]     |       |        |
//  +---+   |       |          |       +--------+
//  | B |---+       +----------+
//  +---+
//

package main

import (
	"fmt"
	"github.com/guitarvydas/collate"
	"github.com/guitarvydas/ip"
	"github.com/guitarvydas/printer"
	"github.com/guitarvydas/readfile"
)

func main() {
	fmt.Println("Ports as explicit func arguments")
	msource := make(chan string)
	dsource := make(chan string)
	ctl := make(chan string)

	inCollate := arrayPort(2, 5) // array[2], bounds=5
	outCollate := make(chan ip.IP, 5)

	// wire up the components and start them
	go collate.Collate("Collate", ctl, inCollate, outCollate)
	go readfile.Read("Read Master", msource, inCollate[0])
	go readfile.Read("Read Details", dsource, inCollate[1])

	// send initialize's
	ctl <- "3,2,5"
	msource <- "mfile.txt"
	dsource <- "dfile.txt"

	printer.Print(outCollate)
}

func arrayPort(n, bound int) []chan ip.IP {
	a := make([]chan ip.IP, n)
	for i := 0; i < n; i++ {
		a[i] = make(chan ip.IP, bound)
	}
	return a
}
