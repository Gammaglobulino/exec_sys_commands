package rsa_key_in_PEM

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

func GeneratePrivateRSAKey(keySize int) (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, err
	}
	return privateKey, err

}
func GeneratePrivatePEMKey(key *rsa.PrivateKey) (*pem.Block, error) {
	encodedPK := x509.MarshalPKCS1PrivateKey(key)
	if encodedPK == nil {
		return nil, errors.New("Issue during key encoding process")
	}
	var privatePEM = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: encodedPK,
	}

	return privatePEM, nil
}

func GeneratePublicPEMFromKey(publicKey rsa.PublicKey) (*pem.Block, error) {
	encodedPubKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return nil, err
	}
	var publicPem = &pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   encodedPubKey,
	}
	return publicPem, nil
}

func SavePemToFile(pemBlock *pem.Block, filename string) error {
	pemfile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer pemfile.Close()
	err = pem.Encode(pemfile, pemBlock)
	if err != nil {
		return err
	}
	return nil
}
