package file

import (
	"github.com/flimzy/web/blob"
	"github.com/gopherjs/gopherjs/js"
	"time"
)

type File interface {
	blob.Blob
	LastModifiedDate() time.Time
	Name() string
}

type file struct {
	js.Object
}

func (f *file) LastModifiedDate() time.Time {
	return f.Get("lastModifiedDate").Interface().(time.Time)
}

func (f *file) Name() string {
	return f.Get("name").String()
}
