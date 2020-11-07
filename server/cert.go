package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/sirupsen/logrus"
	"math/big"
	"net"
	"os"
	"strings"
	"time"
)

var (
	validFrom = ""
	validFor  = 365 * 24 * time.Hour
	isCA      = true
	rsaBits   = 2048
)

func checkIfCertExists(certPath string, keyPath string) error {
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		return err
	} else if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		return err
	}
	return nil
}

func generateCerts(certPath string, keyPath string, host string) error {
	var priv interface{}
	var err error

	priv, err = rsa.GenerateKey(rand.Reader, rsaBits)

	if err != nil {
		logrus.Errorf("failed to generate private key: %v", err)
		return err
	}

	var notBefore time.Time
	if len(validFrom) == 0 {
		notBefore = time.Now()
	} else {
		notBefore, err = time.Parse("Jan 2 15:04:05 2006", validFrom)
		if err != nil {
			logrus.Errorf("Failed to parse creation date: %v", err)
			return err
		}
	}

	notAfter := notBefore.Add(validFor)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		logrus.Errorf("failed to generate serial number: %v", err)
		return err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Seknox"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	hosts := strings.Split(host, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	if isCA {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
		logrus.Errorf("Failed to create certificate: %s", err)
		return err
	}

	certOut, err := os.Create(certPath)
	if err != nil {
		logrus.Errorf("failed to open %s for writing: %v", certPath, err)
		return err
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()
	logrus.Info("written tras-server.crt")

	keyOut, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		logrus.Errorf("failed to open %s for writing: %v", keyPath, err)
		return err
	}
	pem.Encode(keyOut, pemBlockForKey(priv))
	keyOut.Close()
	logrus.Info("written trasa-server.key")
	return nil
}

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func pemBlockForKey(priv interface{}) *pem.Block {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			logrus.Errorf("Unable to marshal ECDSA private key: %v", err)
			os.Exit(2)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}
