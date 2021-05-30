
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