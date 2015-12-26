// Paul Morrison's "Collate" example on page 91 of "Flow-Based Programming, 2nd Edition"
//
//  +---+
//  | A |---+       +----------+
//  +---+   |       |          |       +--------+
//          +------>|0         |       |        |
//                  |    C  out|------>|    D   |
//          +------>|1         |       |        |
//  +---+   |       |          |       +--------+
//  | B |---+       +----------+
//  +---+
//
// Component A produces "headers".  Component B produces data keyed by headers (sorted).  Component C
// Collates the output from A and B and sends the merged output to D.  C reads a header from A, then keeps
// reading data records from B until the key changes.  Then, C reads the next header from A and repeats.
//
// The implementation below shows only the "happy" case.  Edge cases, e.g. what happens when B skips over
// a header, are not implemented below, but it should be clear how such logic could be added.
//
// This implementation explicitly creates Go channels for A to C, B to C and C to D, then passes theses
// channels as parameters to the various components, thereby "wiring up" the diagram.
//

package main

import "fmt"

func main() {
	fmt.Println("Ports as arrays of channels")
	c1 := make(chan string, 5)
	c2 := make(chan string, 5)
	c3 := make(chan string, 5)
	go collate([]<-chan string{c1, c2}, []chan<- string{c3})
	go headerReader(nil, []chan<- string{c1})
	go recordReader(nil, []chan<- string{c2})
	resultPrinter([]<-chan string{c3}, nil)
}

func headerReader(in []<-chan string, out []chan<- string) {
	// fake reading of headers from a file, just generate some headers
	out[0] <- "A"
	out[0] <- "B"
	out[0] <- "C"
}

func recordReader(in []<-chan string, out []chan<- string) {
	// fake reading records - first character is the key (header)
	for i := 0; i < 5; i++ {
		out[0] <- fmt.Sprintf("A ++%v++", i)
	}
	for i := 15; i < 21; i++ {
		out[0] <- fmt.Sprintf("B ++%v++", i)
	}
	for i := 25; i < 32; i++ {
		out[0] <- fmt.Sprintf("C ++%v++", i)
	}
	out[0] <- "EOF"
}

func resultPrinter(in []<-chan string, out []chan<- string) {
	for {
		merged := <-in[0]
		if merged == "EOF" {
			break
		}
		fmt.Println(merged)
	}
}

func collate(in []<-chan string, out []chan<- string) {
	var hdr, rec string
	for {
		switch {
		case rec == "EOF":
			out[0] <- rec
			break
		case hdr == "":
			hdr = <-in[0]
			out[0] <- fmt.Sprintf("header %s", hdr)
		case rec != "" && rec[0] == hdr[0]:
			out[0] <- fmt.Sprintf("record /%s/", rec[1:])
			rec = <-in[1]
		case rec != "" && rec[0] != hdr[0]:
			hdr = ""
		case rec == "":
			rec = <-in[1]
		}
	}
}
