package html

import (
	"bytes"
	"io"

	"golang.org/x/net/html"
)

// Render is almost an alias for x/net/html's Render function but ensures
// files end with a trailing '\n' character.
func Render(w io.Writer, n *Node) error {
	if err := html.Render(w, n); err != nil {
		return err
	}
	_, err := w.Write([]byte{'\n'})
	return err
}

// String renders the *Node to a string and returns it. In case of error,
// the return value is the error string instead. If handling this error
// is important to you, use Render instead.
func String(n *Node) string {
	var b bytes.Buffer
	err := html.Render(&b, n)
	if err != nil {
		return err.Error()
	}
	return b.String()
}
