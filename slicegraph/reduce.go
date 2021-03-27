package slicegraph

import (
	"fmt"
	"sort"
)

// BackingArray is an array that stits behind the slice and
// consists of slice of addresses as well as value representations
// that sit in these addresses.
type BackingArray struct {
	Addresses []string
	Values    []string
}

// Header represents a slice with a name, that points to backing array,
// has length and capacity. It's a model of go's slice.
type Header struct {
	Pointer string
	Len     int
	Cap     int
}

// Reduce takes many slices as variablename -> slice map, finds
// common backing arrays and returns only unique backing arrays
// and all slice headers pointing to these arrays.
func Reduce(slices map[string][]int) (map[string]Header, []BackingArray) {
	if len(slices) == 0 {
		return nil, nil
	}
	var representations []sliceRepresentation
	for k, v := range slices {
		representations = append(representations, reprSingle(k, v))
	}
	return reduce(representations)
}

type sliceRepresentation struct {
	array  *BackingArray //optional as there might not be any backing array
	header Header
	name   string
}

func reprSingle(name string, slice []int) sliceRepresentation {
	if slice == nil {
		return sliceRepresentation{
			header: Header{
				Pointer: "nil",
				Len:     0,
				Cap:     0,
			},
			name: name,
		}
	}
	if len(slice) == 0 {
		return sliceRepresentation{
			header: Header{
				Pointer: "0x0",
				Len:     0,
				Cap:     0,
			},
			name: name,
		}
	}
	sr := sliceRepresentation{
		header: Header{
			Pointer: fmt.Sprintf("%p", &slice[0]),
			Len:     len(slice),
			Cap:     cap(slice),
		},
		name: name,
	}
	sr.array = &BackingArray{} // initialization because of ptr to it. (because of optionality)
	for i := range slice {
		sr.array.Addresses = append(sr.array.Addresses, fmt.Sprintf("%p", &slice[i]))
		sr.array.Values = append(sr.array.Values, fmt.Sprintf("%d", slice[i]))
	}

	return sr
}

func reduce(sr []sliceRepresentation) (map[string]Header, []BackingArray) {
	var arrays []BackingArray
	headers := make(map[string]Header)

	for _, slice := range sr {
		if slice.array != nil {
			arrays = append(arrays, *slice.array)
		}
		headers[slice.name] = slice.header
	}

	// no backing arrays, maybe only empty/nil slices.
	if len(arrays) == 0 {
		return headers, arrays
	}

	// sort arrays from longest to shortest.
	sort.Slice(arrays, func(i, j int) bool {
		return len(arrays[i].Addresses) > len(arrays[j].Addresses)
	})

	var result []BackingArray
	result = append(result, arrays[0])
	arrays = arrays[1:]

	// iterate over arrays, and figure out, whether next array
	// is new one (then attach to result), or maybe is it just
	// subarray other array (ignore it).
	for _, array := range arrays {
		var indeedSubslice bool

		for _, r := range result {
			if isSubslice(array.Addresses, r.Addresses) {
				indeedSubslice = true
				break // it is sublslice, we can just ignore it
			}
		}
		if !indeedSubslice {
			result = append(result, array)
		}

	}
	return headers, result
}

func isSubslice(s1 []string, s2 []string) bool {
	if len(s1) > len(s2) {
		return false
	}
	for _, e := range s1 {
		if !contains(s2, e) {
			return false
		}
	}
	return true
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
