// Package utils containes helper functions
package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"mini-docker/internal/models"
)

func GenJSON(id string, st string, pid int, command []string, createdAt time.Time) ([]byte, string) {
	if createdAt.IsZero() {
		createdAt = time.Now()
	}

	if len(command) == 0 {
		command = []string{"/bin/sh"}
	}

	newContainer := models.Container{ID: id, Status: st, CreatedAt: createdAt, PID: pid, Command: command}

	jsonData, err := json.MarshalIndent(newContainer, "", "  ")
	Handle(err)
	return jsonData, id
}

func SaveFile(json []byte, id string) {
	directory := filepath.Join("/var/lib/minidocker/containers/", id)
	Handle(os.MkdirAll(directory, 0o755))
	Handle(os.WriteFile(filepath.Join(directory, "config.json"), json, 0o755))
}

func Handle(err error) {
	if err != nil {
		fmt.Println("error : ", err)
		os.Exit(1)
	}
}
