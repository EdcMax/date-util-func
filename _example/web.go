//go:build ignore
// +build ignore

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

var addr = flag.String("a", ":8080", "address")

func 