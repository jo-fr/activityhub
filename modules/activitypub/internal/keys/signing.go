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

	"strings"
	"time"

	"github.com/jo-fr/activityhub/modules/api/httputil"
	"github.com/pkg/errors"
)

const (
	// defautltSignatureHeaders defines which http headers are include in the signature.
	// for more information check https://docs.joinmastodon.org/spec/security/#http
	defautltSignatureHeaders = "(request-target) host date digest content-type"

	defaultContentType = "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\""
)

// SignRequesst adds required headers and signs the request with the given private key.
// In order for this function to work properly make sure that the request URL and body are set.
func SignRequest(r *http.Request, privateKeyPEM []byte, actorURI string) (*http.Request, error) {
	privateKey, err := parsePrivateKey(privateKeyPEM)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse private key")
	}

	body, err := httputil.GetBody(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read request body")
	}

	tNow := time.Now().UTC()
	digest := createBodyDigest(body.Bytes())
	signatureString := createSignatureString(digest, r.Method, r.URL.Path, r.Host, tNow)
	signatureBytes, err := signWithRSA(privateKey, signatureString)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign signature string")
	}

	signature := base64.StdEncoding.EncodeToString(signatureBytes)
	keyID := fmt.Sprintf("%s#main-key", actorURI)

	signatureHeaderString := fmt.Sprintf(`keyId="%s",headers="%s",signature="%s"`, keyID, defautltSignatureHeaders, signature)

	// set headers
	r.Header.Set("Date", tNow.Format(http.TimeFormat))
	r.Header.Set("Digest", digest)
	r.Header.Set("Signature", signatureHeaderString)
	r.Header.Set("Content-Type", defaultContentType)

	return r, nil
}

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
