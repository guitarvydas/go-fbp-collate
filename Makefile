# set GOPATH and GOBIN as per https://golang.org/doc/install

../../bin/go-fbp-collate : go-fbp-collate.go
	go install ./go-fbp-collate.go
