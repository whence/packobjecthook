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
	"strings"
)

func runcmd(cmd *exec.Cmd, capturedDir string) error {
	if err := os.MkdirAll(capturedDir, 0755); err != nil {
		return err
	}

	// key
	fkey, err := os.Create(path.Join(capturedDir, "key"))
	if err != nil {
		return err
	}
	skey := sha256.New()
	defer fkey.Close()

	// cmd
	fcmd, err := os.Create(path.Join(capturedDir, "cmd"))
	if err != nil {
		return err
	}
	defer fcmd.Close()
	args := strings.TrimSpace(strings.Join(cmd.Args, " "))
	fmt.Fprintln(fcmd, args)
	fmt.Fprintln(skey, args)

	// stdin
	fstdin, err := os.Create(path.Join(capturedDir, "stdin"))
	if err != nil {
		return err
	}
	defer fstdin.Close()

	var bstdin bytes.Buffer
	io.Copy(io.MultiWriter(&bstdin, fstdin, skey), os.Stdin)
	cmd.Stdin = &bstdin
	key := hex.EncodeToString(skey.Sum(nil))
	fmt.Fprintln(fkey, key)

	// stdout
	fstdout, err := os.Create(path.Join(capturedDir, "stdout"))
	if err != nil {
		return err
	}
	defer fstdout.Close()
	cmd.Stdout = io.MultiWriter(os.Stdout, fstdout)

	// stderr
	fstderr, err := os.Create(path.Join(capturedDir, "stderr"))
	if err != nil {
		return err
	}
	defer fstderr.Close()
	cmd.Stderr = io.MultiWriter(os.Stderr, fstderr)

	// env
	fenv, err := os.Create(path.Join(capturedDir, "env"))
	if err != nil {
		return err
	}
	defer fenv.Close()
	for _, e := range os.Environ() {
		fmt.Fprintln(fenv, e)
	}

	return cmd.Run()
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
