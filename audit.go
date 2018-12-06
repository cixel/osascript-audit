package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	f, err := os.OpenFile("/usr/local/var/log/osascript-audit.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer func() {
		_ = f.Close()
	}()

	log.SetOutput(f)

	osbin, err := exec.LookPath("/usr/bin/osascript")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log.Println("osascript invoked")
	log.Printf("ppid: %d", os.Getppid())
	log.Printf("gid:  %d", os.Getgid())
	log.Println("argv: ", os.Args[1:])

	osascript := exec.Command(osbin, os.Args[1:]...)
	osascript.Stdout = os.Stdout
	osascript.Stderr = os.Stderr
	osascript.Stdin = os.Stdin

	err = osascript.Run()

	var exitCode int
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			ws := exitErr.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		} else {
			exitCode = 1
		}
	}

	os.Exit(exitCode)
}
