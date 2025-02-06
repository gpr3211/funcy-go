package fp

import (
	"reflect"
	"testing"
)

func TestFilter(t *testing.T) {
	tests := []struct {
		name      string
		input     []int
		predicate Predicate[int]
		expected  []int
	}{
		{
			name:      "Filter even numbers",
			input:     []int{1, 2, 3, 4, 5, 6},
			predicate: func(n int) bool { return n%2 == 0 },
			expected:  []int{2, 4, 6},
		},
		{
			name:      "Filter numbers greater than 3",
			input:     []int{1, 2, 3, 4, 5, 6},
			predicate: func(n int) bool { return n > 3 },
			expected:  []int{4, 5, 6},
		},
		{
			name:      "Filter with always true predicate",
			input:     []int{1, 2, 3},
			predicate: func(n int) bool { return true },
			expected:  []int{1, 2, 3},
		},
		{
			name:      "Filter with always false predicate",
			input:     []int{1, 2, 3},
			predicate: func(n int) bool { return false },
			expected:  []int{},
		},
		{
			name:      "Filter empty slice",
			input:     []int{},
			predicate: func(n int) bool { return n > 0 },
			expected:  []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Filter(tt.input, tt.predicate)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Filter() = %v\n want %v\n", result, tt.expected)
			}
		})
	}
}
