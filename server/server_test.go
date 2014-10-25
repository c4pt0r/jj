package server

import "testing"

func TestServer(t *testing.T) {
	s := NewServer(":9999")
	s.Run()
}
