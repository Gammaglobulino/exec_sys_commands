package tls_self_signed_certificate

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"time"
)

func SetupCertificateTemplate(isCA bool, ip net.IP) (x509.Certificate, error) {
	notBefore := time.Now()
	notAfter := notBefore.Add(time.Hour * 24 * 365)           //after one year
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128) //very big number
	randomNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return x509.Certificate{}, err
	}
	nameInfo := pkix.Name{
		Country:            []string{"IT"},
		Organization:       []string{"Gammapan"},
		OrganizationalUnit: []string{"Computer Science Dep"},
		Locality:           []string{"Poppi"},
		Province:           []string{"Arezzo"},
		StreetAddress:      []string{"Via Panoramica"},
	}
	certTemplate := x509.Certificate{
		SerialNumber:          randomNumber,
		Subject:               nameInfo,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
		BasicConstraintsValid: true,
		IsCA:                  false,
		IPAddresses:           []net.IP{ip},
		DNSNames:              []string{"localhost"},
	}
	if isCA {
		certTemplate.IsCA = true
		certTemplate.KeyUsage = certTemplate.KeyUsage | x509.KeyUsageCertSign
	}
	return certTemplate, nil
}

func WriteCertificateToPemFileName(pemFileName string, certBytes []byte) error {
	certPem := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}
	file, err := os.Create(pemFileName)
	if err != nil {
		return err
	}
	defer file.Close()
	pem.Encode(file, certPem)
	return err
}
