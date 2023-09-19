package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

func cliExec(timeout time.Duration, name string, arg ...string) (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel() 
	
	cmd := exec.Command(name, arg...)

	out, err := cmd.CombinedOutput()

	if (err != nil) {
		return string(out), err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("Timeout exceeded")
	}

	return string(out), nil
}