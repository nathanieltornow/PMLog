package main

import "testing"

func TestHallo(t *testing.T) {
	m := make(map[uint64]chan bool)
	m[1] = make(chan bool, 1024)
	for i := 0; i < 50; i++ {
		go stuff(m)
	}
}
