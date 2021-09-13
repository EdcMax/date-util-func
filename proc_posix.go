
//go:build !windows
// +build !windows

package main

import (
	"os"
	"os/signal"

	"golang.org/x/sys/unix"
)

const sigint = unix.SIGINT
const sigterm = unix.SIGTERM
const sighup = unix.SIGHUP

var cmdStart = []string{"/bin/sh", "-c"}
var procAttrs = &unix.SysProcAttr{Setpgid: true}

func terminateProc(proc *procInfo, signal os.Signal) error {