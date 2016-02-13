package filereader

import (
    "github.com/gopherjs/gopherjs/js"
	"github.com/flimzy/event"
)

type FileReader struct {
    js.Object
}

type ReadyState int

const (
    EMPTY ReadyState = 0
    LOADING ReadyState = 1
    DONE ReadyState = 2
)

func New() *FileReader {
	fr := js.Global.Get("FileReader").New()
	return &FileReader{*fr}
}

func (fr *FileReader) Error() string {
    return fr.Get("error").String()
}

func (fr *FileReader) ReadyState() ReadyState {
    return ReadyState( fr.Get("readyState").Int() )
}

func (fr *FileReader) IsEmpty() bool {
    return fr.ReadyState() == EMPTY
}

func (fr *FileReader) IsLoading() bool {
    return fr.ReadyState() == LOADING
}

func (fr *FileReader) IsDone() bool {
    return fr.ReadyState() == DONE
}

type progressUpdate struct {
    Event    event.Event
    FileName string
}

type ProgressEvent interface {
	event.Event
}

type ProgressFeed func() (ProgressEvent, string)

func (fr *FileReader) Progress() ProgressFeed {
	progChan := make(chan progressUpdate)
	fr.Set("onabort", func(e *js.Object) {
		progChan <- progressUpdate{
			Event: &event.BasicEvent{*e},
		}
	})
	fr.Set("onprogress", func(e *js.Object, fn string) {
		progChan <- progressUpdate{
			Event: &event.BasicEvent{*e},
			FileName: fn,
		}
		
	})
	return ProgressFeed(func() (ProgressEvent, string) {
		update := <-progChan
		event := update.Event
		return event, update.FileName
	})
}
