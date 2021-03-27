package slicegraph

import (
	"fmt"
	"io"

	"github.com/goccy/go-graphviz"
)

// AsGraph takes slice of headers, and all backing arrays of the headers, and
// creates graphviz graph for them.
// TODO: cluster them and concentrate
// compound=True;
// size="8,15!";
//		rankdir=TB;
//		edge[weight=1.2];
//		node [shape=plaintext];
//
func AsGraph(slices map[string][]int, w io.Writer) error {
	headers, arrays := Reduce(slices)
	header := `digraph G {
		compound=True;
		concentrate=True;
		rankdir=LR;
		size="8,15!";
		node [shape=plaintext];
	`
	footer := `}`

	var sliceHeaders string
	for name, header := range headers {
		sliceHeaders += fmt.Sprintf(
			`header%s [label=<
				<TABLE BORDER="0" CELLBORDER="1" CELLSPACING="0">
				<TR><TD BGCOLOR="gray">%s</TD></TR>
				<TR><TD PORT="ptr">%s</TD></TR>
				<TR><TD PORT="len">len: %d</TD></TR>
				<TR><TD PORT="cap">cap: %d</TD></TR>
			</TABLE>>]
		`, name, name, header.Pointer, header.Len, header.Cap)
	}

	var backingArrays string
	arrayAddresses := make(map[string]string) // find an address

	for idx, a := range arrays {
		backingArrays += fmt.Sprintf(`array%d [label=<<TABLE BORDER="0" CELLBORDER="1" CELLSPACING="0">`, idx)

		firstRow := "<TR>"
		secondRow := "<TR>"
		for arrayIdx := 0; arrayIdx < len(a.Addresses); arrayIdx++ {
			// for addressIndex, address := range a {
			arrayAddresses[a.Addresses[arrayIdx]] = fmt.Sprintf("array%d:address%s", idx, a.Addresses[arrayIdx])
			firstRow += fmt.Sprintf(`<TD PORT="address%s" >%s</TD>`, a.Addresses[arrayIdx], a.Addresses[arrayIdx])
			secondRow += fmt.Sprintf(`<TD PORT="value%s" >%s</TD>`, a.Values[arrayIdx], a.Values[arrayIdx])
		}
		firstRow += "</TR>"
		secondRow += "</TR>"
		backingArrays += firstRow
		backingArrays += secondRow
		backingArrays += "</TABLE>>];\n\n"
	}

	var connections string
	for name, header := range headers {
		if header.Pointer == "nil" || header.Pointer == "0x0" {
			continue
		}
		connections += fmt.Sprintf(`header%s:ptr -> %s;
		`, name, arrayAddresses[header.Pointer])
	}

	parsed, err := graphviz.ParseBytes([]byte(header + sliceHeaders + backingArrays + connections + footer))
	if err != nil {
		return err
	}
	g := graphviz.New()
	return g.Render(parsed, graphviz.SVG, w)
}
