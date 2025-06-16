package crypto

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"filippo.io/age"
)

func DecryptAgeFile(inputPath, outputPath, identityPath string) error {
	identityFile, err := os.ReadFile(identityPath)
	if err != nil {
		return fmt.Errorf("failed to read identity: %w", err)
	}

	identities, err := age.ParseIdentities(bytes.NewReader(identityFile))
	if err != nil {
		return fmt.Errorf("failed to parse identities: %w", err)
	}

	in, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open encrypted file: %w", err)
	}
	defer in.Close()

	r, err := age.Decrypt(in, identities...)
	if err != nil {
		return fmt.Errorf("failed to decrypt: %w", err)
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, r); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}

	return nil
}

func DecryptFile(inputPath string, identities []age.Identity) ([]byte, error) {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read encrypted file: %w", err)
	}
	return DecryptBytes(data, identities)
}

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
