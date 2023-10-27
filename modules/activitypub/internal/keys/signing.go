package keys

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
	"time"

	"github.com/pkg/errors"
)

const (
	// signatureHeaders defines which http headers are include in the signature.
	// for more information check https://docs.joinmastodon.org/spec/security/#http
	signatureHeaders = "(request-target) host date digest"
)

// CreateSignature creates the Signature header string.
func CreateSignature(bodyDigest string, hostURL string, keyId string, privateKeyPem []byte) (string, error) {
	key, err := parsePrivateKey(privateKeyPem)
	if err != nil {
		return "", err
	}

	signatureString := createSignatureString(bodyDigest, http.MethodPost, "/users/test/inbox", hostURL, signatureHeaders)
	signatureBytes, err := signWithRSA(key, signatureString)
	if err != nil {
		return "", errors.Wrap(err, "failed to sign signature string")
	}

	signature := base64.StdEncoding.EncodeToString(signatureBytes)

	return fmt.Sprintf(`keyId="%s",headers="%s",signature="%s"`, keyId, signatureHeaders, signature), nil
}

func parsePrivateKey(privateKeyPem []byte) (*rsa.PrivateKey, error) {

	pemBlock, _ := pem.Decode(privateKeyPem)
	privKey, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return privKey, nil
}

func createSignatureString(bodyDigest string, method string, target string, hostURL string, headers string) string {
	signatureString := fmt.Sprintf("(request-target): %s %s\n", method, target)
	signatureString += fmt.Sprintf("host: %s\n", hostURL)
	signatureString += fmt.Sprintf("date: %s\n", time.Now().UTC().Format(http.TimeFormat))
	signatureString += fmt.Sprintf("digest: %s\n", bodyDigest)
	return signatureString
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

func CreateBodyDigest(requestBody string) string {
	h := sha256.New()
	h.Write([]byte(requestBody))
	digest := "sha-256=" + base64.StdEncoding.EncodeToString(h.Sum(nil))
	return digest
}
