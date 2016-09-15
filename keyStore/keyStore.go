package keyStore

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/echocat/caretakerd/errors"
	"github.com/echocat/caretakerd/panics"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"
)

// KeyStore represents a keystore that holds certificates, CAs and private keys.
type KeyStore struct {
	enabled    bool
	config     Config
	pem        []byte
	ca         []*x509.Certificate
	cert       *x509.Certificate
	privateKey interface{}
}

// NewKeyStore create an new instance of KeyStore.
func NewKeyStore(enabled bool, conf Config) (*KeyStore, error) {
	err := conf.Validate()
	if err != nil {
		return nil, err
	}
	if !enabled {
		return &KeyStore{
			enabled: false,
			config:  conf,
		}, nil
	}
	switch conf.Type {
	case FromFile:
		return newFomFile(conf)
	case FromEnvironment:
		return newFromEnvironment(conf)
	case Generated:
		return newGenerated(conf)
	}
	return nil, errors.New("Unknown keyStore type %v.", conf.Type)
}

func generatePrivateKey(conf Config) (privateKey interface{}, privateKeyBytes []byte, publicKey interface{}, err error) {
	plainAlgorithm := conf.GetHintsArgument("algorithm")
	if len(plainAlgorithm) > 0 && strings.ToLower(plainAlgorithm) != "rsa" {
		return nil, []byte{}, nil, errors.New("Unsupported algorithm: %s", plainAlgorithm)
	}
	bits := 1024
	plainBits := conf.GetHintsArgument("bits")
	if len(plainBits) > 0 {
		if bits, err = strconv.Atoi(plainBits); err != nil || bits <= 0 {
			return nil, []byte{}, nil, errors.New("Unsupported algorithm bits: %s", plainBits)
		}
	}
	plainPrivateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, []byte{}, nil, errors.New("Could not generate private key.").CausedBy(err)
	}
	privateKeyBytes = x509.MarshalPKCS1PrivateKey(plainPrivateKey)
	return plainPrivateKey, privateKeyBytes, &plainPrivateKey.PublicKey, nil
}

func generateCertificate(conf Config, privateKey interface{}, publicKey interface{}) ([]byte, error) {
	notBefore := time.Now()
	notAfter := notBefore.Add(8760 * time.Hour)

	template := x509.Certificate{
		SerialNumber: newSerialNumber(),
		Subject: pkix.Name{
			CommonName: "caretakerd",
		},
		IsCA:                  true,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey, privateKey)
	if err != nil {
		return nil, errors.New("Failed to create certificate.").CausedBy(err)
	}
	return derBytes, nil
}

func generatePem(conf Config) ([]byte, *x509.Certificate, interface{}, error) {
	privateKey, privateKeyBytes, publicKey, err := generatePrivateKey(conf)
	if err != nil {
		return []byte{}, nil, nil, errors.New("Could not generate private key.").CausedBy(err)
	}
	certificateDerBytes, err := generateCertificate(conf, privateKey, publicKey)
	if err != nil {
		return []byte{}, nil, nil, errors.New("Could not generate certificate.").CausedBy(err)
	}
	cert, err := x509.ParseCertificate(certificateDerBytes)
	if err != nil {
		return []byte{}, nil, nil, errors.New("Wow! Could not parse right now created certificate?").CausedBy(err)
	}

	pemBytes := []byte{}
	pemBytes = append(pemBytes, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certificateDerBytes})...)
	pemBytes = append(pemBytes, pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privateKeyBytes})...)

	return pemBytes, cert, privateKey, nil
}

func newFomFile(conf Config) (*KeyStore, error) {
	pem, err := ioutil.ReadFile(conf.PemFile.String())
	if err != nil {
		return nil, errors.New("Could not read pem from '%v'.", conf.PemFile).CausedBy(err)
	}
	return newPemFromBytes(conf, pem)
}

func newFromEnvironment(conf Config) (*KeyStore, error) {
	pem := os.Getenv("CTD_PEM")
	if len(strings.TrimSpace(pem)) <= 0 {
		return nil, errors.New("There is an %v keyStore confgiured but the CTD_PEM environment varaible is empty.", conf.Type)
	}
	return newPemFromBytes(conf, []byte(pem))
}

func newPemFromBytes(conf Config, pem []byte) (*KeyStore, error) {
	ca, err := buildWholeCAsBy(conf, pem)
	if err != nil {
		return nil, errors.New("Could not build ca for keyStore config.").CausedBy(err)
	}
	certs, err := loadCertificatesFrom(pem)
	if err != nil {
		return nil, errors.New("Could not load certs from PEM.").CausedBy(err)
	}
	if len(certs) <= 0 {
		return nil, errors.New("The provieded PEM does not contain a certificate.")
	}
	privateKey, err := loadPrivateKeyFrom(pem)
	if err != nil {
		return nil, err
	}
	return &KeyStore{
		enabled:    true,
		config:     conf,
		pem:        pem,
		ca:         ca,
		cert:       certs[0],
		privateKey: privateKey,
	}, nil
}

func newGenerated(conf Config) (*KeyStore, error) {
	pem, cert, privateKey, err := generatePem(conf)
	if err != nil {
		return nil, errors.New("Could not generate pem for keyStore config.").CausedBy(err)
	}
	ca, err := buildWholeCAsBy(conf, pem)
	if err != nil {
		return nil, errors.New("Could not build CA bundle for keyStore config.").CausedBy(err)
	}
	return &KeyStore{
		enabled:    true,
		config:     conf,
		pem:        pem,
		ca:         ca,
		cert:       cert,
		privateKey: privateKey,
	}, nil
}

func buildWholeCAsBy(conf Config, p []byte) ([]*x509.Certificate, error) {
	result := []*x509.Certificate{}
	if !conf.CaFile.IsTrimmedEmpty() {
		fileContent, err := ioutil.ReadFile(conf.CaFile.String())
		if err != nil {
			return nil, errors.New("Could not read certificates from %v.", conf.CaFile).CausedBy(err)
		}
		casFromFile, err := loadCertificatesFrom(fileContent)
		if err != nil {
			return nil, errors.New("Could not parse certificates from %v.", conf.CaFile).CausedBy(err)
		}
		for _, candidate := range casFromFile {
			if candidate.IsCA {
				result = append(result, candidate)
			}
		}
	}
	casFromP, err := loadCertificatesFrom(p)
	if err != nil {
		return nil, err
	}
	return append(result, casFromP...), nil
}

func loadCertificatesFrom(p []byte) ([]*x509.Certificate, error) {
	result := []*x509.Certificate{}
	if len(p) > 0 {
		rp := p
		block := new(pem.Block)
		for block != nil && len(rp) > 0 {
			block, rp = pem.Decode(rp)
			if block != nil && block.Type == "CERTIFICATE" {
				candidates, err := x509.ParseCertificates(block.Bytes)
				if err != nil {
					return nil, errors.New("Could not parse certificates.").CausedBy(err)
				}
				for _, candidate := range candidates {
					if candidate.IsCA {
						result = append(result, candidate)
					}
				}
			}
		}
	}
	return result, nil
}

func loadPrivateKeyFrom(p []byte) (interface{}, error) {
	if len(p) > 0 {
		rp := p
		block := new(pem.Block)
		for block != nil && len(rp) > 0 {
			block, rp = pem.Decode(rp)
			if block != nil && block.Type == "RSA PRIVATE KEY" {
				privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
				if err != nil {
					return nil, errors.New("Could not parse privateKey.").CausedBy(err)
				}
				return privateKey, nil
			}
		}
	}
	return nil, errors.New("The PEM does not contain a valid private key.")
}

func newSerialNumber() *big.Int {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		panics.New("Could not generate serial number.").CausedBy(err).Throw()
	}
	return serialNumber
}

// LoadCertificateFromFile loads a certificate from the given filename and returns it.
func LoadCertificateFromFile(filename string) (*x509.Certificate, error) {
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.New("Could not read certificate from %v.", filename).CausedBy(err)
	}
	certificates, err := loadCertificatesFrom(fileContent)
	if err != nil {
		return nil, errors.New("Could not read certificate from %v.", filename).CausedBy(err)
	}
	if len(certificates) <= 0 {
		return nil, errors.New("File %v does not contain a valid certificate.", filename)
	}
	return certificates[0], nil
}

func (instance KeyStore) generateClientCertificate(name string, publicKey interface{}, privateKey interface{}) ([]byte, error) {
	notBefore := time.Now()
	notAfter := notBefore.Add(8760 * time.Hour)

	template := x509.Certificate{
		SerialNumber: newSerialNumber(),
		Issuer:       instance.cert.Subject,
		Subject: pkix.Name{
			CommonName: name,
		},
		IsCA:                  true,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: false,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, instance.cert, publicKey, instance.privateKey)
	if err != nil {
		return []byte{}, errors.New("Failed to create certificate for '%v'.", name).CausedBy(err)
	}
	return derBytes, nil
}

// GeneratePem generates a new PEM with the config of the current KeyStore instance and returns it.
// This PEM will be stored in the KeyStore instance.
func (instance KeyStore) GeneratePem(name string) ([]byte, *x509.Certificate, error) {
	if !instance.enabled {
		return []byte{}, nil, errors.New("KeyStore is not enabled.")
	}
	privateKey, privateKeyBytes, publicKey, err := generatePrivateKey(instance.Config())
	if err != nil {
		return []byte{}, nil, errors.New("Could not generate pem for '%v'.", name).CausedBy(err)
	}
	certificateDerBytes, err := instance.generateClientCertificate(name, publicKey, privateKey)
	if err != nil {
		return []byte{}, nil, err
	}

	cert, err := x509.ParseCertificate(certificateDerBytes)
	if err != nil || cert == nil {
		return []byte{}, nil, errors.New("Wow! Could not parse right now created certificate for '%v'?", name).CausedBy(err)
	}

	pemBytes := []byte{}
	pemBytes = append(pemBytes, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certificateDerBytes})...)
	pemBytes = append(pemBytes, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: instance.cert.Raw})...)
	pemBytes = append(pemBytes, pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privateKeyBytes})...)

	return pemBytes, cert, nil
}

// PEM returns the contained PEM instance of this KeyStore.
// If there is no PEM the result is empty.
func (instance KeyStore) PEM() []byte {
	return instance.pem
}

// CA returns all contained CAs of this KeyStore.
func (instance KeyStore) CA() []*x509.Certificate {
	return instance.ca
}

// Type returns the Type of this KeyStore.
func (instance KeyStore) Type() Type {
	return instance.config.Type
}

// Config returns the Config instance this KeyStore was created with.
func (instance KeyStore) Config() Config {
	return instance.config
}

// IsCA returns "true" if the contained certificate could be used to create new certificates.
func (instance KeyStore) IsCA() bool {
	cert := instance.cert
	return cert != nil && cert.IsCA
}

// IsEnabled returns "true" if this KeyStore is configured and usable.
func (instance KeyStore) IsEnabled() bool {
	return instance.enabled
}
