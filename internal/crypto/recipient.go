package crypto

import (
	"bytes"
	"fmt"
	"strings"

	"filippo.io/age"
)

func StringsToRecipients(strs []string) ([]age.Recipient, error) {
	var b bytes.Buffer
	for _, s := range strs {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		// age.ParseRecipients accepts multiple lines...
		b.WriteString(s)
		b.WriteByte('\n')
	}
	recips, err := age.ParseRecipients(&b)
	if err != nil {
		return nil, fmt.Errorf("parse recipients: %w", err)
	}
	return recips, nil
}
