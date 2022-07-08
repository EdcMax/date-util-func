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
		if r := recov