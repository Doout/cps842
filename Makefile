OUTPUT_ROOT := output

build:
	$(MAKE) build-linux
	$(MAKE) build-mac
	$(MAKE) build-win

.PHONY: all test clean


build-linux:
	dep ensure
	GOOS=linux  CGO_ENABLED=0 go build -o $(OUTPUT_ROOT)/cps842-linux .

.PHONY: test clean

build-mac:
	dep ensure
	GOOS=darwin  CGO_ENABLED=0 go build -o $(OUTPUT_ROOT)/cps842-mac .

.PHONY: test clean

build-win:
	dep ensure
	GOOS=windows  CGO_ENABLED=0 go build -o $(OUTPUT_ROOT)/cps842-win .

.PHONY: test clean

