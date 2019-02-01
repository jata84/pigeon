package core

import (
	"testing"
)

func TestServer_Init(t *testing.T) {
	LoadConfig()

	t.Log(&buf)
}
