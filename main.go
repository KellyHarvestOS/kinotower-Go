package main

import (
	"log"
	"os"
	"os/exec"
)

// Compatibility wrapper: the production entrypoint is ./cmd/api.
func main() {
	cmd := exec.Command("go", "run", "./cmd/api")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
