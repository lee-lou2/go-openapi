package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	_ = os.Setenv("TESTING", "true")
	code := m.Run()
	_ = os.Unsetenv("TESTING")
	os.Exit(code)
}
