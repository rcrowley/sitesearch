//go:build !lambda

package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
)

const (
	TmplFilename = "sitesearch.html" // relative pathname within the zip file only
	ZipFilename  = "sitesearch.zip"
)

// Zip writes the starting-point zip file that contains the compiled Lambda
// function, named bootstrap per the Lambda provided.al* runtime requirement,
// to disk, and adds the file with the given relative pathname to it.
func Zip(
	zipPathname string, // output
	idxPathname, tmplPathname string, // input
) error {

	f, err := os.Create(zipPathname)
	if err != nil {
		return err
	}
	if _, err := f.Write(bootstrapZip); err != nil {
		if err2 := f.Close(); err2 != nil {
			return fmt.Errorf("%w %w", err, err2)
		}
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	if err := exec.Command("zip", "-X", "-r", zipPathname, idxPathname).Run(); err != nil {
		return err
	}

	if err := exec.Command("zip", "-X", zipPathname, tmplPathname).Run(); err != nil { // TODO put tmplPathname at TmplFilename (sitesearch.html)
		return err
	}

	return nil
}

//go:generate env GOARCH=arm64 GOOS=linux go build -o bootstrap -tags lambda
//go:generate touch -t 202402040000.00 bootstrap
//go:generate zip -X bootstrap.zip bootstrap
//go:embed bootstrap.zip
var bootstrapZip []byte
