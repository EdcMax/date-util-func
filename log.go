
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