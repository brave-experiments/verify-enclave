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
	"net/url"
	"os"

	"github.com/hf/nitrite"
)

const (
	nonceSize = 20 // Nonce size in bytes.
)

func fetchAttestationDocument(nonce, attestationEndpoint string) ([]byte, error) {

	resp, err := http.PostForm(attestationEndpoint, url.Values{"nonce": {nonce}})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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
	var err error
	var rawDoc []byte

	flag.StringVar(&attestationFile, "file", "", "File that contains a Base64-encoded attestation document.")
	flag.StringVar(&attestationEndpoint, "url", "", "Attestation endpoint of the Nitro enclave..")
	flag.Parse()

	nonce := make([]byte, nonceSize)
	if _, err := rand.Read(nonce); err != nil {
		log.Fatal(err)
	}
	log.Printf("Created random nonce: %x", nonce)

	if attestationFile != "" {
		rawDoc, err = os.ReadFile(attestationFile)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Read attestation document from file.")
	} else {
		rawDoc, err = fetchAttestationDocument(fmt.Sprintf("%x", nonce), attestationEndpoint)
		if err != nil {
			log.Fatal(err)
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

	pcr0, exists := doc.Document.PCRs[0]
	if !exists {
		log.Fatal("Attestation document does not contain PCR0.")
	}
	log.Println("Remote PCR0:")
	fmt.Printf("%x\n", pcr0)
}
