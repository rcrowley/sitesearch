package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rcrowley/sitesearch/index"
)

func TestCopyRecursive(t *testing.T) {
	must(os.RemoveAll(IdxFilename))
	idx := must2(index.Open(IdxFilename))
	must(idx.IndexHTMLFile("index/test.html"))
	must(idx.Close())
	defer os.RemoveAll(IdxFilename)

	if err := CopyRecursive(IdxFilename, filepath.Join("/tmp", IdxFilename)); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(filepath.Join("/tmp", IdxFilename))

	if _, err := os.Stat(filepath.Join("/tmp", IdxFilename, "index_meta.json")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join("/tmp", IdxFilename, "store/root.bolt")); err != nil {
		t.Fatal(err)
	}
}
