// Package version reports the running knein version and the latest release.
package version

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

// version is set by release builds with -ldflags
// "-X github.com/chris-cmsoft/knein/internal/version.version=v1.2.3".
// Builds without it fall back to the module version Go records in the binary.
var version string

const (
	repository    = "chris-cmsoft/knein"
	releasesPage  = "https://github.com/" + repository + "/releases/latest"
	lookupTimeout = 5 * time.Second
)

// Endpoint is the GitHub API URL queried for the latest release. It is a
// variable so tests can point it at a local server.
var Endpoint = "https://api.github.com/repos/" + repository + "/releases/latest"

// Current reports the version of the running binary.
func Current() string {
	if version != "" {
		return version
	}

	info, ok := debug.ReadBuildInfo()
	if !ok || info.Main.Version == "" {
		return "unknown"
	}
	return info.Main.Version
}

// Latest returns the tag of the most recent published release.
func Latest(ctx context.Context) (string, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, Endpoint, nil)
	if err != nil {
		return "", err
	}
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	request.Header.Set("User-Agent", "knein/"+Current())

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	switch response.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return "", errors.New("no releases published yet")
	default:
		return "", fmt.Errorf("GitHub returned %s", response.Status)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&release); err != nil {
		return "", fmt.Errorf("read GitHub response: %w", err)
	}
	if release.TagName == "" {
		return "", errors.New("GitHub response contained no tag name")
	}

	return release.TagName, nil
}

// Report writes the running version and the latest release to out. A
// failed lookup is reported in place of the latest version instead of failing
// the command, so --version still answers the question it can answer offline.
func Report(ctx context.Context, out io.Writer) error {
	ctx, cancel := context.WithTimeout(ctx, lookupTimeout)
	defer cancel()

	latest, lookupErr := Latest(ctx)
	return write(out, Current(), latest, lookupErr)
}

// sameRelease reports whether current is a build of the latest release. Build
// metadata such as the "+dirty" Go records for a modified checkout does not
// change which release a build came from.
func sameRelease(current string, latest string) bool {
	base, _, _ := strings.Cut(current, "+")
	return base == latest
}

func write(out io.Writer, current string, latest string, lookupErr error) error {
	if _, err := fmt.Fprintf(out, "%-7s %s\n", "knein", current); err != nil {
		return err
	}

	switch {
	case lookupErr != nil:
		_, err := fmt.Fprintf(out, "%-7s unknown (%v)\n", "latest", lookupErr)
		return err
	case sameRelease(current, latest):
		_, err := fmt.Fprintf(out, "%-7s %s (up to date)\n", "latest", latest)
		return err
	default:
		_, err := fmt.Fprintf(out, "%-7s %s (%s)\n", "latest", latest, releasesPage)
		return err
	}
}
