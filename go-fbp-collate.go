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
	cac := make(chan string, 5) // channel between A and C
	cbc := make(chan string, 5) // channel between B and C
	ccd := make(chan string, 5) // channel betwen C and D
	go collate(cac, cbc, ccd)
	go headerReader(cac)
	go recordReader(cbc)
	resultPrinter(ccd)
}

func headerReader(out chan<- string) {
	// fake reading of headers from a file, just generate some headers
	out <- "A"
	out <- "B"
	out <- "C"
}

func recordReader(out chan<- string) {
	// fake reading records - first character is the key (header)
	for i := 0; i < 5; i++ {
		out <- fmt.Sprintf("A ++%v++", i)
	}
	for i := 15; i < 21; i++ {
		out <- fmt.Sprintf("B ++%v++", i)
	}
	for i := 25; i < 32; i++ {
		out <- fmt.Sprintf("C ++%v++", i)
	}
	out <- "EOF"
}

func resultPrinter(in <-chan string) {
	for {
		merged := <-in
		if merged == "EOF" {
			break
		}
		fmt.Println(merged)
	}
}

func collate(in0 <-chan string, in1 <-chan string, out chan<- string) {
	var hdr, rec string
	for {
		switch {
		case rec == "EOF":
			out <- rec
			break
		case hdr == "":
			hdr = <-in0
			out <- fmt.Sprintf("header %s", hdr)
		case rec != "" && rec[0] == hdr[0]:
			out <- fmt.Sprintf("record /%s/", rec[1:])
			rec = <-in1
		case rec != "" && rec[0] != hdr[0]:
			hdr = ""
		case rec == "":
			rec = <-in1
		}
	}
}
