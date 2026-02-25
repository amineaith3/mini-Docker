// Package models defines types related to containerd
package models

import "time"

type Container struct {
	ID        string    `json:"id"`
	PID       int       `json:"pid,omitempty"`
	Status    string    `json:"status"`
	Command   []string  `json:"command,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
