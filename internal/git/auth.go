package git

import (
	"os"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func GetAuth() (transport.AuthMethod, error) {

	if sshKey, exists := os.LookupEnv("NOX_GIT_SSH_KEY_FILE"); exists {
		pemBytes, err := os.ReadFile(sshKey)
		if err != nil {
			return nil, err
		}
		passPhrase := ""
		if sshPass, exists := os.LookupEnv("NOX_GIT_SSH_KEY_PASSWORD"); exists {
			passPhrase = sshPass
		}

		pubKey, err := ssh.NewPublicKeys("git", pemBytes, passPhrase)
		if err != nil {
			return nil, err
		}

		return pubKey, nil
	}

	if gitToken, exists := os.LookupEnv("NOX_GIT_TOKEN"); exists {
		return &http.BasicAuth{
			Username: "nox",
			Password: gitToken,
		}, nil
	}
	return nil, nil
}
