package main

import "testing"

func TestRun(t *testing.T) {
	// when
	err := run()
	if err != nil {
		t.Error("failed run()")
	}

	// then
}