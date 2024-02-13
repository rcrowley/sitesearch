package index

import (
	"os"
	"testing"

	"github.com/rcrowley/mergician/html"
)

func TestIndex(t *testing.T) {
	idx := setup(t)
	defer teardown(t, idx)

	idx.Index("/test.html", "cool")

	testSearch(t, idx)
}

func TestIndexHTML(t *testing.T) {
	idx := setup(t)
	defer teardown(t, idx)

	n, err := html.ParseFile("test.html")
	if err != nil {
		t.Fatal(err)
	}
	idx.IndexHTML("/test.html", n)

	testSearch(t, idx)
}

// TestIndexGoFile shows that this will work just fine, albeit the entire
// contents of the file will end up being one big TextOnlyNode.
func TestIndexGoFile(t *testing.T) {
	idx := setup(t)
	defer teardown(t, idx)

	if err := idx.IndexHTMLFile("index_test.go"); err != nil {
		t.Fatal(err)
	}

	result, err := idx.Search("cool")
	if err != nil {
		t.Fatal(err)
	}
	if result.Hits[0].ID != "/index_test.go" {
		t.Fatal(result)
	}
	//t.Log(result)
}

func TestIndexHTMLFile(t *testing.T) {
	idx := setup(t)
	defer teardown(t, idx)

	if err := idx.IndexHTMLFile("test.html"); err != nil {
		t.Fatal(err)
	}

	testSearch(t, idx)
}

func TestIndexHTMLFiles(t *testing.T) {
	idx := setup(t)
	defer teardown(t, idx)

	if err := idx.IndexHTMLFiles([]string{"test.html"}); err != nil {
		t.Fatal(err)
	}

	testSearch(t, idx)
}

func TestPKForPathname(t *testing.T) {
	for pathname, pk := range map[string]string{
		"foo/bar.html":       "/foo/bar.html",
		"foo/bar/index.html": "/foo/bar/",
		"foo.html":           "/foo.html",
		"foo/index.html":     "/foo/",
		"index.html":         "/",
	} {
		if pkForPathname(pathname) != pk {
			t.Fatalf("pathname %q produced pk %q instead of %q", pathname, pkForPathname(pathname), pk)
		}
	}
}

func setup(t *testing.T) *Index {
	if err := os.RemoveAll("test.idx"); err != nil {
		t.Fatal(err)
	}
	idx, err := Open("test.idx")
	if err != nil {
		t.Fatal(err)
	}
	return idx
}

func teardown(t *testing.T, idx *Index) {
	if err := idx.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll("test.idx"); err != nil {
		t.Fatal(err)
	}
}

func testSearch(t *testing.T, idx *Index) {
	result, err := idx.Search("cool")
	if err != nil {
		t.Fatal(err)
	}
	if result.Hits[0].ID != "/test.html" {
		t.Fatal(result)
	}
	if title, ok := result.Hits[0].Fields["Title"].(string); ok && title != "My cool webpage" {
		t.Fatal(result.Hits[0].Fields)
	}
	if summary, ok := result.Hits[0].Fields["Summary"].(string); ok && summary != "Stuff" {
		t.Fatal(result.Hits[0].Fields)
	}
	//t.Log(result)
}
