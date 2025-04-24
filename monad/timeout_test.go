package monad

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestHttpFutureWithTimeout(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	tests := []struct {
		name         string
		handler      http.HandlerFunc
		timeout      time.Duration
		expectError  bool
		expectStatus int
	}{
		{
			name: "fast success",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("hello"))
			},
			timeout:      2 * time.Second,
			expectError:  false,
			expectStatus: http.StatusOK,
		},
		{
			name: "slow response triggers timeout",
			handler: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(3 * time.Second)
				w.WriteHeader(http.StatusOK)
			},
			timeout:      1 * time.Second,
			expectError:  true,
			expectStatus: 0,
		},
		{
			name: "internal server error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			timeout:      2 * time.Second,
			expectError:  false,
			expectStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := http.Client{
				Timeout: 5 * time.Second, // Higher than our Future timeout
			}

			req, _ := http.NewRequest("GET", server.URL, nil)

			future := NewFuture(func() (*http.Response, error) {
				resp, err := client.Do(req)
				if err != nil {
					return nil, err
				}
				return resp, nil
			})

			resp, err := future.GetWithTimeout(tc.timeout)
			if tc.expectError {
				if err == nil {
					t.Fatal("expected an error but got nil")
				}
				t.Logf("got expected error: %v", err)
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if resp.StatusCode != tc.expectStatus {
					t.Fatalf("expected status %d but got %d", tc.expectStatus, resp.StatusCode)
				}
			}
		})
	}
}
