package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rcrowley/mergician/html"
	"github.com/rcrowley/sitesearch/index"
	"golang.org/x/net/html/atom"
)

const (
	IdxFilename  = "sitesearch.idx"
	TmplFilename = "sitesearch.html"
)

func Search(q string) (*html.Node, error) {

	idx, err := index.OpenReadOnly(IdxFilename)
	if err != nil {
		return nil, err
	}
	defer idx.Close()

	result, err := idx.Search(q)

	body := &html.Node{
		DataAtom: atom.Body,
		Data:     "body",
		Type:     html.ElementNode,
	}
	ol := &html.Node{
		Attr: []html.Attribute{
			{Key: "class", Val: "sitesearch"},
		},
		DataAtom: atom.Ol,
		Data:     "ol",
		Type:     html.ElementNode,
	}
	body.AppendChild(newlineIndent(1))
	body.AppendChild(form(q))
	body.AppendChild(newlineIndent(1))
	body.AppendChild(ol)
	body.AppendChild(newlineIndent(1))
	body.AppendChild(form(q))
	body.AppendChild(newline())

	for _, hit := range result.Hits {
		li := &html.Node{
			DataAtom: atom.Li,
			Data:     "li",
			Type:     html.ElementNode,
		}
		ol.AppendChild(newlineIndent(2))
		ol.AppendChild(li)
		h3 := &html.Node{
			DataAtom: atom.H3,
			Data:     "h3",
			Type:     html.ElementNode,
		}
		li.AppendChild(newlineIndent(3))
		li.AppendChild(h3)
		a := &html.Node{
			Attr: []html.Attribute{
				{Key: "href", Val: hit.ID},
			},
			DataAtom: atom.A,
			Data:     "a",
			Type:     html.ElementNode,
		}
		h3.AppendChild(a)
		a.AppendChild(&html.Node{
			Data: "TODO title",
			Type: html.TextNode,
		})
		a.AppendChild(&html.Node{
			DataAtom: atom.Br,
			Data:     "br",
			Type:     html.ElementNode,
		})
		kbd := &html.Node{
			DataAtom: atom.Kbd,
			Data:     "kbd",
			Type:     html.ElementNode,
		}
		a.AppendChild(kbd)
		kbd.AppendChild(&html.Node{
			Data: hit.ID,
			Type: html.TextNode,
		})
		p := &html.Node{
			DataAtom: atom.P,
			Data:     "p",
			Type:     html.ElementNode,
		}
		li.AppendChild(newlineIndent(3))
		li.AppendChild(p)
		li.AppendChild(newlineIndent(2))
		p.AppendChild(&html.Node{
			Data: "TODO some kind of summary text of the result",
			Type: html.TextNode,
		})
	}
	ol.AppendChild(newlineIndent(1))

	return body, nil
}

func SearchHandler(ctx context.Context, req events.LambdaFunctionURLRequest) (resp events.LambdaFunctionURLResponse, err error) {
	resp.StatusCode = http.StatusOK
	resp.Headers = map[string]string{"Content-Type": "text/html; charset=utf-8"}
	var n, tmpl *html.Node

	if q := req.QueryStringParameters["q"]; q == "" {
		n = &html.Node{
			DataAtom: atom.Body,
			Type:     html.ElementNode,
		}
		n.AppendChild(newlineIndent(1))
		n.AppendChild(form(q))
	} else {
		n, err = Search(q)
	}
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		n = errorNode(err)
	}

	if tmpl, err = html.ParseFile(TmplFilename); err == nil {
		if n, err = html.Merge([]*html.Node{tmpl, n}, html.DefaultRules()); err != nil {
			return
		}
	}

	resp.Body = html.String(n)
	return
}

func errorNode(err error) *html.Node {
	body := &html.Node{
		DataAtom: atom.Body,
		Data:     "body",
		Type:     html.ElementNode,
	}
	pre := &html.Node{
		Attr: []html.Attribute{
			{Key: "class", Val: "sitesearch"},
		},
		DataAtom: atom.Pre,
		Data:     "pre",
		Type:     html.ElementNode,
	}
	body.AppendChild(newlineIndent(1))
	body.AppendChild(pre)
	body.AppendChild(newline())
	pre.AppendChild(&html.Node{
		Data: err.Error(),
		Type: html.TextNode,
	})
	return body
}

func form(q string) *html.Node {
	form := &html.Node{
		Attr: []html.Attribute{
			{Key: "class", Val: "sitesearch"},
		},
		DataAtom: atom.Form,
		Data:     "form",
		Type:     html.ElementNode,
	}
	form.AppendChild(newlineIndent(2))
	form.AppendChild(&html.Node{
		Attr: []html.Attribute{
			{Key: "name", Val: "q"},
			{Key: "placeholder", Val: "query"},
			{Key: "type", Val: "text"},
			{Key: "value", Val: q},
		},
		DataAtom: atom.Input,
		Data:     "input",
		Type:     html.ElementNode,
	})
	form.AppendChild(newlineIndent(2))
	form.AppendChild(&html.Node{
		Attr: []html.Attribute{
			{Key: "type", Val: "submit"},
			{Key: "value", Val: "Search"},
		},
		DataAtom: atom.Input,
		Data:     "input",
		Type:     html.ElementNode,
	})
	form.AppendChild(newlineIndent(1))
	return form
}

func indent(tabs int) *html.Node {
	return &html.Node{
		Data: strings.Repeat("    ", tabs),
		Type: html.TextNode,
	}
}

func newline() *html.Node {
	return &html.Node{
		Data: "\n",
		Type: html.TextNode,
	}
}

func newlineIndent(tabs int) *html.Node {
	return &html.Node{
		Data: "\n" + strings.Repeat("    ", tabs),
		Type: html.TextNode,
	}
}
