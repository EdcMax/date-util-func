
package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

var sleep string

func TestMain(m *testing.M) {
	var dir string
	var err error
	sleep, err = exec.LookPath("sleep")
	if err != nil {
		if runtime.GOOS != "windows" {
			panic(err)
		}

		code := `package main;import ("os";"strconv";"time");func main(){i,_:=strconv.ParseFloat(os.Args[1]);time.Sleep(time.Duration(i)*time.Second)}`
		dir, err := os.MkdirTemp("", "goreman-test")
		if err != nil {
			panic(err)
		}
		sleep = filepath.Join(dir, "sleep.exe")
		src := filepath.Join(dir, "sleep.go")
		err = os.WriteFile(src, []byte(code), 0644)
		if err != nil {
			panic(err)
		}
		b, err := exec.Command("go", "build", "-o", sleep, src).CombinedOutput()
		if err != nil {
			panic(string(b))
		}
		oldpath := os.Getenv("PATH")
		os.Setenv("PATH", dir+";"+oldpath)
		defer os.Setenv("PATH", oldpath)
	}
	r := m.Run()

	if dir != "" {
		os.RemoveAll(dir)
	}
	os.Exit(r)
}

func startGoreman(ctx context.Context, t *testing.T, ch <-chan os.Signal, file []byte) error {
	t.Helper()
	f, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.Write(file); err != nil {
		t.Fatal(err)
	}
	cfg := &config{
		ExitOnError: true,