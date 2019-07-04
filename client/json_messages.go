package main

import "github.com/google/uuid"

type Connections struct {
	uuid uuid.UUID
	date string
}

type ConnectionList struct {
	list []Connections
}
