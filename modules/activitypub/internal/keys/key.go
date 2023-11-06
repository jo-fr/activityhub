package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"

	"github.com/jo-fr/activityhub/modules/activitypub/models"
	"github.com/pkg/errors"
)

// GenerateRSAKeyPair generates a new key pair of specified byte size
func GenerateRSAKeyPair(bitSize int) (*models.KeyPair, error) {
	privateKey, err := generatePrivateKey(bitSize)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate private key")
	}

	return &models.KeyPair{
		PrivKeyPEM: encodePrivateKeyToPEM(privateKey),
		PubKeyPEM:  encodePublicKeyToPEM(&privateKey.PublicKey),
	}, nil
}

// generatePrivateKey creates a RSA Private Key of specified byte size
func generatePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate private key")
	}

	log.Println("Private Key generated")
	return privateKey, nil
}

// encodePrivateKeyToPEM encodes Private Key from RSA to PEM format
func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privBlock)

	return privatePEM
}

// encodePublicKeyToPEM encodes Private Key from RSA to PEM format
func encodePublicKeyToPEM(publicKey *rsa.PublicKey) []byte {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PublicKey(publicKey)

	// pem.Block
	privBlock := pem.Block{
		Type:    "RSA PUBLIC KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privBlock)

	return privatePEM
}
