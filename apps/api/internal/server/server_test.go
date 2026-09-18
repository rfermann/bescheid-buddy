package server_test

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/bored-bots/bored-template/internal/config"
	"github.com/bored-bots/bored-template/internal/server"
	"github.com/bored-bots/bored-template/internal/services"
	"github.com/stretchr/testify/require"
)

func TestServer_Shutdown(t *testing.T) {
	t.Run("graceful shutdown", func(t *testing.T) {
		signals := []syscall.Signal{
			syscall.SIGINT,
			syscall.SIGTERM,
		}

		for _, sig := range signals {
			sig := sig // capture range variable
			t.Run(fmt.Sprintf("shutdown with %s", sig), func(t *testing.T) {
				configPath := filepath.Join("..", "..", "testdata", "config.test.yaml")
				cfg, err := config.LoadConfig(configPath)
				require.NoError(t, err)

				app := server.NewServer(cfg, &services.Services{}, nil, nil)
				done := make(chan error, 1)

				go func() {
					done <- app.Start()
				}()

				// Give server time to start
				time.Sleep(50 * time.Millisecond)

				// Send signal
				proc, err := os.FindProcess(os.Getpid())
				if err != nil {
					t.Fatalf("Failed to find process: %v", err)
				}
				proc.Signal(sig)

				select {
				case err := <-done:
					if err != nil {
						t.Errorf("Unexpected error: %v", err)
					}
				case <-time.After(2 * time.Second):
					t.Error("Shutdown timed out")
				}
			})
		}
	})

	t.Run("shutdown timeout", func(t *testing.T) {
		configPath := filepath.Join("..", "..", "testdata", "config.test.yaml")
		cfg, err := config.LoadConfig(configPath)
		require.NoError(t, err)

		app := server.NewServer(cfg, &services.Services{}, nil, nil)

		// Add a long-running handler
		app.Router.Get("/slow", func(w http.ResponseWriter, r *http.Request) {
			t.Log("Slow handler started")
			time.Sleep(2 * time.Second)
			w.WriteHeader(http.StatusOK)
		})

		done := make(chan error, 1)
		go func() {
			done <- app.Start()
		}()

		// Give server time to start
		time.Sleep(50 * time.Millisecond)

		// Start a slow request
		go func() {
			res, err := http.Get(fmt.Sprintf("http://%s/slow", cfg.GetAddr()))
			if err != nil {
				t.Log("err", err)
			}
			t.Log("res.StatusCode", res.StatusCode)
		}()

		// Give time for request to start
		time.Sleep(50 * time.Millisecond)

		// Send interrupt signal
		proc, _ := os.FindProcess(os.Getpid())
		proc.Signal(syscall.SIGINT)

		select {
		case err := <-done:
			t.Log(err)
			if err == nil || !strings.Contains(err.Error(), "context deadline exceeded") {
				t.Errorf("Expected context deadline error, got: %v", err)
			}
		case <-time.After(5 * time.Second):
			t.Error("Test timed out")
		}
	})
}
