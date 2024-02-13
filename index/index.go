package index

import (
	"path"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/rcrowley/mergician/html"
)

const (
	Title   = "Title"
	Summary = "Summary"
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
		im := bleve.NewIndexMapping()
		im.DefaultMapping.AddFieldMappingsAt(Title, bleve.NewTextFieldMapping())
		im.DefaultMapping.AddFieldMappingsAt(Summary, bleve.NewTextFieldMapping())
		idx, err = bleve.New(pathname, im)
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
	return idx.Index(pk, struct {
		Title, Summary string
		Text           html.TextOnlyNode
	}{
		Title:   html.Title(n),
		Summary: html.FirstParagraph(n),
		Text:    html.Text(n),
	})
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
	req := bleve.NewSearchRequest(bleve.NewMatchQuery(q))
	req.Fields = []string{Title, Summary}

	// This is a small-scale search engine. 1,000 results should pretty much
	// always be all of the results. And even if it's not, who's going to scroll
	// through more than 1,000 search results?
	req.Size = 1000

	result, err := idx.idx.Search(req)
	if err != nil {
		return nil, err
	}
	return &Result{*result}, nil
}

func pkForPathname(pathname string) (pk string) {
	pk = path.Join("/", pathname)
	if strings.HasSuffix(pk, "/index.html") { // if it's a directory index...
		pk = strings.TrimSuffix(pk, "index.html") // ...strip the filename but keep the trailing '/'
	}
	return
}
