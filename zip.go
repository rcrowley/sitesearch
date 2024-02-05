//go:build !lambda

package main

import (
	"archive/zip"
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// Zip writes the starting-point zip file that contains the compiled Lambda
// function, named bootstrap per the Lambda provided.al2 runtime requirement,
// to disk, and adds the file with the given relative pathname to it.
func Zip(pathname string) (err error) {
	var (
		fIdx, fZip *os.File
		r          *zip.Reader
		rBin       io.ReadCloser
		w          *zip.Writer
		wBin, wIdx io.Writer
	)

	if r, err = zip.NewReader(bytes.NewReader(bootstrapZip), int64(len(bootstrapZip))); err != nil {
		log.Print(err)
		return
	}

	if fZip, err = os.Create("sitesearch.zip"); err != nil {
		return
	}
	w = zip.NewWriter(fZip)

	if wBin, err = w.Create("bootstrap"); err != nil {
		goto Error
	}
	for _, f := range r.File {
		if f.Name != "bootstrap" {
			panic(fmt.Sprintf(
				"there should not be any files except bootstrap in this zip file; found %s",
				f.Name,
			))
		}
		if rBin, err = f.Open(); err != nil {
			goto Error
		}
		if _, err = io.Copy(wBin, rBin); err != nil {
			goto Error
		}
		if err = rBin.Close(); err != nil {
			goto Error
		}
	}

	if fIdx, err = os.Open(pathname); err != nil {
		return
	}
	if wIdx, err = w.Create(filepath.Base(pathname)); err != nil {
		goto Error
	}
	if _, err = io.Copy(wIdx, fIdx); err != nil {
		goto Error
	}
	if err = fIdx.Close(); err != nil {
		goto Error
	}

Error:
	if err = w.Close(); err != nil {
		return
	}
	if err = fZip.Close(); err != nil {
		return
	}
	return
}

//go:generate env GOARCH=arm64 GOOS=linux go build -o bootstrap -tags lambda
//go:generate touch -t 202402040000.00 bootstrap
//go:generate zip -X bootstrap.zip bootstrap
//go:embed bootstrap.zip
var bootstrapZip []byte
