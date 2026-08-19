package handler

import (
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCreateReportHandler_Delays(t *testing.T) {
	delays := []string{
		"0s",
		"3s",
		"7s",
		"12s",
	}

	dbDSN := os.Getenv("DB_DSN")
	if dbDSN == "" {
		t.Fatal("DB_DSN is not set")
	}

	projectRoot := findProjectRoot(t)

	// Build the API once.
	binaryPath := filepath.Join(projectRoot, "report-test-server")

	build := exec.Command(
		"go",
		"build",
		"-o",
		binaryPath,
		"./cmd/api",
	)

	build.Dir = projectRoot
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr

	if err := build.Run(); err != nil {
		t.Fatalf("failed to build server: %v", err)
	}

	// Remove the temporary executable when the test finishes.
	defer os.Remove(binaryPath)

	for _, delay := range delays {
		t.Run(delay, func(t *testing.T) {

			// Start the actual API executable.
			cmd := exec.Command(
				binaryPath,
				"-db-dsn="+dbDSN,
				"-report-delay="+delay,
			)

			cmd.Dir = projectRoot
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Start(); err != nil {
				t.Fatalf("failed to start server: %v", err)
			}

			// Make absolutely sure this server is stopped
			// before moving to the next delay.
			defer func() {
				if cmd.Process != nil {
					_ = cmd.Process.Kill()
				}

				_ = cmd.Wait()
			}()

			// Wait until this server is accepting connections.
			waitForServer(t)

			// Start measuring the actual API request.
			start := time.Now()

			curl := exec.Command(
				"curl",
				"--silent",
				"--show-error",
				"--output",
				"/dev/null",
				"--write-out",
				"%{http_code} %{time_total} %{size_download}",
				"--request",
				"POST",
				"--header",
				"Content-Type: application/json",
				"--data",
				`{
					"consumer_id": "0198f000-0000-7000-8000-000000000001",
					"from": "2026-01-01T00:00:00Z",
					"to": "2027-01-01T00:00:00Z"
				}`,
				"http://localhost:4000/api/reports",
			)

			output, err := curl.Output()

			elapsed := time.Since(start)

			if err != nil {
				t.Fatalf("curl failed: %v", err)
			}

			fields := strings.Fields(string(output))

			if len(fields) != 3 {
				t.Fatalf("unexpected curl output: %q", output)
			}

			status := fields[0]
			responseTime := fields[1]
			bytes := fields[2]

			t.Logf(
				"delay=%s status=%s time=%ss bytes=%s actual_elapsed=%s",
				delay,
				status,
				responseTime,
				bytes,
				elapsed.Round(time.Millisecond),
			)
		})
	}
}

func waitForServer(t *testing.T) {
	t.Helper()

	client := http.Client{
		Timeout: 200 * time.Millisecond,
	}

	for i := 0; i < 50; i++ {
		resp, err := client.Get("http://localhost:4000")

		if err == nil {
			resp.Body.Close()
			return
		}

		time.Sleep(100 * time.Millisecond)
	}

	t.Fatal("server did not become available")
}

func findProjectRoot(t *testing.T) string {
	t.Helper()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)

		if parent == dir {
			t.Fatal("could not find project root containing go.mod")
		}

		dir = parent
	}
}
