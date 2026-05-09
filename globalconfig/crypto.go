package globalconfig

import (
	"archive/zip"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"strings"
)

const (
	encryptedConfigFileName = "config.json.enc"
	signatureFileName       = "signature.sig"
	versionFileName         = "version"
)

type PackageFiles struct {
	EncryptedConfig string
	Signature       string
	Version         string
}

func PackageHash(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func DecodePackage(data []byte, privateKey string) (*ConfigPackage, error) {
	files, err := ExtractPackageFiles(data)
	if err != nil {
		return nil, err
	}

	plainText, err := DecryptConfig(files.EncryptedConfig, privateKey)
	if err != nil {
		return nil, err
	}

	config, err := parseConfigPackage(plainText)
	if err != nil {
		return nil, err
	}
	if files.Version != "" && strings.TrimSpace(files.Version) != config.ConfigVersion {
		return nil, fmt.Errorf("package version mismatch: version file %s, config json %s", strings.TrimSpace(files.Version), config.ConfigVersion)
	}
	return config, nil
}

func parseConfigPackage(plainText string) (*ConfigPackage, error) {
	var config ConfigPackage
	if err := json.Unmarshal([]byte(plainText), &config); err != nil {
		return nil, fmt.Errorf("parse config json failed: %w", err)
	}
	config.RawJSON = plainText
	return &config, nil
}

func ExtractPackageFiles(data []byte) (*PackageFiles, error) {
	if len(data) > 4 && data[0] == 0x50 && data[1] == 0x4B && data[2] == 0x03 && data[3] == 0x04 {
		return extractPackageFilesFromZip(data)
	}
	return &PackageFiles{EncryptedConfig: string(data)}, nil
}

func extractPackageFilesFromZip(data []byte) (*PackageFiles, error) {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}

	files := &PackageFiles{}
	for _, file := range reader.File {
		switch file.Name {
		case encryptedConfigFileName, signatureFileName, versionFileName:
		default:
			continue
		}
		content, err := readZipFile(file)
		if err != nil {
			return nil, err
		}
		switch file.Name {
		case encryptedConfigFileName:
			files.EncryptedConfig = string(content)
		case signatureFileName:
			files.Signature = string(content)
		case versionFileName:
			files.Version = string(content)
		}
	}

	if strings.TrimSpace(files.EncryptedConfig) == "" {
		return nil, fmt.Errorf("zip package missing %s", encryptedConfigFileName)
	}
	return files, nil
}

func readZipFile(file *zip.File) ([]byte, error) {
	rc, err := file.Open()
	if err != nil {
		return nil, err
	}
	content, readErr := io.ReadAll(rc)
	closeErr := rc.Close()
	if readErr != nil {
		return nil, readErr
	}
	if closeErr != nil {
		return nil, closeErr
	}
	return content, nil
}

func DecryptConfig(encContent string, privateKey string) (string, error) {
	parts := strings.Split(strings.TrimSpace(encContent), "|")
	if len(parts) != 4 {
		return "", fmt.Errorf("invalid encrypted config format, want encryptedAesKey|cipherText|nonce|tag")
	}

	encryptedAesKey, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return "", fmt.Errorf("decode encryptedAesKey failed: %w", err)
	}
	cipherText, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("decode cipherText failed: %w", err)
	}
	nonce, err := base64.StdEncoding.DecodeString(parts[2])
	if err != nil {
		return "", fmt.Errorf("decode nonce failed: %w", err)
	}
	tag, err := base64.StdEncoding.DecodeString(parts[3])
	if err != nil {
		return "", fmt.Errorf("decode tag failed: %w", err)
	}

	rsaPrivKey, err := parseRSAPrivateKey(privateKey)
	if err != nil {
		return "", err
	}

	aesKeyBytes, err := rsa.DecryptPKCS1v15(nil, rsaPrivKey, encryptedAesKey)
	if err != nil {
		return "", fmt.Errorf("rsa decrypt aes key failed: %w", err)
	}

	aesKey, err := base64.StdEncoding.DecodeString(string(aesKeyBytes))
	if err != nil {
		aesKey = aesKeyBytes
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", fmt.Errorf("create aes cipher failed: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("create gcm failed: %w", err)
	}

	ciphertextWithTag := make([]byte, 0, len(cipherText)+len(tag))
	ciphertextWithTag = append(ciphertextWithTag, cipherText...)
	ciphertextWithTag = append(ciphertextWithTag, tag...)

	plainBytes, err := gcm.Open(nil, nonce, ciphertextWithTag, nil)
	if err != nil {
		return "", fmt.Errorf("aes-gcm decrypt failed: %w", err)
	}

	return string(plainBytes), nil
}

func parseRSAPrivateKey(privateKeyPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(strings.TrimSpace(privateKeyPEM)))
	if block == nil {
		return nil, fmt.Errorf("private key pem decode failed")
	}

	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse private key failed: %w", err)
	}
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not rsa")
	}
	return rsaKey, nil
}
