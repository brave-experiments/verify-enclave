module github.com/brave-experiments/attest-enclave

go 1.17

require github.com/hf/nitrite v0.0.0-20210106181132-25252415087f

require (
	github.com/fxamacker/cbor/v2 v2.3.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	go.mozilla.org/cose v0.0.0-20200930124131-25dc96df8228 // indirect
)

replace github.com/hf/nitrite => ../nitrite
