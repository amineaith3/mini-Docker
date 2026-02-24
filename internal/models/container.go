// Package models defines types related to containerd
package models

import "time"

type Container struct {
	ID        string    `json:"id"`
	PID       int       `json:"pid"`
	Status    string    `json:"status"`
	Command   []string  `json:"command"`
	CreatedAt time.Time `json:"created_at"`
}
