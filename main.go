package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/hf/nitrite"
)

const (
	nonceSize = 20 // Nonce size in bytes.
)

func fetchAttestationDocument(nonce, attestationEndpoint string) ([]byte, error) {

	attestationEndpoint += fmt.Sprintf("?nonce=%s", nonce)
	resp, err := http.Get(attestationEndpoint)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Failed to close HTTP response body: %s", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("attestation handler returned HTTP %d", resp.StatusCode)
	}

	b64Doc, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// Decode Base64-encoded attestation document.
	rawDoc := make([]byte, base64.StdEncoding.DecodedLen(len(b64Doc)))
	if _, err = base64.StdEncoding.Decode(rawDoc, b64Doc); err != nil {
		return nil, err
	}

	return rawDoc, nil
}

func main() {
	var attestationFile, attestationEndpoint string
	var debug bool
	var err error
	var rawDoc []byte

	flag.StringVar(&attestationFile, "file", "", "File that contains a Base64-encoded attestation document.")
	flag.StringVar(&attestationEndpoint, "url", "", "Attestation endpoint of the Nitro enclave..")
	flag.BoolVar(&debug, "debug", false, "Print extra debug information.")
	flag.Parse()

	nonce := make([]byte, nonceSize)
	if _, err := rand.Read(nonce); err != nil {
		log.Fatalf("Failed to read random bytes: %s", err)
	}
	log.Printf("Created random nonce: %x", nonce)

	if attestationFile != "" {
		rawDoc, err = os.ReadFile(attestationFile)
		if err != nil {
			log.Fatalf("Failed to read attestatil file: %s", err)
		}
		log.Printf("Read attestation document from file.")
	} else {
		rawDoc, err = fetchAttestationDocument(fmt.Sprintf("%x", nonce), attestationEndpoint)
		if err != nil {
			log.Fatalf("Failed to fetch attestation document: %s", err)
		}
		log.Printf("Fetched attestation document from enclave.")
	}

	doc, err := nitrite.Verify(rawDoc, nitrite.VerifyOptions{})
	if err != nil {
		log.Fatalf("Failed to verify attestation document: %v", err)
	}
	log.Println("Signature verification successful.")

	if !bytes.Equal(doc.Document.Nonce, nonce) {
		log.Fatalf("Expected nonce %v but got %v; enclave liveness check failed.", nonce, doc.Document.Nonce)
	}
	log.Println("Correct nonce found in attestation document.")

	if debug {
		// See the following page for a description of the following values:
		// https://docs.aws.amazon.com/enclaves/latest/user/set-up-attestation.html
		for _, i := range []uint{0, 1, 2, 3, 4, 8} {
			fmt.Printf("PCR[%d]: %x\n", i, doc.Document.PCRs[i])
		}
	} else {
		pcr0, exists := doc.Document.PCRs[0]
		if !exists {
			log.Fatal("Attestation document does not contain PCR0.")
		}
		log.Println("Remote PCR0:")
		fmt.Printf("%x\n", pcr0)
	}
}
