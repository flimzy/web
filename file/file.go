package file

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/flimzy/web/blob"
	"time"
)

type File interface {
	blob.Blob
	LastModifiedDate() time.Time
	Name() string
}

type FileObject struct {
	blob.BlobObject
}

// Internalize internalizes a standard *js.Object to a File object
func Internalize(o *js.Object) *FileObject {
	return &FileObject{ blob.BlobObject{*o} }
}

var _ File = &FileObject{}

func (f *FileObject) LastModifiedDate() time.Time {
	return f.Get("lastModifiedDate").Interface().(time.Time)
}

func (f *FileObject) Name() string {
	return f.Get("name").String()
}
