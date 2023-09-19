package main

import (
	"fmt"
	"os/exec"
	"syscall"
	"time"
)

func cliExec(timeout time.Duration, command string, args string) (string, error) {

	cmd := exec.Command(command)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
		CmdLine: " " + args,
		CreationFlags: 0,
	}

	fmt.Print(cmd.String() + " " + args)

	out, err := cmd.CombinedOutput()

	if (err != nil) {
		return string(out), err
	}

	return string(out), nil
}