package main

import (
	"archive/zip"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/rcrowley/sitesearch/index"
)

func TestZip(t *testing.T) {
	const idxFilename = "sitesearch.idx"

	// Known-good sitesearch.idx to add to the zip file.
	must(os.RemoveAll(idxFilename))
	idx := must2(index.Open(idxFilename))
	must(idx.Close())
	defer os.RemoveAll(idxFilename)

	zipFilename, err := Zip(idxFilename)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(zipFilename)

	r, err := zip.OpenReader(zipFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	var count int
	for _, f := range r.File {
		fi := f.FileInfo()
		switch f.Name { // f.Name is the relative pathname; fi.Name() is the basename
		case "bootstrap":
			count++
			if fi.Size() < 1000000 {
				t.Fatalf("%s is suspiciously small (%d bytes)", f.Name, fi.Size())
			}
		case fmt.Sprintf("%s/", idxFilename):
			count++
		case fmt.Sprintf("%s/index_meta.json", idxFilename):
			count++
			if fi.Size() != 42 {
				t.Fatalf("%s should contain 42 bytes but contains %d bytes", f.Name, fi.Size())
			}
		case fmt.Sprintf("%s/store/", idxFilename):
			count++
		case fmt.Sprintf("%s/store/root.bolt", idxFilename):
			count++
			if fi.Size() != 65536 {
				t.Fatalf("%s should contain 65536 bytes but contains %d bytes", f.Name, fi.Size())
			}
		default:
			t.Fatal(f.Name)
		}
	}
	if count != 5 {
		t.Fatalf("only found %d out of the expected 5 files in %s", count, zipFilename)
	}

	/*
		stdout, err := exec.Command("unzip", "-l", zipFilename).Output()
		if err != nil {
			t.Fatal(err)
		}
		os.Stdout.Write(stdout)
	*/

}
