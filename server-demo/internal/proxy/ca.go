package proxy

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

type CAFiles struct {
	Dir      string
	CertPath string
	KeyPath  string
}

func DefaultCAFiles() CAFiles {
	// 相对运行目录创建 certs 目录
	dir := filepath.Clean("../certs")
	return CAFiles{
		Dir:      dir,
		CertPath: filepath.Join(dir, "proxy_root_ca.pem"),
		KeyPath:  filepath.Join(dir, "proxy_root_ca.key"),
	}
}

func EnsureCAExists() (certPEM, keyPEM []byte, files CAFiles, err error) {
	files = DefaultCAFiles()
	if err = os.MkdirAll(files.Dir, 0755); err != nil {
		return nil, nil, files, err
	}
	certPEM, _ = os.ReadFile(files.CertPath)
	keyPEM, _ = os.ReadFile(files.KeyPath)
	if len(certPEM) > 0 && len(keyPEM) > 0 {
		return certPEM, keyPEM, files, nil
	}
	return GenerateCA()
}

func GenerateCA() (certPEM, keyPEM []byte, files CAFiles, err error) {
	files = DefaultCAFiles()
	if err = os.MkdirAll(files.Dir, 0755); err != nil {
		return nil, nil, files, err
	}

	// 生成 RSA 私钥
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, files, err
	}

	// 生成自签名根证书
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, serialNumberLimit)
	tmpl := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"CetiProbe Local Proxy Root CA"},
			Country:      []string{"CN"},
			CommonName:   "CetiProbe Root CA",
		},
		NotBefore:             time.Now().Add(-10 * time.Minute),
		NotAfter:              time.Now().AddDate(5, 0, 0),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            2,
		// 添加更好的兼容性设置
		SubjectKeyId:   []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		AuthorityKeyId: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, &priv.PublicKey, priv)
	if err != nil {
		return nil, nil, files, err
	}

	certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyPEM = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	if err = os.WriteFile(files.CertPath, certPEM, 0644); err != nil {
		return nil, nil, files, err
	}
	if err = os.WriteFile(files.KeyPath, keyPEM, 0600); err != nil {
		return nil, nil, files, err
	}
	return certPEM, keyPEM, files, nil
}

func LoadCAFromDisk() (certPEM, keyPEM []byte, files CAFiles, err error) {
	files = DefaultCAFiles()
	certPEM, err = os.ReadFile(files.CertPath)
	if err != nil {
		return nil, nil, files, err
	}
	keyPEM, err = os.ReadFile(files.KeyPath)
	if err != nil {
		return nil, nil, files, err
	}
	if len(certPEM) == 0 || len(keyPEM) == 0 {
		return nil, nil, files, errors.New("empty CA files")
	}
	return certPEM, keyPEM, files, nil
}

// ParseCA 将PEM解析为x509根证书与私钥
func ParseCA(certPEM, keyPEM []byte) (*x509.Certificate, *rsa.PrivateKey, error) {
	block, _ := pem.Decode(certPEM)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, nil, errors.New("invalid ca cert pem")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, nil, err
	}

	kblk, _ := pem.Decode(keyPEM)
	if kblk == nil {
		return nil, nil, errors.New("invalid ca key pem")
	}
	var key *rsa.PrivateKey
	if kblk.Type == "RSA PRIVATE KEY" {
		key, err = x509.ParsePKCS1PrivateKey(kblk.Bytes)
		if err != nil {
			return nil, nil, err
		}
	} else if kblk.Type == "PRIVATE KEY" {
		pk, err2 := x509.ParsePKCS8PrivateKey(kblk.Bytes)
		if err2 != nil {
			return nil, nil, err2
		}
		var ok bool
		key, ok = pk.(*rsa.PrivateKey)
		if !ok {
			return nil, nil, errors.New("unsupported private key type")
		}
	} else {
		return nil, nil, errors.New("unsupported key pem type")
	}
	return cert, key, nil
}

// SignHostCert 使用根CA为指定主机名签发叶子证书（含 SAN）
func SignHostCert(caCert *x509.Certificate, caKey *rsa.PrivateKey, host string, validFor time.Duration) (tlsCertPEM, tlsKeyPEM []byte, err error) {
	// 生成叶子私钥
	leafKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, serialNumberLimit)
	tmpl := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      pkix.Name{CommonName: host, Organization: []string{"CetiProbe MITM"}},
		NotBefore:    time.Now().Add(-5 * time.Minute),
		NotAfter:     time.Now().Add(validFor),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		DNSNames:     []string{host},
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, caCert, &leafKey.PublicKey, caKey)
	if err != nil {
		return nil, nil, err
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(leafKey)})
	return certPEM, keyPEM, nil
}
