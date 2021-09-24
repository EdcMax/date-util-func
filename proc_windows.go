package main

import (
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sys/windows"
)

var cmdStart = []st