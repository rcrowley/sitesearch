//go:build !lambda

package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCLI(t *testing.T) {}

// TestPublishedPaths checks that only currently-published documents are kept
// for indexing.
func TestPublishedPaths(t *testing.T) {
	dir := t.TempDir()
	write := func(name, body string) string {
		p := filepath.Join(dir, name)
		doc := "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n<title>x</title>\n</head>\n<body>\n" + body + "\n</body>\n</html>\n"
		if err := os.WriteFile(p, []byte(doc), 0666); err != nil {
			t.Fatal(err)
		}
		return p
	}

	live := write("live.html", `<article class="body"><time class="published" datetime="2026-01-01 00:00:00"></time><h1>Live</h1></article>`)
	scheduled := write("scheduled.html", `<article class="body"><time class="published" datetime="2026-12-31 00:00:00"></time><h1>Soon</h1></article>`)
	draft := write("draft.html", `<article class="body draft"><h1>Draft</h1></article>`)

	now := time.Date(2026, 6, 4, 0, 0, 0, 0, time.UTC)
	got := publishedPaths([]string{live, scheduled, draft}, now)
	if len(got) != 1 || got[0] != live {
		t.Errorf("publishedPaths = %v, want just %s", got, live)
	}
}
