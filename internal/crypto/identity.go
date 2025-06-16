package crypto

import (
	"bytes"
	"fmt"
	"os"

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
