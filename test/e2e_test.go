package test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var (
	binaryPath string
)

// TestMain builds the binary before running tests
func TestMain(m *testing.M) {
	// Build the binary
	fmt.Println("Building Alexandria binary...")
	binaryPath = filepath.Join("..", "alexandria-test")
	cmd := exec.Command("go", "build", "-o", binaryPath, "..")
	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("Failed to build binary: %v\n%s\n", err, output)
		os.Exit(1)
	}

	// Run tests using standard database connection
	code := m.Run()

	// Cleanup binary
	os.Remove(binaryPath)

	os.Exit(code)
}

// runCommand executes the Alexandria binary with given arguments
func runCommand(t *testing.T, args ...string) (string, string, error) {
	cmd := exec.Command(binaryPath, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func TestBinaryBuilds(t *testing.T) {
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Fatal("Binary was not built")
	}
}

func TestHelpCommand(t *testing.T) {
	stdout, _, err := runCommand(t, "--help")
	if err != nil {
		t.Fatalf("Help command failed: %v", err)
	}

	if !strings.Contains(stdout, "alexandria") {
		t.Error("Help output doesn't contain 'alexandria'")
	}
}

func TestSourceStatusCommand(t *testing.T) {
	stdout, _, err := runCommand(t, "source", "--status")
	if err != nil {
		t.Fatalf("Source status command failed: %v", err)
	}

	if !strings.Contains(stdout, "Database Type:") {
		t.Error("Source status output doesn't contain 'Database Type:'")
	}

	// Should default to sqlite
	if !strings.Contains(stdout, "sqlite") {
		t.Error("Default database type should be sqlite")
	}
}

func TestSourceSwitchToSQLite(t *testing.T) {
	stdout, stderr, err := runCommand(t, "source", "sqlite")
	if err != nil {
		t.Fatalf("Source switch to sqlite failed: %v\nStderr: %s", err, stderr)
	}

	if !strings.Contains(stdout, "Successfully switched to sqlite") {
		t.Errorf("Expected success message, got: %s", stdout)
	}
}

func TestSourceSwitchInvalidType(t *testing.T) {
	_, stderr, err := runCommand(t, "source", "invalid")
	if err == nil {
		t.Error("Expected error when switching to invalid database type")
	}

	combinedOutput := stderr
	if !strings.Contains(combinedOutput, "invalid database type") {
		t.Errorf("Expected error message about invalid database type, got: %s", combinedOutput)
	}
}

func TestCreateTicket(t *testing.T) {
	stdout, stderr, err := runCommand(t, "create",
		"--title", "Test Ticket",
		"--description", "This is a test ticket",
		"--type", "task",
		"--priority", "medium",
		"--project", "TestProject",
	)

	if err != nil {
		t.Fatalf("Create ticket failed: %v\nStdout: %s\nStderr: %s", err, stdout, stderr)
	}

	if !strings.Contains(stdout, "created successfully") && !strings.Contains(stdout, "Ticket") {
		t.Errorf("Expected success message for ticket creation, got: %s", stdout)
	}
}

func TestListTickets(t *testing.T) {
	// First create a ticket
	runCommand(t, "create",
		"--title", "List Test Ticket",
		"--type", "bug",
		"--priority", "high",
		"--project", "TestProject",
	)

	// Then list tickets
	stdout, stderr, err := runCommand(t, "list")
	if err != nil {
		t.Fatalf("List tickets failed: %v\nStderr: %s", err, stderr)
	}

	// Should show tickets in table format
	if !strings.Contains(stdout, "ID") || !strings.Contains(stdout, "TITLE") {
		t.Errorf("List output should contain table headers, got: %s", stdout)
	}
}

func TestListTicketsJSON(t *testing.T) {
	stdout, stderr, err := runCommand(t, "list", "-o", "json")
	if err != nil {
		t.Fatalf("List tickets with JSON output failed: %v\nStderr: %s", err, stderr)
	}

	// Should output JSON
	if !strings.Contains(stdout, "[") && !strings.Contains(stdout, "No tickets found") {
		t.Errorf("Expected JSON array or 'No tickets found', got: %s", stdout)
	}
}

func TestListTicketsWithFilter(t *testing.T) {
	_, stderr, err := runCommand(t, "list", "--status", "open", "--type", "bug")
	if err != nil {
		t.Fatalf("List tickets with filter failed: %v\nStderr: %s", err, stderr)
	}

	// Should complete without error
	if stderr != "" && !strings.Contains(stderr, "WARNING") {
		t.Errorf("Unexpected stderr: %s", stderr)
	}
}

func TestViewTicket(t *testing.T) {
	// Create a ticket first
	stdout, stderr, err := runCommand(t, "create",
		"--title", "View Test Ticket",
		"--description", "Detailed description for viewing",
		"--type", "feature",
		"--priority", "low",
		"--project", "TestProject",
	)
	if err != nil {
		t.Fatalf("Failed to create ticket: %v\nStdout: %s\nStderr: %s", err, stdout, stderr)
	}

	// View ticket
	stdout, stderr, err = runCommand(t, "view", "--title", "View Test Ticket", "--project", "TestProject")
	if err != nil {
		t.Fatalf("View ticket failed: %v\nStderr: %s", err, stderr)
	}

	if !strings.Contains(stdout, "View Test Ticket") && !strings.Contains(stdout, "Ticket") {
		t.Errorf("Expected ticket details in view output, got: %s", stdout)
	}
}

func TestUpdateTicket(t *testing.T) {
	// Create a ticket first
	stdout, stderr, err := runCommand(t, "create",
		"--title", "Update Test Ticket",
		"--type", "task",
		"--priority", "medium",
		"--project", "TestProject",
	)
	if err != nil {
		t.Fatalf("Failed to create ticket: %v\nStdout: %s\nStderr: %s", err, stdout, stderr)
	}

	// Update the ticket
	stdout, stderr, err = runCommand(t, "update",
		"--project", "TestProject",
		"--title", "Update Test Ticket",
		"--status", "in-progress",
		"--priority", "high",
	)

	if err != nil {
		t.Fatalf("Update ticket failed: %v\nStdout: %s\nStderr: %s", err, stdout, stderr)
	}

	if !strings.Contains(stdout, "updated successfully") && !strings.Contains(stdout, "Ticket") {
		t.Errorf("Expected success message for ticket update, got: %s", stdout)
	}
}

func TestDeleteTicket(t *testing.T) {
	// Create a ticket to delete
	stdout, stderr, err := runCommand(t, "create",
		"--title", "Delete Test Ticket",
		"--type", "task",
		"--priority", "low",
		"--project", "TestProject",
	)
	if err != nil {
		t.Fatalf("Failed to create ticket: %v\nStdout: %s\nStderr: %s", err, stdout, stderr)
	}

	// Delete the ticket
	stdout, stderr, err = runCommand(t, "delete",
		"--project", "TestProject",
		"--title", "Delete Test Ticket",
	)
	if err != nil {
		t.Fatalf("Delete ticket failed: %v\nStdout: %s\nStderr: %s", err, stdout, stderr)
	}

	if !strings.Contains(stdout, "deleted successfully") && !strings.Contains(stdout, "Ticket") {
		t.Errorf("Expected success message for ticket deletion, got: %s", stdout)
	}
}

func TestDatabaseSwitchingWorkflow(t *testing.T) {
	// Start with SQLite
	t.Log("Switching to SQLite...")
	stdout, stderr, err := runCommand(t, "source", "sqlite")
	if err != nil {
		t.Fatalf("Failed to switch to SQLite: %v\nStderr: %s", err, stderr)
	}
	if !strings.Contains(stdout, "sqlite") {
		t.Errorf("Expected confirmation of SQLite switch, got: %s", stdout)
	}

	// Verify we're on SQLite
	t.Log("Verifying SQLite is active...")
	stdout, _, err = runCommand(t, "source", "--status")
	if err != nil {
		t.Fatalf("Failed to get source status: %v", err)
	}
	if !strings.Contains(stdout, "sqlite") {
		t.Errorf("Expected database type to be sqlite, got: %s", stdout)
	}

	// Create a ticket on SQLite
	t.Log("Creating ticket on SQLite...")
	runCommand(t, "create",
		"--title", "SQLite Ticket",
		"--type", "task",
		"--priority", "medium",
		"--project", "TestProject",
	)

	// List tickets
	stdout, _, _ = runCommand(t, "list")
	if !strings.Contains(stdout, "SQLite Ticket") && !strings.Contains(stdout, "Total:") {
		t.Log("Note: Ticket may not be visible in list output")
	}
}
