package main

import (
	"fmt"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	t.Run("server should shut down after a SIGINT", func(t *testing.T) {
		httpPort := "5000"
		t.Setenv("HTTP_PORT", httpPort)

		go func() {
			main()
		}()

		// wait for the server to start
		time.Sleep(2 * time.Second)

		p, err := os.FindProcess(os.Getpid())
		if err != nil {
			t.Fatalf("Unable to find process: %v", err)
		}

		if err := p.Signal(syscall.SIGINT); err != nil {
			t.Fatalf("Unable to send SIGINT: %v", err)
		}

		time.Sleep(10 * time.Second)

		// Check if server is no longer running
		resp, err := http.Get(fmt.Sprintf("http://localhost:%s", httpPort))
		if err == nil {
			resp.Body.Close()
			t.Fatalf("Expected server to be shut down, but it is still running")
		}
	})

}
