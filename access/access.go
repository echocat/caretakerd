package access

import (
	"github.com/echocat/caretakerd/errors"
	"os"
	"os/user"
	"github.com/echocat/caretakerd/keyStore"
	"crypto/x509"
	"reflect"
)

type Access struct {
	name              string
	t                 Type
	permission        Permission
	pem               []byte
	cert              *x509.Certificate
	temporaryFilename *string
}

func NewAccess(conf Config, name string, ks *keyStore.KeyStore) (*Access, error) {
	err := conf.Validate()
	if err != nil {
		return nil, err
	}
	if !ks.IsEnabled() {
		return newNoneInstance(name)
	}
	switch conf.Type {
	case None:
		return newNoneInstance(name)
	case Trusted:
		return newTrustedInstance(conf, name, ks)
	case GenerateToEnvironment:
		return newGenerateToEnvironmentInstance(conf, name, ks)
	case GenerateToFile:
		return newGenerateToFileInstance(conf, name, ks)
	}
	return nil, errors.New("Unknown access type %v.", conf.Type)
}

func newNoneInstance(name string) (*Access, error) {
	return &Access{
		t: None,
		permission: Forbidden,
		name: name,
	}, nil
}

func newTrustedInstance(conf Config, name string, ks *keyStore.KeyStore) (*Access, error) {
	if len(ks.Ca()) == 0 {
		return nil, errors.New("If there is valid caFile configured %v access could not work.", Trusted)
	}
	var cert *x509.Certificate
	if !conf.PemFile.IsTrimmedEmpty() {
		var err error
		cert, err = ks.LoadCertificateFromFile(conf.PemFile.String())
		if err != nil {
			return nil, errors.New("Could not load certificate from pemFile %v of service %v.", conf.PemFile, name)
		}
	}
	return &Access{
		t: Trusted,
		permission: conf.Permission,
		name: name,
		cert: cert,
	}, nil
}

func checkForIsCa(name string, sec *keyStore.KeyStore) error {
	if !sec.IsCA() {
		return errors.New("It is not possible to generate a new certificate for service '%v' with a caretakerd certificate that is not a CA. " +
		"Use trusted access for service '%v', configure caretakerd to generate its own certificate or provide a CA enabled certificate for caretakerd.", name, name)
	}
	return nil
}

func newGenerateToEnvironmentInstance(conf Config, name string, ks *keyStore.KeyStore) (*Access, error) {
	if err := checkForIsCa(name, ks); err != nil {
		return nil, err
	}
	pem, cert, err := ks.GeneratePem(name)
	if err != nil {
		return nil, errors.New("Could not generate pem for '%v'.", name).CausedBy(err)
	}
	return &Access{
		t: GenerateToEnvironment,
		permission: conf.Permission,
		name: name,
		pem: pem,
		cert: cert,
	}, nil
}

func newGenerateToFileInstance(conf Config, name string, ks *keyStore.KeyStore) (*Access, error) {
	if err := checkForIsCa(name, ks); err != nil {
		return nil, err
	}
	pem, cert, err := ks.GeneratePem(name)
	if err != nil {
		return nil, errors.New("Could not generate pem for '%v'.", name).CausedBy(err)
	}
	file, err := generateFileForPem(conf, pem)
	if err != nil {
		return nil, errors.New("Could not generate pem file for '%v'.", name).CausedBy(err)
	}
	return &Access{
		t: GenerateToFile,
		permission: conf.Permission,
		name: name,
		pem: pem,
		cert: cert,
		temporaryFilename: &file,
	}, nil
}

func generateFileForPem(conf Config, pem []byte) (string, error) {
	permission := conf.PemFilePermission.ThisOrDefault().AsFileMode()
	f, err := os.OpenFile(conf.PemFile.String(), os.O_WRONLY | os.O_CREATE | os.O_TRUNC, permission)
	if err != nil {
		return "", errors.New("Could not create pemFile '%s'.", conf.PemFile).CausedBy(err)
	}
	defer f.Close()
	if ! conf.PemFileUser.IsEmpty() {
		_, lerr := user.Lookup(conf.PemFileUser.String())
		if lerr != nil {
			return "", errors.New("Could not set ownership of pemFile '%s' to '%s'.", conf.PemFile, conf.PemFileUser).CausedBy(err)
		}
		//f.Chown(kfu.Uid, kfu.Gid) TODO!
	}
	f.Write(pem)
	f.Sync()
	return conf.PemFile.String(), nil
}

func (instance Access) Pem() []byte {
	return instance.pem
}

func (instance Access) Type() Type {
	return instance.t
}

func (instance Access) Cleanup() {
	if instance.temporaryFilename != nil {
		os.Remove(*instance.temporaryFilename)
	}
}

func (instance Access) HasReadPermission() bool {
	permission := instance.permission
	return permission == ReadOnly || permission == ReadWrite
}

func (instance Access) HasWritePermission() bool {
	permission := instance.permission
	return permission == ReadWrite
}

func (instance *Access) IsCertValid(cert *x509.Certificate) bool {
	instanceCert := instance.cert
	if instance.t == None {
		return false
	} else if cert == nil && instanceCert == nil {
		return false
	} else if cert != nil && instanceCert != nil {
		thatPublicKey := cert.PublicKey
		instancePublicKey := instanceCert.PublicKey
		result := reflect.DeepEqual(instancePublicKey, thatPublicKey)
		return result
	} else if instance.Type() == Trusted {
		thatName := cert.Subject.CommonName
		result := instance.name == thatName
		return result
	} else {
		return false
	}
}
