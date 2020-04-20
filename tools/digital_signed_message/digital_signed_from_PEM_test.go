package digital_signed_message

import (
	"../digital_signed_message"
	"../rsa_key_in_PEM"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"testing"
)

func TestSignMessage(t *testing.T) {
	privateKey, err := rsa_key_in_PEM.LoadPrivateKFromPEMFile("privatePEM")
	assert.Nil(t, err)
	assert.NotEmpty(t, privateKey)
	signature, err := digital_signed_message.SignMessage(privateKey, []byte("Sir Andrea Mazzanti da Poppi"))
	assert.Nil(t, err)
	assert.NotEmpty(t, signature)
	ioutil.WriteFile("signature", signature, 0666)
	log.Println(signature)
}

func TestVerifySignature(t *testing.T) {
	signature, err := ioutil.ReadFile("signature")
	assert.Nil(t, err)
	publicKey, err := rsa_key_in_PEM.LoadPublicKFromPEMFile("publicPEM")
	assert.Nil(t, err)
	assert.True(t, digital_signed_message.VerifySignature(signature, []byte("Sir Andrea Mazzanti da Poppi"), publicKey))
}
