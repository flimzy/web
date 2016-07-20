package worker

import (
	"github.com/gopherjs/gopherjs/js"
)

type Worker struct {
	o          *js.Object
	terminated bool
	fromWorker <-chan work
}

type work struct {
	err     error
	message interface{}
}

func New(file string) *Worker {
	worker := js.Global.Get("Worker").New(file)
	c := make(chan work)
	worker.Set("onerror", func(e *js.Error) {
		go func() {
			c <- work{err: e}
		}()
	})
	worker.Set("onmessage", func(e *js.Object) {
		go func() {
			c <- work{message: e.Get("data").Interface()}
		}()
	})
	return &Worker{
		o:          worker,
		fromWorker: c,
	}
}

func (w *Worker) Receive() (interface{}, error) {
	// First see if there's anything to fetch
	select {
	case msg := <-w.fromWorker:
		return msg.message, msg.err
	}
	// If we made it this far, it means the worker has terminated, and we've
	// already drained the channel, so just return nothing
	if w.terminated {
		return nil, nil
	}
	// Wait for something from the worker
	msg := <-w.fromWorker
	return msg.message, msg.err
}

func (w *Worker) Send(m interface{}) {
	w.o.Call("postMessage", m)
}

func (w Worker) Terminate() {
	w.o.Call("terminate")
}
