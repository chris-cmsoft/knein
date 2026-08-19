package cmd

import (
	"errors"
	"io"
	"path/filepath"
	"testing"

	kubecontext "github.com/chris-cmsoft/gotool-kubecontext-picker"
)

func runRootCommand(t *testing.T, args ...string) error {
	t.Helper()

	command := newRootCommand()
	command.SetArgs(args)
	command.SetOut(io.Discard)
	command.SetErr(io.Discard)

	return command.Execute()
}

func TestRootRejectsLimitsBelowOne(t *testing.T) {
	err := runRootCommand(t, "--limit", "0")

	if !errors.Is(err, kubecontext.ErrInvalidLimit) {
		t.Fatalf("error = %v, want %v", err, kubecontext.ErrInvalidLimit)
	}
}

func TestRootRejectsArguments(t *testing.T) {
	if err := runRootCommand(t, "prd-euw1"); err == nil {
		t.Fatal("error = nil, want error")
	}
}

func TestRootFailsForUnreadableKubeconfig(t *testing.T) {
	t.Setenv("KUBECONFIG", "")

	err := runRootCommand(t, "--kubeconfig", filepath.Join(t.TempDir(), "missing"))
	if err == nil {
		t.Fatal("error = nil, want error")
	}
}
