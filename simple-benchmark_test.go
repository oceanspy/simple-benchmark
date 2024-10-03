package main

import (
	"testing"
)

func TestContainsForbiddenCmd(t *testing.T) {
	if ContainsForbiddenCmd("echo 'hello world'") {
		t.Errorf("echo should not be a forbidden command")
	}

	if !ContainsForbiddenCmd("rm -rf /truc") {
		t.Errorf("rm should be a forbidden command")
	}

	if !ContainsForbiddenCmd("sudo rm -rf /truc") {
		t.Errorf("rm should be a forbidden command")
	}
}
