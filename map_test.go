package fp

import (
	"reflect"
	"testing"
)

func TestMap(t *testing.T) {
	t.Run("Integer operations", func(t *testing.T) {
		tests := []struct {
			name     string
			input    []int
			mapFunc  MapFunc[int]
			expected []int
		}{
			{
				name:     "Double integers",
				input:    []int{1, 2, 3, 4, 5},
				mapFunc:  func(n int) int { return n * 2 },
				expected: []int{2, 4, 6, 8, 10},
			},
			{
				name:     "Add one to integers",
				input:    []int{0, 1, 2, 3, 4},
				mapFunc:  func(n int) int { return n + 1 },
				expected: []int{1, 2, 3, 4, 5},
			},
			{
				name:     "Map empty integer slice",
				input:    []int{},
				mapFunc:  func(n int) int { return n * 2 },
				expected: []int{},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, _ := Map(tt.input, tt.mapFunc)
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("Map() = %v, want %v", result, tt.expected)
				}
			})
		}
	})

	t.Run("String operations", func(t *testing.T) {
		tests := []struct {
			name     string
			input    []string
			mapFunc  MapFunc[string]
			expected []string
		}{
			{
				name:     "Uppercase strings",
				input:    []string{"hello", "world", "go"},
				mapFunc:  func(s string) string { return s + "!" },
				expected: []string{"hello!", "world!", "go!"},
			},
			{
				name:     "Prefix strings",
				input:    []string{"apple", "banana", "cherry"},
				mapFunc:  func(s string) string { return "fruit_" + s },
				expected: []string{"fruit_apple", "fruit_banana", "fruit_cherry"},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, _ := Map(tt.input, tt.mapFunc)
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("Map() = %v, want %v", result, tt.expected)
				}
			})
		}
	})

	t.Run("Custom type operations", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		tests := []struct {
			name     string
			input    []Person
			mapFunc  MapFunc[Person]
			expected []Person
		}{
			{
				name: "Increase age",
				input: []Person{
					{Name: "Alice", Age: 30},
					{Name: "Bob", Age: 25},
					{Name: "Charlie", Age: 35},
				},
				mapFunc: func(p Person) Person {
					p.Age++
					return p
				},
				expected: []Person{
					{Name: "Alice", Age: 31},
					{Name: "Bob", Age: 26},
					{Name: "Charlie", Age: 36},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := Map(tt.input, tt.mapFunc)
				if err != nil {
					t.Errorf("Map() = %v, nil value found by assert", result)
				}
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("Map() = %v, want %v", result, tt.expected)
				}
			})
		}
	})
}
