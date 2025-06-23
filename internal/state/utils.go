package state

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// HashContent returns the SHA256 hash of the given data.
func HashContent(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func GenerateKey(appName, file string) string {
	return fmt.Sprintf("%s:%s", appName, file)
}
