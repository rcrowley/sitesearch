package main

import (
	"archive/zip"
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/rcrowley/sitesearch/index"
)

func TestZip(t *testing.T) {

	// Known-good sitesearch.idx to add to the zip file.
	must(os.RemoveAll(IdxFilename))
	idx := must2(index.Open(IdxFilename))
	must(idx.Close())
	defer os.RemoveAll(IdxFilename)

	zipFile, err := Zip(ZipFilename, IdxFilename, "index/test.html")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(ZipFilename)
	if len(zipFile) < 1000000 {
		t.Fatalf("zipFile is suspiciously small (%d bytes)", len(zipFile))
	}

	r, err := zip.OpenReader(ZipFilename)
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
		case fmt.Sprintf("%s/", IdxFilename):
			count++
		case fmt.Sprintf("%s/index_meta.json", IdxFilename):
			count++
			if fi.Size() != 42 {
				t.Fatalf("%s should contain 42 bytes but contains %d bytes", f.Name, fi.Size())
			}
		case fmt.Sprintf("%s/store/", IdxFilename):
			count++
		case fmt.Sprintf("%s/store/root.bolt", IdxFilename):
			count++
			if fi.Size() != 65536 {
				t.Fatalf("%s should contain 65536 bytes but contains %d bytes", f.Name, fi.Size())
			}
		case fmt.Sprintf("%s", "test.html"): // not "index/test.html", because `zip -j`
			count++
		default:
			t.Fatal(f.Name)
		}
	}
	if count != 6 {
		stdout, err := exec.Command("unzip", "-l", ZipFilename).Output()
		if err != nil {
			t.Fatal(err)
		}
		os.Stdout.Write(stdout)
		t.Fatalf("only found %d out of the expected 6 files in %s", count, ZipFilename)
	}

}
