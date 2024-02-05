package main

import (
	"os"
	"strings"
	"testing"

	"github.com/rcrowley/mergician/html"
	"github.com/rcrowley/sitesearch/index"
)

func TestSearch(t *testing.T) {
	must(os.RemoveAll(IdxFilename))
	idx := must2(index.Open(IdxFilename))
	must(idx.IndexHTMLFile("index/test.html"))
	must(idx.Close())
	defer os.RemoveAll(IdxFilename)

	n, err := Search("cool")
	if err != nil {
		t.Fatal(err)
	}

	s := html.String(n)
	if !strings.Contains(s, "1 matches") {
		t.Fatal(s)
	}
	if !strings.Contains(s, "1. index/test.html") {
		t.Fatal(s)
	}

}
