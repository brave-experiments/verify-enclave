.PHONY: test lint verify clean

binary = fetch-attestation

SCRIPT_PATH:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

test:
	go test -cover ./...

lint:
	golangci-lint run -E gofmt -E golint --exclude-use-default=false

verify:
	@go build -o $(binary) .
	@$(SCRIPT_PATH)/attest-enclave.sh $(CODE) $(ENCLAVE)

clean:
	rm -f $(binary) Dockerfile
