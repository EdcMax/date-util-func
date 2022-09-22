package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"sync"
	"time"
)

// Goreman is RPC server
type Goreman struct {
	rpcChan chan<- *rpcMessage
}

type rpcMessage struct {
	Msg  string
	Args []string
	// sending error (if any) when the task completes
	ErrCh chan error
}

// Start do start
func (r *Goreman) Start(args []string, ret *string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	for _, arg := range args {
		if err = startProc(arg, nil, nil); err != nil {
			break
		}
	}
	return err
}

// Stop do stop
func (r *Goreman) Stop(args []string, ret *string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	errChan := make(chan error, 1)
	r.rpcChan <- &rpcMessage{
		Msg:   "stop",
		Args:  args,
		ErrCh: errChan,
	}
	err = <-errChan
	return
}

// StopAll do stop all
func (r *Goreman) StopAll(args []string, ret *string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	for _, proc := range procs {
		if err = stopProc(proc.name, nil); err != nil {
			break
		}
	}
	return err
}

// Restart do restart
func (r *Goreman) Restart(args []string, ret *string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	for _, arg := range args {
		if err = restartProc(arg); err != nil {
			break
		}
	}
	return err
}

// RestartAll do restart all
func (r *Goreman) RestartAll(args []string, ret *string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	for _, proc := range procs {
		if err = restartProc(proc.name); err != nil {
			break
		}
	}
	return err
}

// List do list
func (r *Goreman) List(args []string, ret *string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	*ret = ""
	for _, proc := range procs {
		*ret += proc.name + "\n"
	}
	return err
}

// Status do status
func (r *