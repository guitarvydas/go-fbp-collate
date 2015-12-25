# set GOPATH and GOBIN as per https://golang.org/doc/install

all: $(GOBIN)/go-fbp-collate  $(GOBIN)/fanin

$(GOBIN)/go-fbp-collate : go-fbp-collate.go
	go install ./go-fbp-collate.go

$(GOBINP)/fanin : fanin.go
	go install ./fanin.go
