package tls_self_signed_certificate

import (
	"../../client_server_connection/core/handle_connections"
	"../rsa_key_in_PEM"
	"../tls_self_signed_certificate"
	"crypto/rand"
	"crypto/x509"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestSetupCertificateTemplate(t *testing.T) {
	localIp, err := handle_connections.GetLocalIp04()
	assert.Nil(t, err)
	certificate, err := tls_self_signed_certificate.SetupCertificateTemplate(false, localIp)
	assert.Nil(t, err)
	assert.NotEmpty(t, certificate)
	log.Println(certificate)

}
func TestWriteCertificateToPemFileName(t *testing.T) {
	localIp, err := handle_connections.GetLocalIp04()
	assert.Nil(t, err)
	tCertificate, err := tls_self_signed_certificate.SetupCertificateTemplate(false, localIp)
	assert.Nil(t, err)
	assert.NotEmpty(t, tCertificate)
	privateKey, err := rsa_key_in_PEM.LoadPrivateKFromPEMFile("privatePEM")
	publicKey := privateKey.PublicKey
	certificate, err := x509.CreateCertificate(
		rand.Reader,
		&tCertificate,
		&tCertificate,
		&publicKey,
		privateKey,
	)
	assert.Nil(t, err)
	err = tls_self_signed_certificate.WriteCertificateToPemFileName("PEMCertificate", certificate)
	assert.Nil(t, err)
}
