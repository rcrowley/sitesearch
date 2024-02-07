package index

import (
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/rcrowley/mergician/html"
)

type (
	Index struct {
		idx bleve.Index
	}
	Result struct {
		bleve.SearchResult
	}
)

func Open(pathname string) (*Index, error) {
	idx, err := bleve.Open(pathname)
	if err != nil {
		if err != bleve.ErrorIndexPathDoesNotExist {
			return nil, err
		}
		mapping := bleve.NewIndexMapping()
		idx, err = bleve.New(pathname, mapping)
		if err != nil {
			return nil, err
		}
	}
	return &Index{idx}, nil
}

func OpenReadOnly(pathname string) (*Index, error) {
	idx, err := bleve.OpenUsing(pathname, map[string]interface{}{
		"read_only": true,
	})
	if err != nil {
		return nil, err
	}
	return &Index{idx}, nil
}

func (idx *Index) Close() error {
	return idx.idx.Close()
}

func (idx *Index) Index(pk string, data any) error {
	return idx.idx.Index(pk, data)
}

func (idx *Index) IndexHTML(pk string, n *html.Node) error {
	return idx.Index(pk, html.Text(n))
}

func (idx *Index) IndexHTMLFile(pathname string) error {
	n, err := html.ParseFile(pathname)
	if err != nil {
		return err
	}
	return idx.IndexHTML(pkForPathname(pathname), n)
}

func (idx *Index) IndexHTMLFiles(pathnames []string) error {
	for _, pathname := range pathnames {
		if err := idx.IndexHTMLFile(pathname); err != nil {
			return err
		}
	}
	return nil
}

func (idx *Index) Search(q string) (*Result, error) {
	sr, err := idx.idx.Search(bleve.NewSearchRequest(bleve.NewMatchQuery(q)))
	if err != nil {
		return nil, err
	}
	return &Result{*sr}, nil
}

func pkForPathname(pathname string) (pk string) {
	pk = pathname
	if strings.HasSuffix(pk, "/index.html") { // if it's a directory index...
		pk = strings.TrimSuffix(pk, "index.html") // ...strip the filename but keep the trailing '/'
	}
	if pk == "" || pk == "index.html" { // if if ends up empty or just "index.html"...
		pk = "/" // ...link to it as "/"
	}
	return
}
