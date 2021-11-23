package main

import (
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sys/windows"
)

var cmdStart = []string{"cmd", "/c"}
var procAttrs = &windows.SysProcAttr{
	CreationFlags: windows.CREATE_UNICODE_ENVIRONMENT | windows.CREATE_NEW_PROCESS_GROUP,
}

func terminateProc(proc *procInfo, _ os.Signal) error {
	dll, err := windows.LoadDLL("kernel32.dll")
	if err != nil {
		return err
	}
	defer dll.Release()

	pid := proc.cmd.Process.Pid

	f, err := dll.FindProc("AttachConsole")
	if 