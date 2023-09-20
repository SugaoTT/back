package main

import (
	"os"
	"os/exec"
)

func Pod_exec_by_kubectl() {
	cmd := exec.Command("kubectl", "exec", "-it", "klish", "--", "/bin/sh")

	// Set up input/output
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func main() {
	Pod_exec_by_kubectl()
}
