
package main

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/mattn/go-colorable"
)

type clogger struct {
	idx     int
	name    string
	writes  chan []byte
	done    chan struct{}
	timeout time.Duration // how long to wait before printing partial lines
	buffers buffers       // partial lines awaiting printing
}

var colors = []int{
	32, // green
	36, // cyan
	35, // magenta
	33, // yellow
	34, // blue
	31, // red
}
var mutex = new(sync.Mutex)

var out = colorable.NewColorableStdout()

type buffers [][]byte

func (v *buffers) consume(n int64) {
	for len(*v) > 0 {
		ln0 := int64(len((*v)[0]))
		if ln0 > n {
			(*v)[0] = (*v)[0][n:]
			return
		}
		n -= ln0
		*v = (*v)[1:]
	}
}

func (v *buffers) WriteTo(w io.Writer) (n int64, err error) {