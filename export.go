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
		procfile, err := filepath.Abs(cfg.Procfile)
		if err != nil {
			return err
		}
		b, err := os.ReadFile(filepath.Join(filepath.Dir(procfile), ".env"))
		if err == nil {
			for _, line := range strings.Split(string(b), "\n") {
				token := strings.SplitN(line, "=", 2)
				if len(token) != 2 {
					continue
				}
				if strings.HasPrefix(token[0], "export ") {
					token[0] = token[0][7:]
				}
				token[0] = st