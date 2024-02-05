package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rcrowley/mergician/html"
	"github.com/rcrowley/sitesearch/index"
	"golang.org/x/net/html/atom"
)

const IdxFilename = "sitesearch.idx"

func Search(q string) (*html.Node, error) {

	idx, err := index.Open(IdxFilename)
	if err != nil {
		return nil, err
	}
	defer idx.Close()

	result, err := idx.Search(q)

	n := &html.Node{
		DataAtom: atom.Pre,
		Data:     "pre",
		Type:     html.ElementNode,
	}
	n.AppendChild(&html.Node{
		Data: fmt.Sprint(result),
		Type: html.TextNode,
	})
	return n, nil
}

func SearchHandler(ctx context.Context, req events.LambdaFunctionURLRequest) (resp events.LambdaFunctionURLResponse, err error) {
	panic("not implemented")
}
