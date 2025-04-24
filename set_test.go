package fp

import (
	"reflect"
	"testing"
)

// Product struct for testing
type Prod struct {
	Name     string
	Category string
	Price    float64
}

type Point struct {
	X, Y int
}

type CustomID int

func TestSetWithStructsAndCustomTypes(t *testing.T) {
	t.Run("Struct slice with duplicates", func(t *testing.T) {
		tests := []struct {
			name     string
			input    []Prod
			expected []Prod
		}{
			{
				name: "Duplicate Products",
				input: []Prod{
					{"Kayak", "Watersports", 279},
					{"Kayak", "Watersports", 279}, // Duplicate
					{"Soccer Ball", "Soccer", 19.50},
					{"Chess Board", "Chess", 50},
					{"Soccer Ball", "Soccer", 19.50}, // Duplicate
				},
				expected: []Prod{
					{"Kayak", "Watersports", 279},
					{"Soccer Ball", "Soccer", 19.50},
					{"Chess Board", "Chess", 50},
				},
			},
			{
				name:     "Empty Product slice",
				input:    []Prod{},
				expected: nil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := Set(tt.input)
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("Set(%v) = %v, want %v", tt.input, result, tt.expected)
				}
			})
		}
	})

	t.Run("Structs with custom fields", func(t *testing.T) {
		tests := []struct {
			name     string
			input    []Point
			expected []Point
		}{
			{
				name: "Points with duplicates",
				input: []Point{
					{X: 1, Y: 2},
					{X: 3, Y: 4},
					{X: 1, Y: 2}, // Duplicate
					{X: 5, Y: 6},
				},
				expected: []Point{
					{X: 1, Y: 2},
					{X: 3, Y: 4},
					{X: 5, Y: 6},
				},
			},
			{
				name:     "Empty point slice",
				input:    []Point{},
				expected: nil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := Set(tt.input)
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("Set(%v) = %v, want %v", tt.input, result, tt.expected)
				}
			})
		}
	})
}

func TestUnionWithStructsAndCustomTypes(t *testing.T) {
	tests := []struct {
		name     string
		set1     []Prod
		set2     []Prod
		expected []Prod
	}{
		{
			name: "Union of Product slices",
			set1: []Prod{
				{"Kayak", "Watersports", 279},
				{"Soccer Ball", "Soccer", 19.50},
			},
			set2: []Prod{
				{"Chess Board", "Chess", 50},
				{"Soccer Ball", "Soccer", 19.50}, // Duplicate
			},
			expected: []Prod{
				{"Kayak", "Watersports", 279},
				{"Soccer Ball", "Soccer", 19.50},
				{"Chess Board", "Chess", 50},
			},
		},
		{
			name:     "Union of empty slices",
			set1:     []Prod{},
			set2:     []Prod{},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Union(tt.set1, tt.set2)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Union(%v, %v) = %v, want %v", tt.set1, tt.set2, result, tt.expected)
			}
		})
	}
}

func TestIntersectionWithStructs(t *testing.T) {
	tests := []struct {
		name     string
		set1     []Prod
		set2     []Prod
		expected []Prod
	}{
		{
			name: "Intersection of Product slices",
			set1: []Prod{
				{"Kayak", "Watersports", 279},
				{"Soccer Ball", "Soccer", 19.50},
				{"Chess Board", "Chess", 50},
			},
			set2: []Prod{
				{"Soccer Ball", "Soccer", 19.50},
				{"Chess Board", "Chess", 50},
			},
			expected: []Prod{
				{"Soccer Ball", "Soccer", 19.50},
				{"Chess Board", "Chess", 50},
			},
		},
		{
			name:     "Empty intersection",
			set1:     []Prod{},
			set2:     []Prod{},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Intersection(tt.set1, tt.set2)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Intersection(%v, %v) = %v, want %v", tt.set1, tt.set2, result, tt.expected)
			}
		})
	}
}

func TestDifferenceWithStructs(t *testing.T) {
	tests := []struct {
		name     string
		set1     []Prod
		set2     []Prod
		expected []Prod
	}{
		{
			name: "Difference of Product slices",
			set1: []Prod{
				{"Kayak", "Watersports", 279},
				{"Soccer Ball", "Soccer", 19.50},
			},
			set2: []Prod{
				{"Soccer Ball", "Soccer", 19.50},
			},
			expected: []Prod{
				{"Kayak", "Watersports", 279},
			},
		},
		{
			name:     "Empty difference",
			set1:     []Prod{},
			set2:     []Prod{},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Difference(tt.set1, tt.set2)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Difference(%v, %v) = %v, want %v", tt.set1, tt.set2, result, tt.expected)
			}
		})
	}
}
