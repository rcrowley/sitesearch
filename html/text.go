package html

import (
	"encoding/json"
	"strings"

	"golang.org/x/net/html"
)

// TextNode could be structured as a tree of linked lists, like the DOM, but
// if we do that then the somewhat naive indexing libraries won't be able to
// follow. So we have to structure the tree using slices, which reflection,
// JSON encoding, etc. can follow.
type TextNode struct {
	Nodes []TextNode
	Text  string
}

func Text(in *Node) (out TextNode) {
	for i := in.FirstChild; i != nil; i = i.NextSibling {
		o := Text(i)
		if len(o.Nodes) > 0 || o.Text != "" {
			out.Nodes = append(out.Nodes, o)
		}
	}
	if in.Type == html.TextNode && strings.TrimSpace(in.Data) != "" {
		out.Text = in.Data
	}
	return
}

func (n TextNode) String() string {
	b, err := json.MarshalIndent(n, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(b)
}
