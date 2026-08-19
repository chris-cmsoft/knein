// Package k9s starts k9s for a Kubernetes context.
package k9s

import (
	"os"
	"os/exec"
)

// Open starts k9s for context in the current terminal.
func Open(context string) error {
	command := exec.Command("k9s", "--context="+context)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}
