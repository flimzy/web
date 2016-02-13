/*
Package blob provides GopherJS bindings for the JavaScript BlobObject objects.

Read more about JavaScript BlobObjects here: https://developer.mozilla.org/en-US/docs/Web/API/BlobObject
*/
package blob

import (
	"github.com/gopherjs/gopherjs/js"
	"sync"
)

type Blob interface {
	IsClosed() bool
	Size() int64
	Type() string
	Close()
	Slice(int, int, string) Blob
	Bytes() []byte
}

var _ Blob = &BlobObject{}

// BlobObject wraps a js.Object
type BlobObject struct {
	js.Object
}

type Options struct {
	Type    string `js:"type"`
	Endings string `js:"endings"`
}

// New returns a newly created BlobObject object whose content consists of the
// concatenation of the array of values given in parameter.
func New(parts []interface{}, opts Options) *BlobObject {
	blob := js.Global.Get("BlobObject").New(parts, opts)
	return &BlobObject{*blob}
}

// Internalize internalizes a standard *js.Object to a GlobObj
func Internalize(o *js.Object) *BlobObject {
	return &BlobObject{*o}
}

// IsClosed returns true if the Close() method (or the underlying JavaScript
// BlobObjectl.close() method) has been called on the blob. Closed blobs can not be
// read.
func (b *BlobObject) IsClosed() bool {
	return b.Get("isClosed").Bool()
}

// Size returns the size, in bytes, of the data contained in the BlobObject object.
func (b *BlobObject) Size() int64 {
	return b.Get("size").Int64()
}

// Type returns a string indicating the MIME type of the data contained in
// the BlobObject. If the type is unknown, this string is empty.
func (b *BlobObject) Type() string {
	return b.Get("type").String()
}

// Close closes the blob object, possibly freeing underlying resources.
func (b *BlobObject) Close() {
	b.Call("close")
}

// Slice returns a new BlobObject object containing the specified range of bytes of the source BlobObject.
func (b *BlobObject) Slice(start, end int, contenttype string) Blob {
	newBlobObject := b.Call("slice", start, end, contenttype)
	return &BlobObject{*newBlobObject}
}

// Bytes returns a slice of the contents of the BlobObject.
func (b *BlobObject) Bytes() []byte {
	fileReader := js.Global.Get("FileReader").New()
	var wg sync.WaitGroup
	var buf []byte
	wg.Add(1)
	fileReader.Set("onload", js.MakeFunc(func(this *js.Object, _ []*js.Object) interface{} {
		defer wg.Done()
		buf = js.Global.Get("Uint8Array").New(this.Get("result")).Interface().([]uint8)
		return nil
	}))
	fileReader.Call("readAsArrayBuffer", b)
	wg.Wait()
	return buf
}
