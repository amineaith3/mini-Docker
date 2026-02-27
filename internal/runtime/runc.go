package runtime

import (
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"

	"mini-docker/internal/utils"
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
	utils.Handle(os.MkdirAll("/sys/fs/cgroup/miniDocker", 0o755))
	utils.Handle(os.WriteFile("/sys/fs/cgroup/miniDocker/memory.max", []byte("104857600"), 0o700)) // 100 mo of memory

	// generate meta data
	createdAt := time.Now()
	json, id := utils.GenJSON("", "CREATED", 0, args[2:], createdAt)
	utils.SaveFile(json, id)

	utils.Handle(cmd.Start())
	pid := cmd.Process.Pid
	json, id = utils.GenJSON(id, "RUNNING", pid, args[2:], createdAt)
	utils.SaveFile(json, id)

	utils.Handle(os.WriteFile("/sys/fs/cgroup/miniDocker/cgroup.procs", []byte(strconv.Itoa(pid)), 0o700))

	utils.Handle(cmd.Wait())
	json, id = utils.GenJSON(id, "EXITED", pid, args[2:], createdAt)
	utils.SaveFile(json, id)
	os.RemoveAll("/sys/fs/cgroup/miniDocker") // cleanUP
}

func child(args []string) {
	utils.Handle(syscall.Mount("", "/", "", syscall.MS_REC|syscall.MS_PRIVATE, "")) // Mount and unmount will not propagate to parent
	utils.Handle(syscall.Chroot("/home/houcinee/Downloads/newRoot/"))
	utils.Handle(syscall.Chdir("/"))
	utils.Handle(syscall.Mount("proc", "/proc", "proc", 0, ""))

	os.Setenv("PATH", "/bin:/sbin:/usr/bin:/usr/sbin")
	binary, err := exec.LookPath(args[3]) // syscall.Exec() requires a full path to the binary
	utils.Handle(err)
	arguments := args[3:]

	utils.Handle(syscall.Exec(binary, arguments, os.Environ()))
}
