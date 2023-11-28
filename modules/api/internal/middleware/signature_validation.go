package middleware

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"

	"fmt"
	"io"

	"net/http"
	"strings"

	"github.com/jo-fr/activityhub/modules/api/internal/render"
	"github.com/jo-fr/activityhub/pkg/errutil"
	"github.com/jo-fr/activityhub/pkg/log"
	"github.com/jo-fr/activityhub/pkg/util"
	"github.com/jo-fr/activityhub/pkg/util/httputil"
	"github.com/pkg/errors"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// predefined errors
var (
	ErrSignatureHeaderMissing = errutil.NewError(errutil.TypeMissingHeader, "signature header missing")
	ErrSignatureNotValid      = errutil.NewError(errutil.TypeBadRequest, "signature not valid")
)

// ValidateSignature validates the signature header of the request. The signature header is expected to be in the following format:
// Signature: keyId="https://my-example.com/actor#main-key",headers="(request-target) host date",signature="Y2FiYW...IxNGRiZDk4ZA=="
func ValidateSignature(log *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			request := r

			// parse header
			signatureHeader := request.Header.Get("Signature")
			if signatureHeader == "" {
				render.Error(r.Context(), ErrSignatureHeaderMissing, rw, log)
				return
			}

			signatureHeaderMap := parseSignatureHeader(signatureHeader)

			keyID := signatureHeaderMap["keyId"]
			headers := signatureHeaderMap["headers"]
			signatureBytes, _ := util.DecodeBase64(signatureHeaderMap["signature"])

			actor, err := fetchActor(keyID)
			if err != nil {
				render.Error(r.Context(), err, rw, log)
				return
			}

			rsaPubKey, err := parsePublicKey(actor)
			if err != nil {
				render.Error(r.Context(), err, rw, log)
				return
			}

			comparisonString := buildComparisonString(request, headers)
			comparisonHash := util.HashSHA256(comparisonString)

			if err = verifySignature(rsaPubKey, signatureBytes, comparisonHash); err != nil {
				render.Error(r.Context(), ErrSignatureNotValid, rw, log)
				return
			}

			next.ServeHTTP(rw, r)
		})
	}
}

func parseSignatureHeader(signatureHeader string) map[string]string {
	parts := strings.Split(signatureHeader, ",")
	headerMap := make(map[string]string)

	for _, pair := range parts {
		keyValue := strings.SplitN(pair, "=", 2)
		key := strings.Trim(keyValue[0], " \t\"")
		value := strings.Trim(keyValue[1], " \t\"")
		headerMap[key] = value
	}

	return headerMap
}

func buildComparisonString(request *http.Request, headers string) string {
	signedHeaders := strings.Split(headers, " ")
	var comparisonStrings []string

	for _, signedHeaderName := range signedHeaders {
		switch signedHeaderName {
		case "(request-target)":
			comparisonStrings = append(comparisonStrings, "(request-target): "+strings.ToLower(request.Method)+" "+request.URL.Path)
		case "host":
			comparisonStrings = append(comparisonStrings, "host: "+request.Host)
		default:
			capitalizedHeaderName := cases.Title(language.AmericanEnglish).String(strings.ToLower(signedHeaderName))
			headerValue := request.Header.Get(capitalizedHeaderName)
			comparisonStrings = append(comparisonStrings, fmt.Sprintf("%s: %s", signedHeaderName, headerValue))
		}
	}

	signatureString := strings.Join(comparisonStrings, "\n")

	return signatureString
}

func verifySignature(rsaPubKey *rsa.PublicKey, signatureBytes []byte, comparisonHash []byte) error {
	return rsa.VerifyPKCS1v15(rsaPubKey, crypto.SHA256, comparisonHash, signatureBytes)
}

func fetchActor(keyID string) (map[string]interface{}, error) {

	req, err := http.NewRequest("GET", keyID, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send request")
	}
	defer response.Body.Close()

	if !httputil.StatusOK(response.StatusCode) {
		return nil, errors.Errorf("request failed with status code %s", response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	var actor map[string]interface{}
	if err := json.Unmarshal(body, &actor); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response body")
	}

	return actor, nil
}

func parsePublicKey(actor map[string]interface{}) (*rsa.PublicKey, error) {
	publicKeyPEM := actor["publicKey"].(map[string]interface{})["publicKeyPem"].(string)

	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the public key")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse DER encoded public key")
	}

	pk, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("failed to parse DER encoded public key")
	}

	return pk, nil
}
