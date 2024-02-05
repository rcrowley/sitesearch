//go:build !lambda

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/rcrowley/sitesearch/index"
)

func init() {
	log.SetFlags(0)
}

func main() {
	//output := flag.String("o", "-", "write to this file instead of standard output")
	// TODO -n Lambda function name, no -o (?), -r AWS region
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: sitesearch [-o <output>] <input>[...]
  -o <output>   write to this file instead of standard output
  <input>[...]  pathname to one or more input HTML or Markdown files
`)
	}
	flag.Parse()

	const idxFilename = "sitesearch.idx"

	idx := must2(index.Open(idxFilename))
	defer idx.Close()

	must(idx.IndexHTMLFiles(flag.Args()))

	zipFilename := must2(Zip(idxFilename))

	log.Print(zipFilename) // TODO upload zipFilename to Lambda and set a function URL

}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func must2[T any](v T, err error) T {
	must(err)
	return v
}
