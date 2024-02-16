package thrift

import (
	"bytes"
	"testing"
)

// TestDecodeEmptyList checks that, when decoding a list, we get an empty slice
// rather than a nil. While a nil and an empty list are often hard to tell apart
// (e.g. len() gives the same result for them) they create problems downstream
// (e.g. when serializing a struct using the go json package).
//
// One may be lured to think that it's possible to hide the issue just by adding
// `json:omitempty` annotations to list fields (a common thing to do), so that the
// nil values don't become json nulls. This works only on the surface, as e.g.
// a list<list<i32>> clearly would have the problem - i.e. [[], [1,2]] would be
// serialized as [null, [1,2]].
func TestDecodeEmptyList(t *testing.T) {
	buf := &bytes.Buffer{}
	s := TestStruct{
		List: []string{},
	}
	if len(s.List) != 0 || s.List == nil {
		t.Fatal("Unexpected error")
	}
	err := EncodeStruct(NewBinaryProtocolWriter(buf, true), s)
	if err != nil {
		t.Fatal(err)
	}

	d := TestStruct{}
	err = DecodeStruct(NewBinaryProtocolReader(buf, true), &d)
	if err != nil {
		t.Fatal(err)
	}
	if len(d.List) != 0 {
		t.Fatalf("Decoding error - wrong length for d.List = %d (expected 0)", len(d.List))
	}
	if d.List == nil {
		t.Fatal("Decoding error - d.List is nil (expected [])")
	}
}
