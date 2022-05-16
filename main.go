package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

// WriteCounter counts the number of bytes written to it.
type WriteCounter struct {
	Total int64 // Total # of bytes written
}

// Write implements the io.Writer interface.
//
// Always completes and never returns an error.
func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += int64(n)
	return n, nil
}

func runcmd(cmd *exec.Cmd, capturedDir string) error {
	skey := sha256.New()

	args := strings.TrimSpace(strings.Join(cmd.Args, " "))
	fmt.Fprintln(skey, args)

	var bstdin1, bstdin2 bytes.Buffer
	io.Copy(io.MultiWriter(&bstdin1, &bstdin2, skey), os.Stdin)
	cmd.Stdin = &bstdin1
	stdin := strings.ReplaceAll(bstdin2.String(), "\n", " ")

	key := hex.EncodeToString(skey.Sum(nil))

	cstdout := &WriteCounter{}
	cmd.Stdout = io.MultiWriter(os.Stdout, cstdout)

	cstderr := &WriteCounter{}
	cmd.Stderr = io.MultiWriter(os.Stderr, cstderr)

	exitCode := "?"
	runError := cmd.Run()
	if runError != nil {
		if exitError, ok := runError.(*exec.ExitError); ok {
			exitCode = strconv.Itoa(exitError.ExitCode())
		}
	} else {
		exitCode = "0"
	}

	if f, err := os.OpenFile(path.Join(capturedDir, key), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		fmt.Fprintf(f, "%s %s | %s out=%d err=%d exit=%s\n", time.Now().Format(time.RFC3339), args, stdin, cstdout.Total, cstderr.Total, exitCode)
		defer f.Close()
	}

	return runError
}

func main() {
	// dont print out any extra information
	log.SetFlags(0)

	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatalln("No command to execute.")
	}

	capturedDir, ok := os.LookupEnv("CAPTURED_DIR")
	if !ok {
		log.Fatalln("CAPTURED_DIR not set.")
	}

	if err := os.MkdirAll(capturedDir, 0755); err != nil {
		log.Fatalf("Failed to create CAPTURED_DIR (%s).\n", capturedDir)
	}

	var exitCode int
	defer func() {
		recover()
		os.Exit(exitCode)
	}()

	cmd := exec.Command(args[0], args[1:]...)
	err := runcmd(cmd, capturedDir)
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			log.Printf("Failed to execute %s. Error: %v\n", cmd, err)
			exitCode = 1
		}
	} else {
		exitCode = 0
	}
}
