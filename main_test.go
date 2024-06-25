package main

import (
	"os"
	"os/exec"
	"testing"
	"time"
)

// TestIntegrations runs integration tests using Newman (Postman CLI).
// It starts the server, waits for it to initialize, and then executes the tests.
func TestIntegrations(t *testing.T) {
	// Start the server in a separate goroutine to allow it to run concurrently with the tests.
	go func() {
		main()
	}()

	// Allow some time for the server to start
	time.Sleep(1 * time.Second)

	// Define the command and arguments for running the integration tests with Newman.
	cmd := "npx"
	args := []string{
		"newman",
		"run",
		"./tests/integration_tests.json",
		"--reporters",
		"cli,junit",
		"--reporter-junit-export",
		"integration_report.xml",
		"--insecure",
	}

	//Create a command
	command := exec.Command(cmd, args...)
	command.Stdout = os.Stdout // Set the standard output to the console.
	command.Stderr = os.Stderr // Set the standard error to the console.

	// Run the command and check for errors.
	err := command.Run()
	if err != nil {
		t.Errorf("Expected All Integration Tests to pass but got error")
	}
}
