package main

import (
	"os"
	"strings"
	"testing"

	"github.com/rcrowley/mergician/html"
	"github.com/rcrowley/sitesearch/index"
)

func TestSearch(t *testing.T) {
	must(os.Chdir("testdata"))
	defer func() { must(os.Chdir("..")) }()

	must(os.RemoveAll(IdxFilename))
	idx := must2(index.Open(IdxFilename))
	must(idx.IndexHTMLFile("../index/testdata/test.html", nil))
	must(idx.Close())
	defer os.RemoveAll(IdxFilename)

	n, err := Search("cool")
	if err != nil {
		t.Fatal(err)
	}

	s := html.String(n)
	if !strings.Contains(s, "/index/testdata/test.html") {
		t.Fatal(s)
	}

}
