package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func exportUpstart(cfg *config, path string) error {
	for i, proc := range procs {
		f, err := os.Create(filepath.Join(path, "app-"+proc.name+".conf"))
		if err != nil {
			return err
		}

		fmt.Fprintf(f, "start on starting app-%s\n", proc.name)
		fmt.Fprintf(f, "stop on stopping app-%s\n", proc.name)
		fmt.Fprintf(f, "respawn\n")
		fmt.Fprintf(f, "\n")

		env := map[string]string{}
		procfi