package digital_signed_message

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

func VerifySignature(signature []byte, message []byte, publicKey *rsa.PublicKey) bool {
	hashed := sha256.Sum256(message)
	err := rsa.VerifyPKCS1v15(
		publicKey,
		crypto.SHA256,
		hashed[:],
		signature,
	)
	if err != nil {
		return false
	}
	return true
}

func SignMessage(privateKey *rsa.PrivateKey, message []byte) ([]byte, error) {
	hashed := sha256.Sum256(message)
	signature, err := rsa.SignPKCS1v15(
		rand.Reader,
		privateKey,
		crypto.SHA256,
		hashed[:],
	)
	if err != nil {
		return nil, err
	}
	return signature, nil

}
