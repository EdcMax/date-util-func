
package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

// spawnProc starts the specified proc, and returns any error from running it.
func spawnProc(name string, errCh chan<- error) {
	proc := findProc(name)
	logger := createLogger(name, proc.colorIndex)

	cs := append(cmdStart, proc.cmdline)
	cmd := exec.Command(cs[0], cs[1:]...)
	cmd.Stdin = nil
	cmd.Stdout = logger
	cmd.Stderr = logger
	cmd.SysProcAttr = procAttrs

	if proc.setPort {
		cmd.Env = append(os.Environ(), fmt.Sprintf("PORT=%d", proc.port))
		fmt.Fprintf(logger, "Starting %s on port %d\n", name, proc.port)
	}
	if err := cmd.Start(); err != nil {
		select {
		case errCh <- err:
		default:
		}
		fmt.Fprintf(logger, "Failed to start %s: %s\n", name, err)