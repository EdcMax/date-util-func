package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func exportUpstart(cfg *config, path string) error {
	for i, proc := range procs {
		f, err := os.Create(filepat