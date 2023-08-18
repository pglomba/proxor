package certificates

import (
	"bytes"
	"crypto/tls"
	"io"
	"os"
)

type Certificate struct {
	Name                 string
	CertificateFile      string
	CertificateChainFile string
	CertificateKey       string
}

func readFileToBuffer(fileName string) ([]byte, error) {
	f, err := os.Open(fileName)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, f)

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func LoadCertificates(configCertificates []Certificate) ([]tls.Certificate, error) {
	var certificates []tls.Certificate

	for _, configCertificate := range configCertificates {
		certPem, _ := readFileToBuffer(configCertificate.CertificateFile)

		if configCertificate.CertificateChainFile != "" {
			chainPem, _ := readFileToBuffer(configCertificate.CertificateChainFile)
			certPem = append(certPem, chainPem...)
		}

		keyPem, _ := readFileToBuffer(configCertificate.CertificateKey)

		cert, err := tls.X509KeyPair(certPem, keyPem)

		if err != nil {
			return nil, err
		}

		certificates = append(certificates, cert)
	}

	return certificates, nil
}
