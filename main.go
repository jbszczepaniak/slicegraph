package main

import (
	"bytes"
	"fmt"

	"github.com/jbszczepaniak/slicegraph/slicegraph"
)

func main() {
	a := []int{1, 2, 3, 4}
	x := append(a, []int{5, 6, 7, 8}...)
	c := make([]int, 3, 4)
	d := make([]int, len(x))
	e := x[2:3]
	copy(d, c)

	var b []byte
	buff := bytes.NewBuffer(b)

	err := slicegraph.AsGraph(map[string][]int{
		"a": a,
		"x": x,
		"c": c,
		"d": d,
		"e": e,
	}, buff)

	if err != nil {
		panic(err)
	}

	fmt.Println(buff)
}
