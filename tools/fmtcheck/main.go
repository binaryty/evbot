package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "fmtcheck: no files provided")
		os.Exit(1)
	}

	args := append([]string{"-l"}, os.Args[1:]...)

	cmd := exec.Command("gofmt", args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "fmtcheck: gofmt failed: %v\n", err)
		os.Exit(1)
	}

	out := strings.TrimSpace(stdout.String())
	if out == "" {
		fmt.Println("gofmt check passed")
		return
	}

	fmt.Fprintln(os.Stderr, "fmtcheck: the following files need gofmt:")
	for _, file := range strings.Split(out, "\n") {
		fmt.Fprintln(os.Stderr, file)
	}
	os.Exit(1)
}

