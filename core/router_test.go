package core

import (
	"log"
	"testing"
)

func TestRouter(t *testing.T) {
	LoadConfig()
	log.Println("Testing Router")
	status := make(chan string)
	r := NewRouter(status)
	r.Init()

	t.Log(&buf)
}
