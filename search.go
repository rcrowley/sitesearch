package main

import (
	"context"
	"fmt"
	"net/http"

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
	pre := &html.Node{
		DataAtom: atom.Pre,
		Data:     "pre",
		Type:     html.ElementNode,
	}
	body.AppendChild(pre)
	pre.AppendChild(&html.Node{
		Data: fmt.Sprint(result),
		Type: html.TextNode,
	})
	return body, nil
}

func SearchHandler(ctx context.Context, req events.LambdaFunctionURLRequest) (resp events.LambdaFunctionURLResponse, err error) {
	resp.Headers = make(map[string]string)
	var n *html.Node
	if n, err = Search(req.QueryStringParameters["q"]); err == nil {
		var tmpl *html.Node
		if tmpl, err = html.ParseFile(TmplFilename); err == nil {
			if n, err = html.Merge([]*html.Node{tmpl, n}, html.DefaultRules()); err != nil {
				return
			}
		}
		resp.StatusCode = http.StatusOK
		resp.Headers["Content-Type"] = "text/html; charset=utf-8"
		resp.Body = html.String(n)
	} else {
		resp.StatusCode = http.StatusInternalServerError
		resp.Headers["Content-Type"] = "text/plain; charset=utf-8"
		resp.Body = err.Error()
	}
	return
}
