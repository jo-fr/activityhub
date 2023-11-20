package httprequest

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func parsePrivateKey(privateKeyPem []byte) (*rsa.PrivateKey, error) {
	pemBlock, _ := pem.Decode(privateKeyPem)
	privKey, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return privKey, nil
}

func createSignatureString(bodyDigest string, method string, target string, hostURL string, date time.Time) string {
	var signatureStrings []string
	signatureStrings = append(signatureStrings, fmt.Sprintf("(request-target): %s %s", strings.ToLower(method), target))
	signatureStrings = append(signatureStrings, fmt.Sprintf("host: %s", hostURL))
	signatureStrings = append(signatureStrings, fmt.Sprintf("date: %s", date.Format(http.TimeFormat)))
	signatureStrings = append(signatureStrings, fmt.Sprintf("digest: %s", bodyDigest))
	signatureStrings = append(signatureStrings, fmt.Sprintf("content-type: %s", defaultContentType))
	return strings.Join(signatureStrings, "\n")
}

func signWithRSA(key *rsa.PrivateKey, data string) ([]byte, error) {
	h := sha256.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hashed)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

func createBodyDigest(requestBody []byte) string {
	h := sha256.New()
	h.Write(requestBody)
	digest := "sha-256=" + base64.StdEncoding.EncodeToString(h.Sum(nil))
	return digest
}
