package crypto

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"filippo.io/age"
)

// DecryptFile decrypts the given file using the given identities
func DecryptFile(inputPath string, identities []age.Identity) ([]byte, error) {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read encrypted file: %w", err)
	}
	return DecryptBytes(data, identities)
}

// DecryptBytes decrypts the given bytes using the given identities
func DecryptBytes(encrypted []byte, identities []age.Identity) ([]byte, error) {

	dec, err := age.Decrypt(bytes.NewReader(encrypted), identities...)
	if err != nil {
		return nil, fmt.Errorf("age decryption failed: %w", err)
	}

	var out bytes.Buffer
	if _, err := io.Copy(&out, dec); err != nil {
		return nil, fmt.Errorf("failed to read decrypted data: %w", err)
	}

	return out.Bytes(), nil
}
