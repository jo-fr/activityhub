package middleware

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// example Signature header:
// Signature: keyId="https://my-example.com/actor#main-key",headers="(request-target) host date",signature="Y2FiYW...IxNGRiZDk4ZA=="

func ValidateSignature(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		request := r

		// parse header
		signatureHeader := request.Header.Get("Signature")
		if signatureHeader == "" {
			rw.Write([]byte("No signature header"))
			return
		}

		signatureHeaderMap := parseSignatureHeader(signatureHeader)

		keyID := signatureHeaderMap["keyId"]
		headers := signatureHeaderMap["headers"]
		signatureBytes, _ := base64.StdEncoding.DecodeString(signatureHeaderMap["signature"])

		actor, err := fetchActor(keyID)
		if err != nil {
			fmt.Println(err)
			return
		}

		rsaPubKey, err := parsePublicKey(actor)
		if err != nil {
			fmt.Println(err)
			return
		}

		comparisonStrings := buildComparisonStrings(request, headers)
		comparisonHash := calculateHash(comparisonStrings)

		if err = verifySignature(rsaPubKey, signatureBytes, comparisonHash); err != nil {
			rw.Write([]byte("verification failed"))
			return
		}

		next.ServeHTTP(rw, r)
	})
}

func Validate(r *http.Request) error {
	request := r

	// parse header
	signatureHeader := request.Header.Get("Signature")
	if signatureHeader == "" {
		return errors.New("no signature header")
	}

	signatureHeaderMap := parseSignatureHeader(signatureHeader)

	keyID := signatureHeaderMap["keyId"]
	headers := signatureHeaderMap["headers"]
	signatureBytes, err := base64.StdEncoding.DecodeString(signatureHeaderMap["signature"])
	if err != nil {
		return err
	}

	actor, err := fetchActor(keyID)
	if err != nil {
		return err
	}

	rsaPubKey, err := parsePublicKey(actor)
	if err != nil {
		return err
	}

	comparisonStrings := buildComparisonStrings(request, headers)
	comparisonHash := calculateHash(comparisonStrings)

	if err = verifySignature(rsaPubKey, signatureBytes, comparisonHash); err != nil {
		return err
	}
	return nil
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

func buildComparisonStrings(request *http.Request, headers string) string {
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

func calculateHash(comparisonString string) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(comparisonString))

	return hasher.Sum(nil)
}

func verifySignature(rsaPubKey *rsa.PublicKey, signatureBytes []byte, comparisonHash []byte) error {
	return rsa.VerifyPKCS1v15(rsaPubKey, crypto.SHA256, comparisonHash, signatureBytes)
}

func fetchActor(keyID string) (map[string]interface{}, error) {

	req, err := http.NewRequest("GET", keyID, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var actor map[string]interface{}
	if err := json.Unmarshal(body, &actor); err != nil {
		return nil, err
	}

	return actor, nil
}

func parsePublicKey(actor map[string]interface{}) (*rsa.PublicKey, error) {
	publicKeyPEM := actor["publicKey"].(map[string]interface{})["publicKeyPem"].(string)
	publicKeyPEM = "-----BEGIN PUBLIC KEY-----\n" + publicKeyPEM + "\n-----END PUBLIC KEY-----"

	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the public key")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("key is not an RSA public key")
	}

	return rsaPubKey, nil
}
