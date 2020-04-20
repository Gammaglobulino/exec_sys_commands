package rsa_key_in_PEM

import (
	"../rsa_key_in_PEM"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

var defaultKeySize = 2048

func TestRSAPrivateKeyGenerator(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, defaultKeySize)
	assert.Nil(t, err)
	assert.NotEmpty(t, privateKey)
	log.Println(*privateKey)
}

func TestGeneratePrivatePEM(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, defaultKeySize)
	assert.Nil(t, err)
	assert.NotEmpty(t, privateKey)
	encodedPK := x509.MarshalPKCS1PrivateKey(privateKey)
	var privatePEM = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: encodedPK,
	}
	assert.NotEmpty(t, privatePEM)
	log.Println(privatePEM)
}
func TestGeneratePublicPEMFromKey(t *testing.T) {
	pkey, err := rsa_key_in_PEM.GeneratePrivateRSAKey(2048)
	assert.Nil(t, err)
	publicPEM, err := rsa_key_in_PEM.GeneratePublicPEMFromKey(pkey.PublicKey)
	assert.Nil(t, err)
	assert.NotEmpty(t, publicPEM)
	log.Println(pem.EncodeToMemory(publicPEM))
}
func TestGeneratePrivatePEMKey(t *testing.T) {
	pkey, err := rsa_key_in_PEM.GeneratePrivateRSAKey(2048)
	assert.Nil(t, err)
	privatePEM, err := rsa_key_in_PEM.GeneratePrivatePEMKey(pkey)
	assert.Nil(t, err)
	assert.NotEmpty(t, privatePEM)
	log.Println(privatePEM)

}
func TestSavePemToFile(t *testing.T) {
	pkey, err := rsa_key_in_PEM.GeneratePrivateRSAKey(2048)
	assert.Nil(t, err)
	publicPEM, err := rsa_key_in_PEM.GeneratePublicPEMFromKey(pkey.PublicKey)
	assert.Nil(t, err)
	assert.NotEmpty(t, publicPEM)
	log.Println(pem.EncodeToMemory(publicPEM))
	privatePEM, err := rsa_key_in_PEM.GeneratePrivatePEMKey(pkey)
	assert.Nil(t, err)
	assert.NotEmpty(t, privatePEM)
	log.Println(pem.EncodeToMemory(privatePEM))
	err = rsa_key_in_PEM.SavePemToFile(publicPEM, "publicPEM")
	assert.Nil(t, err)
	err = rsa_key_in_PEM.SavePemToFile(privatePEM, "privatePEM")
	assert.Nil(t, err)

}
func TestLoadPKFromPEMFile(t *testing.T) {
	privateKey, err := rsa_key_in_PEM.LoadPrivateKFromPEMFile("privatePEM")
	assert.Nil(t, err)
	assert.NotEmpty(t, privateKey)
	log.Println(privateKey)
}
