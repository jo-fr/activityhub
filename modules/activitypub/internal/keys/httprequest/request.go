package httprequest

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jo-fr/activityhub/pkg/util/httputil"
	"github.com/pkg/errors"
)

const (
	// defautltSignatureHeaders defines which http headers are include in the signature.
	// for more information check https://docs.joinmastodon.org/spec/security/#http
	defautltSignatureHeaders = "(request-target) host date digest content-type"

	defaultContentType = "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\""
)

type Request struct {
	*http.Request
}

// New creates a new request with the given method, url and body.
func New(method string, url string, body io.Reader) (*Request, error) {
	r, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	r.Header.Set("Accept", "application/activity+json")

	return &Request{r}, nil
}

// SetHeader sets the header with the given key and value.
func (r *Request) SetHeader(key string, value string) {
	r.Header.Set(key, value)
}

// Do sends the request and returns the response.
func (r *Request) Do() (*http.Response, error) {
	c := http.Client{}

	return c.Do(r.Request)
}

// Sign adds required headers and signs the request with the given private key.
// In order for this function to work properly make sure that the request URL and body are set.
func (r *Request) Sign(privateKeyPEM []byte, actorURI string) error {
	privateKey, err := parsePrivateKey(privateKeyPEM)
	if err != nil {
		return errors.Wrap(err, "failed to parse private key")
	}

	body, err := httputil.GetBody(r.Request)
	if err != nil {
		return errors.Wrap(err, "failed to read request body")
	}

	tNow := time.Now().UTC()
	digest := createBodyDigest(body.Bytes())
	signatureString := createSignatureString(digest, r.Method, r.URL.Path, r.Host, tNow)
	signatureBytes, err := signWithRSA(privateKey, signatureString)
	if err != nil {
		return errors.Wrap(err, "failed to sign signature string")
	}

	signature := base64.StdEncoding.EncodeToString(signatureBytes)
	keyID := fmt.Sprintf("%s#main-key", actorURI)

	signatureHeaderString := fmt.Sprintf(`keyId="%s",headers="%s",signature="%s"`, keyID, defautltSignatureHeaders, signature)

	// set headers
	r.SetHeader("Date", tNow.Format(http.TimeFormat))
	r.SetHeader("Digest", digest)
	r.SetHeader("Signature", signatureHeaderString)
	r.SetHeader("Content-Type", defaultContentType)

	return nil
}
