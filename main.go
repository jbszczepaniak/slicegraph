package main

import (
	"bytes"
	"fmt"

	"github.com/jbszczepaniak/slicegraph/slicegraph"
)

func main() {
	var a []int
	b := []int{1, 2, 3, 4, 5, 6, 7, 8}
	c := b[4:6]
	d := b[6:]
	e := make([]int, len(b))
	copy(e, b)

	var by []byte
	buff := bytes.NewBuffer(by)
	err := slicegraph.AsGraph(map[string][]int{
		"a": a,
		"b": b,
		"c": c,
		"d": d,
		"e": e,
	}, buff)
	if err != nil {
		panic(err)
	}

	fmt.Println(buff) // redirect this to the file and open it.
}
