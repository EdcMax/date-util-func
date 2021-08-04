
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
		return
	}
	proc.cmd = cmd
	proc.stoppedBySupervisor = false
	proc.mu.Unlock()
	err := cmd.Wait()
	proc.mu.Lock()
	proc.cond.Broadcast()
	if err != nil && !proc.stoppedBySupervisor {
		select {
		case errCh <- err:
		default:
		}
	}
	proc.waitErr = err
	proc.cmd = nil
	fmt.Fprintf(logger, "Terminating %s\n", name)
}

// Stop the specified proc, issuing os.Kill if it does not terminate within 10
// seconds. If signal is nil, os.Interrupt is used.
func stopProc(name string, signal os.Signal) error {
	if signal == nil {
		signal = os.Interrupt
	}
	proc := findProc(name)
	if proc == nil {
		return errors.New("unknown proc: " + name)
	}

	proc.mu.Lock()
	defer proc.mu.Unlock()

	if proc.cmd == nil {
		return nil
	}
	proc.stoppedBySupervisor = true

	err := terminateProc(proc, signal)
	if err != nil {
		return err
	}

	timeout := time.AfterFunc(10*time.Second, func() {
		proc.mu.Lock()
		defer proc.mu.Unlock()
		if proc.cmd != nil {
			err = killProc(proc.cmd.Process)
		}
	})