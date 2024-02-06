//go:build !lambda

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/rcrowley/sitesearch/index"
)

func init() {
	log.SetFlags(0)
}

func main() {
	// TODO -n Lambda function name, -r AWS region
	tmpl := flag.String("t", "", "HTML template for search result pages")
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: sitesearch -t <template> <input>[...]
  -t <template>  HTML template for search result pages
  <input>[...]   pathname to one or more input HTML or Markdown files
`)
	}
	flag.Parse()
	if *tmpl == "" {
		log.Fatal("-t <template> is required")
	}

	tmp, err := os.MkdirTemp("", "sitesearch-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmp)
	log.Print(tmp)

	// Index all the HTML we've been told to. Store the index where the Lambda
	// function is eventually going to look for it.
	// TODO also read pathnames from standard input
	idx := must2(index.Open(filepath.Join(tmp, IdxFilename)))
	must(idx.IndexHTMLFiles(flag.Args()))
	must(idx.Close())

	// Copy the search engine result page template to where the Lambda function
	// is eventually going to look for it.
	must(os.WriteFile(
		filepath.Join(tmp, TmplFilename),
		must2(os.ReadFile(*tmpl)),
		0666,
	))

	// Package up the application, index, and template for service in Lambda.
	oldpwd := must(os.Getwd())
	must(os.Chdir(tmp))
	must(Zip(ZipFilename, IdxFilename, TmplFilename))
	must(os.Chdir(oldpwd))

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
