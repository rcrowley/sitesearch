package main

import (
	"archive/zip"
	"log"
	"os"
	"testing"

	"github.com/rcrowley/sitesearch/index"
)

func TestZip(t *testing.T) {

	// Known-good sitesearch.idx to add to the zip file.
	must(os.RemoveAll("sitesearch.idx"))
	idx := must2(index.Open("sitesearch.idx"))
	must(idx.Close())
	defer os.RemoveAll("sitesearch.idx")

	if err := Zip("sitesearch.idx"); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("sitesearch.zip")

	r, err := zip.OpenReader("sitesearch.zip")
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
				t.Fatalf("%s is suspiciously small (%d bytes)", fi.Name(), fi.Size())
			}
		case "sitesearch.idx/":
			count++
		case "sitesearch.idx/index_meta.json":
			count++
			if fi.Size() != 42 {
				t.Fatalf("%s should contain 42 bytes but contains %d bytes", fi.Name(), fi.Size())
			}
		case "sitesearch.idx/store/":
			count++
		case "sitesearch.idx/store/root.bolt":
			count++
			if fi.Size() != 65536 {
				t.Fatalf("%s should contain 65536 bytes but contains %d bytes", fi.Name(), fi.Size())
			}
		default:
			t.Fatal(fi.Name())
		}
	}
	if count != 5 {
		t.Fatalf("only found %d out of the expected 5 files in sitesearch.zip", count)
	}

	/*
		stdout, err := exec.Command("unzip", "-l", "sitesearch.zip").Output()
		if err != nil {
			t.Fatal(err)
		}
		os.Stdout.Write(stdout)
	*/

}
