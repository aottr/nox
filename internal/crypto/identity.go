package crypto

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"filippo.io/age"
)

// LoadAgeIdentities reads and parses all age identities from the given file
func LoadAgeIdentities(path string) ([]age.Identity, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read age key file: %w", err)
	}

	identities, err := age.ParseIdentities(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("invalid age identity file: %w", err)
	}

	return identities, nil
}

func LoadAgeIdentitiesFromPaths(paths []string) ([]age.Identity, error) {
	var allIdentities []age.Identity

	for _, path := range paths {
		idents, err := LoadAgeIdentities(path)
		if err != nil {
			return nil, fmt.Errorf("failed to load identities from %s: %w", path, err)
		}
		allIdentities = append(allIdentities, idents...)
	}

	return allIdentities, nil
}

func GenerateIdentity(path string) (priv string, pub string, err error) {
	id, err := age.GenerateX25519Identity()
	if err != nil {
		return "", "", fmt.Errorf("generate identity: %w", err)
	}

	priv = id.String()            // "AGE-SECRET-KEY-1..."
	pub = id.Recipient().String() // "age1..."

	// ensure directory exists and apply permissions
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return "", "", fmt.Errorf("mkdir %s: %w", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(priv+"\n"), 0o600); err != nil {
		return "", "", fmt.Errorf("write identity: %w", err)
	}

	// write public key
	if err := os.WriteFile(path+".pub", []byte(pub+"\n"), 0o644); err != nil {
		return "", "", fmt.Errorf("write identity: %w", err)
	}
	return priv, pub, nil
}
