package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestZip(t *testing.T) {
	if err := os.WriteFile("sitesearch.idx", []byte("not a real index, just testing zip file I/O"), 0666); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("sitesearch.idx")
	if err := Zip("sitesearch.idx"); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("sitesearch.zip")
	stdout, err := exec.Command("unzip", "-l", "sitesearch.zip").Output()
	if err != nil {
		t.Fatal(err)
	}
	if string(stdout) != `Archive:  sitesearch.zip
  Length      Date    Time    Name
---------  ---------- -----   ----
  1951131  00-00-1980 00:00   bootstrap
       43  00-00-1980 00:00   sitesearch.idx
---------                     -------
  1951174                     2 files
` {
		t.Fatal(string(stdout))
	}
}
