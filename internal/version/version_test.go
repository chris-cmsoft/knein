package version

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCurrentPrefersBuildFlag(t *testing.T) {
	t.Cleanup(func() { version = "" })
	version = "v9.9.9"

	if got := Current(); got != "v9.9.9" {
		t.Fatalf("Current() = %q, want %q", got, "v9.9.9")
	}
}

func TestCurrentFallsBackToBuildInfo(t *testing.T) {
	if got := Current(); got == "" {
		t.Fatal("Current() is empty without a build flag")
	}
}

func TestLatestReadsTagName(t *testing.T) {
	var accept string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accept = r.Header.Get("Accept")
		_, _ = w.Write([]byte(`{"tag_name":"v1.2.3","name":"v1.2.3"}`))
	}))
	t.Cleanup(server.Close)
	useEndpoint(t, server.URL)

	got, err := Latest(context.Background())
	if err != nil {
		t.Fatalf("Latest returned error: %v", err)
	}
	if got != "v1.2.3" {
		t.Fatalf("Latest() = %q, want %q", got, "v1.2.3")
	}
	if accept != "application/vnd.github+json" {
		t.Fatalf("Accept header = %q, want %q", accept, "application/vnd.github+json")
	}
}

func TestLatestReportsMissingReleases(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(server.Close)
	useEndpoint(t, server.URL)

	_, err := Latest(context.Background())
	if err == nil {
		t.Fatal("Latest returned nil, want error")
	}
	if got, want := err.Error(), "no releases published yet"; got != want {
		t.Fatalf("error = %q, want %q", got, want)
	}
}

func TestLatestReportsUnexpectedStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	t.Cleanup(server.Close)
	useEndpoint(t, server.URL)

	_, err := Latest(context.Background())
	if err == nil {
		t.Fatal("Latest returned nil, want error")
	}
	if !strings.Contains(err.Error(), "403") {
		t.Fatalf("error = %q, want it to mention the status", err)
	}
}

func TestLatestRejectsResponseWithoutTag(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(server.Close)
	useEndpoint(t, server.URL)

	if _, err := Latest(context.Background()); err == nil {
		t.Fatal("Latest returned nil, want error")
	}
}

func TestWrite(t *testing.T) {
	tests := []struct {
		name      string
		current   string
		latest    string
		lookupErr error
		want      []string
	}{
		{
			name:    "up to date",
			current: "v1.2.3",
			latest:  "v1.2.3",
			want:    []string{"knein   v1.2.3\n", "latest  v1.2.3 (up to date)\n"},
		},
		{
			name:    "build metadata is still the same release",
			current: "v1.2.3+dirty",
			latest:  "v1.2.3",
			want:    []string{"knein   v1.2.3+dirty\n", "latest  v1.2.3 (up to date)\n"},
		},
		{
			name:    "newer release links the downloads",
			current: "v1.2.3",
			latest:  "v1.3.0",
			want:    []string{"knein   v1.2.3\n", "latest  v1.3.0 (" + releasesPage + ")\n"},
		},
		{
			name:      "failed lookup still reports the running version",
			current:   "v1.2.3",
			lookupErr: context.DeadlineExceeded,
			want:      []string{"knein   v1.2.3\n", "latest  unknown (context deadline exceeded)\n"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var out bytes.Buffer
			if err := write(&out, test.current, test.latest, test.lookupErr); err != nil {
				t.Fatalf("write returned error: %v", err)
			}
			if got, want := out.String(), strings.Join(test.want, ""); got != want {
				t.Fatalf("write() wrote %q, want %q", got, want)
			}
		})
	}
}

func useEndpoint(t *testing.T, url string) {
	t.Helper()
	previous := Endpoint
	t.Cleanup(func() { Endpoint = previous })
	Endpoint = url
}
