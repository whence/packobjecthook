package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestCommandExec(t *testing.T) {
	dir, err := ioutil.TempDir("", "captured")
	if err != nil {
		t.FailNow()
	}
	t.Cleanup(func() {
		os.RemoveAll(dir)
	})

	cmd := exec.Command("./main", "sed", "-e", "s/world/git/g")
	cmd.Env = append(os.Environ(), "CAPTURED_DIR="+dir)
	t.Logf("setting CAPTURED_DIR to %s", dir)
	cmd.Stdin = strings.NewReader("hello world")

	var bstdout, bstderr bytes.Buffer
	cmd.Stdout = &bstdout
	cmd.Stderr = &bstderr

	err = cmd.Run()
	sstdout := bstdout.String()
	sstderr := bstderr.String()

	if err != nil {
		t.Logf("failed with error: %v %s %s", err, sstderr, sstdout)
		t.FailNow()
	}

	if sstderr != "" {
		t.Logf("failed with stderr: %s", sstderr)
		t.FailNow()
	}

	expected := "hello git"
	if sstdout != "hello git" {
		t.Logf("failed with stdout: %s expected: %s", sstdout, expected)
		t.FailNow()
	}
	t.Logf("OK: output %s", sstdout)
}
