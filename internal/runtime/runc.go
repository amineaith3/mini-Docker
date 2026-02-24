package runtime

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)



func Run(args []string) {
	switch args[1] {
	case "run":
		parent(args)
	case "child":
		child(args)
	default:
		panic("Bad usage")
	}
}

func parent(args []string) {
	arguments := append([]string{"child"}, args[1:]...)


	cmd := exec.Command("/proc/self/exe", arguments...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	// cgroup logic here
	handle(os.MkdirAll("/sys/fs/cgroup/miniDocker", 0755))
	handle(os.WriteFile("/sys/fs/cgroup/miniDocker/memory.max", []byte("104857600"), 0700)) // 100 mo of memory

	handle(cmd.Start())
	pid := cmd.Process.Pid
	handle(os.WriteFile("/sys/fs/cgroup/miniDocker/cgroup.procs", []byte(strconv.Itoa(pid)), 0700))

	handle(cmd.Wait())
	os.RemoveAll("/sys/fs/cgroup/miniDocker") // cleanUP
}

func child(args []string) {

	handle(syscall.Mount("", "/", "", syscall.MS_REC|syscall.MS_PRIVATE, "")) // Mount and unmount will not propagate to parent
	handle(syscall.Chroot("/home/houcinee/Downloads/newRoot/"))
	handle(syscall.Chdir("/"))
	handle(syscall.Mount("proc", "/proc", "proc", 0, ""))

	os.Setenv("PATH", "/bin:/sbin:/usr/bin:/usr/sbin")
	binary, err := exec.LookPath(args[3]) // syscall.Exec() requires a full path to the binary
	handle(err)
	arguments := args[3:]

	handle(syscall.Exec(binary, arguments, os.Environ()))
}

func handle(err error) {
	if err != nil {
		fmt.Println("error : ", err)
		os.Exit(1)
	}
}
