// Package manager is used to manage containers
package manager

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"

	"mini-docker/internal/models"
	"mini-docker/internal/runtime"
	"mini-docker/internal/utils"
)

func Run(args []string) {
	path := "/var/lib/minidocker/containers/"
	switch args[1] {
	case "run", "child":
		runtime.Run(args)

	case "ps":
		list := getContainers(path)
		if len(args) > 2 && args[2] == "-a" {
			printContainers(list)
		} else if len(args) > 2 {
			panic("Bad usage, this option doesn't exist")
		} else {
			printContainers(checkActive(list))
		}

	case "inspect":
		if len(args) < 3 {
			panic("Bad usage, containerID required")
		}
		dir := filepath.Join(path, args[2])
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			panic("The provided ID doesn't exist")
		}
		c, err := readContainerConfig(dir)
		utils.Handle(err)
		live := "EXITED"
		if isProcessActive(c.PID) {
			live = "RUNNING"
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "CONTAINER ID\tCOMMAND\tCREATED\tSTATUS\tPID")
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\n",
			c.ID, strings.Join(c.Command, " "), formatAge(c.CreatedAt), live, c.PID)
		w.Flush()

	default:
		fmt.Println("We will do something soon")
	}
}

func getContainers(path string) []models.Container {
	var list []models.Container
	entries, err := os.ReadDir(path)
	utils.Handle(err)

	for _, v := range entries {
		c, err := readContainerConfig(filepath.Join(path, v.Name()))
		utils.Handle(err)
		list = append(list, c)
	}
	return list
}

// readContainerConfig reads and unmarshals the config.json from a container directory.
func readContainerConfig(dir string) (models.Container, error) {
	content, err := os.ReadFile(filepath.Join(dir, "config.json"))
	if err != nil {
		return models.Container{}, err
	}
	var c models.Container
	return c, json.Unmarshal(content, &c)
}

// isProcessActive reports whether the process with the given PID is alive
// by checking for the existence of its /proc/<pid> directory.
func isProcessActive(pid int) bool {
	if pid <= 0 {
		return false
	}
	_, err := os.Stat(fmt.Sprintf("/proc/%d", pid))
	return err == nil
}

// checkActive returns only containers whose process is currently alive.
func checkActive(list []models.Container) []models.Container {
	var active []models.Container
	for _, c := range list {
		if isProcessActive(c.PID) {
			active = append(active, c)
		}
	}
	return active
}

func formatAge(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return fmt.Sprintf("%d seconds ago", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%d minutes ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%d hours ago", int(d.Hours()))
	default:
		return fmt.Sprintf("%d days ago", int(d.Hours()/24))
	}
}

func printContainers(list []models.Container) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "CONTAINER ID\tCOMMAND\tCREATED\tSTATUS\tPID")
	for _, c := range list {
		live := "EXITED"
		if isProcessActive(c.PID) {
			live = "RUNNING"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\n",
			c.ID, strings.Join(c.Command, " "), formatAge(c.CreatedAt), live, c.PID)
	}
	w.Flush()
}
