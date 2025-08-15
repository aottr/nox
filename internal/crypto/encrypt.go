package crypto

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"filippo.io/age"
)

// EncryptFile encrypts the given file using the given identities
func EncryptFile(path string, recipients []age.Recipient) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return EncryptBytes(data, recipients)
}

func EncryptBytes(data []byte, recipients []age.Recipient) ([]byte, error) {

	src := bytes.NewReader(data)
	dst := new(bytes.Buffer)

	enc, err := age.Encrypt(dst, recipients...)
	if err != nil {
		return nil, fmt.Errorf("initializing age enc failed: %w", err)
	}

	if _, err := io.Copy(enc, src); err != nil {
		enc.Close()
		return nil, fmt.Errorf("encryption failed: %w", err)
	}

	if err := enc.Close(); err != nil {
		return nil, fmt.Errorf("finalizing encryption failed: %w", err)
	}

	return dst.Bytes(), nil
}
