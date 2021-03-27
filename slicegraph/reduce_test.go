package slicegraph_test

import (
	"testing"

	"github.com/jbszczepaniak/slicegraph/slicegraph"
	"github.com/stretchr/testify/assert"
)

func TestReduce(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		h, a := slicegraph.Reduce(nil)
		assert.Empty(t, h)
		assert.Empty(t, a)
	})

	t.Run("single empty", func(t *testing.T) {
		s := make(map[string][]int)
		s["a"] = []int{}
		h, a := slicegraph.Reduce(s)

		assert.Equal(t, slicegraph.Header{
			Pointer: "0x0",
			Len:     0,
			Cap:     0,
		}, h["a"])
		assert.Empty(t, a)
	})

	t.Run("single nil has nil pointer and no backing array", func(t *testing.T) {
		s := make(map[string][]int)
		s["a"] = nil
		h, a := slicegraph.Reduce(s)

		assert.Equal(t, slicegraph.Header{
			Pointer: "nil",
			Len:     0,
			Cap:     0,
		}, h["a"])
		assert.Empty(t, a)
	})

	t.Run("nil slice and with values slice have single backing array", func(t *testing.T) {
		s := make(map[string][]int)
		s["a"] = nil
		s["b"] = []int{1, 2}
		h, a := slicegraph.Reduce(s)

		assert.Equal(t, "nil", h["a"].Pointer)

		// pointer of b points to first element of first array.
		if assert.Len(t, a, 1) {
			assert.Equal(t, h["b"].Pointer, a[0].Addresses[0])
			assert.Equal(t, []string{"1", "2"}, a[0].Values)
		}
	})

	t.Run("two different slices with the same values have different backing arrays", func(t *testing.T) {
		s := make(map[string][]int)
		s["a"] = []int{1, 2}
		s["b"] = []int{1, 2}
		h, a := slicegraph.Reduce(s)
		assert.NotEqual(t, h["a"].Pointer, h["b"].Pointer)

		if assert.Len(t, a, 2) {
			addr0arr0 := a[0].Addresses[0]
			addr0arr1 := a[1].Addresses[0]

			order1 := (addr0arr0 == h["a"].Pointer && addr0arr1 == h["b"].Pointer)
			order2 := (addr0arr0 == h["b"].Pointer && addr0arr1 == h["a"].Pointer)
			// doesn't matter which one, but both pointers need to point to first element of different arrays.

			assert.True(t, order1 || order2,
				"%s and %s don't point to beginning of different backing arrays: [%v], [%v]",
				h["a"].Pointer, h["b"].Pointer, a[0].Addresses, a[1].Addresses)
		}

	})

	t.Run("two different slices pointing to the same point in one backing array", func(t *testing.T) {
		s := make(map[string][]int)
		slice := []int{1, 2}
		s["a"] = slice
		s["b"] = slice
		h, a := slicegraph.Reduce(s)

		if assert.Len(t, a, 1) {
			assert.Equal(t, h["a"].Pointer, a[0].Addresses[0])
			assert.Equal(t, h["b"].Pointer, a[0].Addresses[0])
			assert.Equal(t, []string{"1", "2"}, a[0].Values)
		}
	})

	t.Run("two different slices pointing to different points in one backing array", func(t *testing.T) {
		s := make(map[string][]int)
		s["a"] = []int{1, 2, 3, 4, 5}
		s["b"] = s["a"][3:]

		h, a := slicegraph.Reduce(s)

		if assert.Len(t, a, 1) {
			assert.Equal(t, h["a"].Pointer, a[0].Addresses[0])
			assert.Equal(t, h["b"].Pointer, a[0].Addresses[3])
			assert.Equal(t, []string{"1", "2", "3", "4", "5"}, a[0].Values)
		}
	})

	t.Run("headers have proper len and cap after reduction", func(t *testing.T) {
		testcases := []struct {
			slice []int
			len   int
			cap   int
		}{
			{nil, 0, 0},
			{[]int{}, 0, 0},
			{[]int{1, 2, 3, 4}, 4, 4},
		}

		for _, tc := range testcases {
			h, _ := slicegraph.Reduce(map[string][]int{
				"test": tc.slice,
			})
			assert.Equal(t, tc.len, h["test"].Len)
			assert.Equal(t, tc.cap, h["test"].Cap)
		}
	})
}
