package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		parent()
	case "child":
		child()
	default:
		panic("Bad usage")
	}
}

func parent() {
	arguments := append([]string{"child"}, os.Args[2:]...)

	cmd := exec.Command("/proc/self/exe", arguments...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	err := cmd.Run()
	handle(err)
}

func child() {
	handle(syscall.Mount("", "/", "", syscall.MS_REC|syscall.MS_PRIVATE, "")) // Mount and unmount will not propagate to parent
	handle(syscall.Mount("proc", "/proc", "proc", 0, ""))

	binary, err := exec.LookPath(os.Args[2]) // syscall.Exec() requires a full path to the binary
	handle(err)

	args := os.Args[3:]

	handle(syscall.Exec(binary, args, os.Environ()))
}

func handle(err error) {
	if err != nil {
		fmt.Println("error : ", err)
	}
}
