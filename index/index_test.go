package index

import (
	"os"
	"testing"

	"github.com/rcrowley/sitesearch/html"
)

func TestIndex(t *testing.T) {
	idx := setup(t)
	defer teardown(t, idx)

	idx.Index("../html/test.html", "cool")

	testSearch(t, idx)
}

func TestIndexHTML(t *testing.T) {
	idx := setup(t)
	defer teardown(t, idx)

	n, err := html.ParseFile("../html/test.html")
	if err != nil {
		t.Fatal(err)
	}
	idx.IndexHTML("../html/test.html", n)

	testSearch(t, idx)
}

func TestIndexHTMLFile(t *testing.T) {
	idx := setup(t)
	defer teardown(t, idx)

	if err := idx.IndexHTMLFile("../html/test.html"); err != nil {
		t.Fatal(err)
	}

	testSearch(t, idx)
}

func TestIndexHTMLFiles(t *testing.T) {
	idx := setup(t)
	defer teardown(t, idx)

	if err := idx.IndexHTMLFiles([]string{"../html/test.html"}); err != nil {
		t.Fatal(err)
	}

	testSearch(t, idx)
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
	if result.Hits[0].ID != "../html/test.html" {
		t.Fatal(result)
	}
	//t.Log(result)
}
