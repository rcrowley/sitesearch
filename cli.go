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

	idx := must2(index.Open(IdxFilename))
	defer idx.Close()

	must(idx.IndexHTMLFiles(flag.Args()))

	must(Zip(ZipFilename, IdxFilename, *tmpl))

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
