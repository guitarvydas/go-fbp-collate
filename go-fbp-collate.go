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
	fmt.Println("Ports as maps of channels")
	c1 := make(chan string, 5)
	c2 := make(chan string, 5)
	c3 := make(chan string, 5)
	go collate(map[string]<-chan string{"in0": c1, "in1": c2}, map[string]chan<- string{"out": c3})
	go headerReader(nil, map[string]chan<- string{"out": c1})
	go recordReader(nil, map[string]chan<- string{"out": c2})
	resultPrinter(map[string]<-chan string{"in": c3}, nil)
}

func headerReader(in map[string]<-chan string, out map[string]chan<- string) {
	// fake reading of headers from a file, just generate some headers
	out["out"] <- "A"
	out["out"] <- "B"
	out["out"] <- "C"
}

func recordReader(in map[string]<-chan string, out map[string]chan<- string) {
	// fake reading records - first character is the key (header)
	for i := 0; i < 5; i++ {
		out["out"] <- fmt.Sprintf("A ++%v++", i)
	}
	for i := 15; i < 21; i++ {
		out["out"] <- fmt.Sprintf("B ++%v++", i)
	}
	for i := 25; i < 32; i++ {
		out["out"] <- fmt.Sprintf("C ++%v++", i)
	}
	out["out"] <- "EOF"
}

func resultPrinter(in map[string]<-chan string, out map[string]chan<- string) {
	for {
		merged := <-in["in"]
		if merged == "EOF" {
			break
		}
		fmt.Println(merged)
	}
}

func collate(in map[string]<-chan string, out map[string]chan<- string) {
	var hdr, rec string
	for {
		switch {
		case rec == "EOF":
			out["out"] <- rec
			break
		case hdr == "":
			hdr = <-in["in0"]
			out["out"] <- fmt.Sprintf("header %s", hdr)
		case rec != "" && rec[0] == hdr[0]:
			out["out"] <- fmt.Sprintf("record /%s/", rec[1:])
			rec = <-in["in1"]
		case rec != "" && rec[0] != hdr[0]:
			hdr = ""
		case rec == "":
			rec = <-in["in1"]
		}
	}
}
