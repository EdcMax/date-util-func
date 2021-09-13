
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
	p := proc.cmd.Process
	if p == nil {
		return nil
	}

	pgid, err := unix.Getpgid(p.Pid)
	if err != nil {
		return err
	}

	// use pgid, ref: http://unix.stackexchange.com/questions/14815/process-descendants
	pid := p.Pid
	if pgid == p.Pid {
		pid = -1 * pid
	}

	target, err := os.FindProcess(pid)