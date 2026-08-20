package cmd

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	kubecontext "github.com/chris-cmsoft/gotool-kubecontext-picker"
	"github.com/chris-cmsoft/knein/internal/version"
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

func TestRootCommandReportsVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"tag_name":"v1.2.3"}`))
	}))
	t.Cleanup(server.Close)

	previous := version.Endpoint
	t.Cleanup(func() { version.Endpoint = previous })
	version.Endpoint = server.URL

	var stdout bytes.Buffer
	command := newRootCommand()
	command.SetArgs([]string{"--version"})
	command.SetOut(&stdout)
	command.SetErr(io.Discard)

	if err := command.Execute(); err != nil {
		t.Fatalf("--version returned error: %v", err)
	}

	output := stdout.String()
	if !strings.HasPrefix(output, "knein   ") {
		t.Fatalf("output = %q, want it to start with the binary name", output)
	}
	if !strings.Contains(output, "latest  v1.2.3") {
		t.Fatalf("output = %q, want it to report the latest release", output)
	}
}
