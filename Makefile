.PHONY: test lint verify clean

binary = fetch-attestation

test:
	go test -cover ./...

lint:
	golangci-lint run -E gofmt -E golint --exclude-use-default=false

verify:
	@go build -o $(binary) .
	@./attest-enclave.sh $(CODE)

clean:
	rm -f $(binary)
