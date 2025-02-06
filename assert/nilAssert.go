package assert

import (
	"fmt"
	"reflect"
)

// AssertNotNil checks if any pointer fields in a struct are nil.
// It returns an error listing all nil fields, or nil if all fields are non-nil.
// Generic type T must be a struct.
func AssertNotNil[T any](s T) error {
	// Get the reflect value and type of the struct.
	v := reflect.ValueOf(s)
	t := v.Type()

	// Ensure we're working with a struct.
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("AssertNotNil: type %s is not a struct", t.Name())
	}
	var nilFields []string

	// Iterate through all fields in the struct.
	for i := 0; i < t.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Check if the field is a pointer.
		if field.Kind() == reflect.Pointer {
			if field.IsNil() {
				nilFields = append(nilFields, fieldType.Name)
			}
		}
		// Handle nested structs recursively.
		if field.Kind() == reflect.Struct {
			if err := AssertNotNil(field.Interface()); err != nil {
				return fmt.Errorf("in nested struct %s: %w", fieldType.Name, err)
			}
		}
	}
	if len(nilFields) > 0 {
		return fmt.Errorf("nil fields found: %v", nilFields)
	}
	return nil
}

/*
type Person struct {
    Name    *string
    Age     *int
    Address *string
}

func main() {
    name := "John"
    age := 30

    // This will return an error because Address is nil
    person := Person{
        Name: &name,
        Age:  &age,
        Address: nil,
    }

    if err := AssertNotNil(person); err != nil {
        fmt.Println(err) // Output: nil fields found: [Address]
    }
}
*/
