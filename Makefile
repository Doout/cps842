build: 
	dep ensure
	GOOS=linux  CGO_ENABLED=0 go build .

.PHONY: all test clean
