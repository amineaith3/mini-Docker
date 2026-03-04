package runtime

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
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

// ./miniDocker run <options> <command> <args>
func parent(args []string) {
	containerDirectory := "/home/houcinee/Downloads/newRoot/"
	command, options := parseRunSpec(args[2:])
	volumeMounts, err := parseVolumeSpecs(options["-v"], containerDirectory)
	utils.Handle(err)
	utils.Handle(setupVolumeMounts(volumeMounts))
	defer func() {
		utils.Handle(cleanupVolumeMounts(volumeMounts))
	}()

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
	json, id := utils.GenJSON("", "CREATED", 0, command, createdAt)
	utils.SaveFile(json, id)

	utils.Handle(cmd.Start())
	pid := cmd.Process.Pid
	json, id = utils.GenJSON(id, "RUNNING", pid, command, createdAt)
	utils.SaveFile(json, id)

	utils.Handle(os.WriteFile("/sys/fs/cgroup/miniDocker/cgroup.procs", []byte(strconv.Itoa(pid)), 0o700))

	utils.Handle(cmd.Wait())
	json, id = utils.GenJSON(id, "EXITED", pid, command, createdAt)
	utils.SaveFile(json, id)
	os.RemoveAll("/sys/fs/cgroup/miniDocker") // cleanUP
}

// ./miniDocker child run <options> <command> <args>
func child(args []string) {
	containerDirectory := "/home/houcinee/Downloads/newRoot/"
	command, options := parseRunSpec(args[3:])
	_, err := parseVolumeSpecs(options["-v"], containerDirectory)
	utils.Handle(err)

	for opt := range options {
		switch opt {
		case "-v":
			continue
		default:
			fmt.Println("I can't handle that option yet")
		}
	}

	utils.Handle(syscall.Mount("", "/", "", syscall.MS_REC|syscall.MS_PRIVATE, "")) // Mount and unmount will not propagate to parent
	utils.Handle(syscall.Chroot(containerDirectory))
	utils.Handle(syscall.Chdir("/"))
	utils.Handle(syscall.Mount("proc", "/proc", "proc", 0, ""))

	os.Setenv("PATH", "/bin:/sbin:/usr/bin:/usr/sbin")
	binary, err := exec.LookPath(command[0]) // syscall.Exec() requires a full path to the binary
	utils.Handle(err)

	utils.Handle(syscall.Exec(binary, command, os.Environ()))
}

func parseRunSpec(args []string) ([]string, map[string][]string) {
	var command []string
	options := map[string][]string{}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			if i+1 < len(args) {
				options[arg] = append(options[arg], args[i+1])
				i++
			} else {
				panic("Bad usage")
			}
		} else {
			command = args[i:]
			break
		}
	}
	return command, options
}

func parseVolumeSpecs(values []string, containerRoot string) ([]VolumeMount, error) {
	var mounts []VolumeMount

	for _, val := range values {
		directories := strings.SplitN(val, ":", 2)
		if len(directories) != 2 {
			return nil, errors.New("invalid volume format. Expected /host/path:/container/path")
		}

		src := directories[0]
		target := directories[1]

		if !filepath.IsAbs(src) || !filepath.IsAbs(target) {
			return nil, errors.New("both host and container paths must be absolute (start with '/')")
		}

		mounts = append(mounts, VolumeMount{
			Source:     src,
			HostTarget: filepath.Join(containerRoot, target),
		})
	}

	return mounts, nil
}

func setupVolumeMounts(mounts []VolumeMount) error {
	var mounted []VolumeMount

	for _, mount := range mounts {
		if err := os.MkdirAll(mount.Source, 0o755); err != nil {
			_ = cleanupVolumeMounts(mounted)
			return err
		}

		if err := os.MkdirAll(mount.HostTarget, 0o755); err != nil {
			_ = cleanupVolumeMounts(mounted)
			return err
		}

		if err := syscall.Mount(mount.Source, mount.HostTarget, "bind", syscall.MS_BIND, ""); err != nil {
			_ = cleanupVolumeMounts(mounted)
			return err
		}

		mounted = append(mounted, mount)
	}

	return nil
}

func cleanupVolumeMounts(mounts []VolumeMount) error {
	var errs []error
	for i := len(mounts) - 1; i >= 0; i-- {
		mount := mounts[i]
		if err := syscall.Unmount(mount.HostTarget, 0); err != nil {
			errs = append(errs, fmt.Errorf("unmount %s: %w", mount.HostTarget, err))
		}
	}

	return errors.Join(errs...)
}
