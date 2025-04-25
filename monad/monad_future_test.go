package monad

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Helper functions for testing
func assertEqual[T comparable](t *testing.T, expected, actual T) {
	t.Helper()
	if expected != actual {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

func assertError(t *testing.T, err error, expectedMsg string) {
	t.Helper()
	if err == nil {
		t.Errorf("Expected error with message '%s' but got nil", expectedMsg)
		return
	}
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s' but got '%s'", expectedMsg, err.Error())
	}
}
func TestHttpFutureWithTimeout(t *testing.T) {

	tests := []struct {
		name          string
		handler       http.HandlerFunc
		timeout       time.Duration
		expectError   bool
		expectStatus  int
		expectTimeout bool
	}{
		{
			name: "fast success",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("hello"))
			},
			timeout:       2 * time.Second,
			expectError:   false,
			expectStatus:  http.StatusOK,
			expectTimeout: false,
		},
		{
			name: "slow response triggers timeout",
			handler: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(3 * time.Second)
				w.WriteHeader(http.StatusOK)
			},
			timeout:       1 * time.Second,
			expectError:   false,
			expectStatus:  http.StatusOK,
			expectTimeout: true,
		},
		{
			name: "internal server error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			timeout:       2 * time.Second,
			expectError:   false,
			expectStatus:  http.StatusInternalServerError,
			expectTimeout: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			client := http.Client{
				Timeout: 5 * time.Second, // exceeds our cases(max 3 ?)
			}

			req, _ := http.NewRequest("GET", server.URL, nil)

			// create future request
			future := NewFuture(func() (*http.Response, error) {
				resp, err := client.Do(req)
				if err != nil {
					return nil, err
				}
				return resp, nil
			})

			// execute request with a set limit timeout
			resp, err, ok := future.GetWithTimeout(tc.timeout)
			if tc.expectError {
				if err == nil {
					t.Fatal("expected an error but got nil")
				}
				t.Logf("got expected error: %v", err)
			} else if tc.expectTimeout {
				if ok {
					t.Logf("got expected timeout: %v", err)
				} else {
					t.Fatalf("failed to timeout")
				}

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

func TestFutureBasicOperations(t *testing.T) {
	t.Run("successful future creation and retrieval", func(t *testing.T) {
		f := NewFuture(func() (int, error) {
			return 42, nil
		})

		result, err := f.Get()
		assertEqual(t, 42, result)
		assertEqual(t, nil, err)
	})

	t.Run("failed future creation and retrieval", func(t *testing.T) {
		expectedErr := errors.New("computation failed")
		f := NewFuture(func() (int, error) {
			return 0, expectedErr
		})

		result, err := f.Get()
		assertEqual(t, 0, result)
		assertEqual(t, expectedErr, err)
	})
}
func TestFutureTransformations(t *testing.T) {
	t.Run("successful Map transformation", func(t *testing.T) {
		f := NewFuture(func() (int, error) {
			return 21, nil
		})

		doubled := Map(f, func(x int) int {
			return x * 2
		})

		result, err := doubled.Get()
		assertEqual(t, 42, result)
		assertEqual(t, nil, err)
	})

	t.Run("Map with failed future", func(t *testing.T) {
		expectedErr := errors.New("computation failed")
		f := NewFuture(func() (int, error) {
			return 0, expectedErr
		})

		doubled := Map(f, func(x int) int {
			return x * 2
		})

		result, err := doubled.Get()
		assertEqual(t, 0, result)
		assertEqual(t, expectedErr, err)
	})

	t.Run("successful FlatMap transformation", func(t *testing.T) {
		f := NewFuture(func() (int, error) {
			return 21, nil
		})

		doubled := FlatMap(f, func(x int) *Future[int] {
			return Successful(x * 2)
		})

		result, err := doubled.Get()
		assertEqual(t, 42, result)
		assertEqual(t, nil, err)
	})

	t.Run("FlatMap with failed inner future", func(t *testing.T) {
		f := NewFuture(func() (int, error) {
			return 21, nil
		})

		expectedErr := errors.New("inner computation failed")
		doubled := FlatMap(f, func(x int) *Future[int] {
			return Failed[int](expectedErr)
		})

		result, err := doubled.Get()
		assertEqual(t, 0, result)
		assertEqual(t, expectedErr, err)
	})
}
func TestFutureSequence(t *testing.T) {
	t.Run("sequence of successful futures", func(t *testing.T) {
		futures := []*Future[int]{
			Successful(1),
			Successful(2),
			Successful(3),
		}

		combined := Sequence(futures...)
		result, err := combined.Get()

		assertEqual(t, nil, err)
		assertEqual(t, 3, len(result))
		assertEqual(t, 1, result[0])
		assertEqual(t, 2, result[1])
		assertEqual(t, 3, result[2])
	})

	t.Run("sequence with failed future", func(t *testing.T) {
		expectedErr := errors.New("future 2 failed")
		futures := []*Future[int]{
			Successful(1),
			Failed[int](expectedErr),
			Successful(3),
		}

		combined := Sequence(futures...)
		result, err := combined.Get()

		assertEqual(t, expectedErr, err)
		assertEqual(t, 0, len(result))
	})
}

func TestConcurrentOperations(t *testing.T) {
	t.Run("concurrent modifications", func(t *testing.T) {
		f := NewFuture(func() (int, error) {
			time.Sleep(50 * time.Millisecond)
			return 42, nil
		})

		const goroutines = 10
		done := make(chan bool, goroutines)

		for i := 0; i < goroutines; i++ {
			go func() {
				result, err := f.Get()
				if err != nil || result != 42 {
					t.Errorf("Concurrent access failed: got %v, %v", result, err)
				}
				done <- true
			}()
		}

		for i := 0; i < goroutines; i++ {
			<-done
		}
	})
}

func TestComplexChaining(t *testing.T) {
	t.Run("chain of transformations", func(t *testing.T) {
		// Start with a number
		f := Successful(5)

		// Convert to string
		strFuture := Map(f, func(n int) string {
			return fmt.Sprintf("Number: %d", n)
		})

		// Get length and double it
		lenFuture := FlatMap(strFuture, func(s string) *Future[int] {
			return Successful(len(s) * 2)
		})

		result, err := lenFuture.Get()
		assertEqual(t, nil, err)
		assertEqual(t, 18, result) // "Number: 5" has length 8, doubled to 16
	})

	t.Run("error propagation in chain", func(t *testing.T) {
		expectedErr := errors.New("middle chain error")

		f := Successful(5)

		// This transformation succeeds
		strFuture := Map(f, func(n int) string {
			return fmt.Sprintf("Number: %d", n)
		})

		// This transformation fails
		failedFuture := FlatMap(strFuture, func(s string) *Future[int] {
			return Failed[int](expectedErr)
		})

		// This transformation never executes due to previous error
		finalFuture := Map(failedFuture, func(n int) int {
			return n * 2
		})

		result, err := finalFuture.Get()
		assertEqual(t, 0, result)
		assertEqual(t, expectedErr, err)
	})
}

func TestEdgeCases(t *testing.T) {
	t.Run("zero value types", func(t *testing.T) {
		f := Successful[*int](nil)
		result, err := f.Get()
		assertEqual(t, (*int)(nil), result)
		assertEqual(t, nil, err)
	})

	t.Run("empty sequence", func(t *testing.T) {
		var futures []*Future[int]
		combined := Sequence(futures...)
		result, err := combined.Get()
		assertEqual(t, nil, err)
		assertEqual(t, 0, len(result))
	})

	/*

		func TestFutureTimeout(t *testing.T) {
			t.Run("successful completion within timeout", func(t *testing.T) {
				f := NewFuture(func() (int, error) {
					time.Sleep(50 * time.Millisecond)
					return 42, nil
				})

				result, err, ok := f.GetWithTimeout(100 * time.Millisecond)
				assertEqual(t, 42, result)
				assertEqual(t, nil, err)
			})

			t.Run("timeout exceeded", func(t *testing.T) {
				f := NewFuture(func() (int, error) {
					time.Sleep(200 * time.Millisecond)
					return 42, nil
				})

				_, err, _ := f.GetWithTimeout(100 * time.Millisecond)
				assertError(t, err, "timeout waiting for future")
			})
		}


					t.Run("immediate timeouts", func(t *testing.T) {
						f := NewFuture(func() (int, error) {
							time.Sleep(1 * time.Second)
							return 42, nil
						})

						_, err := f.GetWithTimeout(0)
						assertError(t, err, "timeout waiting for future")
				})
	*/

}
