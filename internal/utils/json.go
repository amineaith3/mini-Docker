// Package utils containes helper functions
package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"mini-docker/internal/models"

	"github.com/google/uuid"
)

func GenJSON(id string, st string, pid int, command []string) ([]byte, string) {
	if len(id) == 0 {
		id = uuid.New().String()
	}
	status := st
	createdAt := time.Now()

	if len(command) == 0 {
		command = []string{"/bin/sh"}
	}

	newContainer := models.Container{ID: id, Status: status, CreatedAt: createdAt, PID: pid, Command: command}

	jsonData, err := json.MarshalIndent(newContainer, "", "  ")
	Handle(err)
	return jsonData, id
}

func SaveFile(json []byte, id string) {
	directory := filepath.Join("/var/lib/minidocker/containers/", id)
	Handle(os.MkdirAll(directory, 0o755))
	os.WriteFile(filepath.Join(directory, "config.json"), json, 0o755)
}

func Handle(err error) {
	if err != nil {
		fmt.Println("error : ", err)
	}
}
