

# FP Go package
  This package features some basic FP functions like Filter, Map, Reduce, Set and others.
## List of functions
###  Filter
```go
// Predicate type takes any type of slice  and a function which returns a bool
type Predicate[A any] func(A) bool

// Filter takes as input a slice of Any and a Predicate type function to be applied to each item in the slice.
// returns a new slice containing only the satisfied elements
func Filter[A any](input []A, pred Predicate[A]) []A  
```
```go
nums := []int{1, 2, 3, 4, 5}
even := Filter(nums, func(n int) bool { return n%2 == 0 })
// even -> [2, 4]
```
### Map
```go

// MapFunc type of Any func(A any) A
type MapFunc[A any] func(A) A

// Map applies MapFunc to each element of the input slice
// returns a new modified slice with the same number of elements
func Map[A any](input []A, m MapFunc[A]) ([]A, error)
```
```go
nums := []int{1, 2, 3}
squared, _ := Map(nums, func(n int) int { return n * n })
// squared -> [1, 4, 9]
```
### FlatMap
```go
// FlatMap is a one-to-many map. Example below
// - input slice of any type
// - func(A) []A to be applied to each element
func FlatMap[A any](input []A, m func(A) []A) []A
```
```go
words := []string{"go", "lang"}
chars := FlatMap(words, func(s string) []string { return strings.Split(s, "") })
// chars -> ["g", "o", "l", "a", "n", "g"]
```
### Reduce
```go
// type reduceFunc applies func to successive elements.
type reduceFunc[A any] func(a1, a2 A) A

func Reduce[A any](input []A, reducer reduceFunc[A]) A
```
```go
words := []string{"Go", "is", "awesome"}
sentence := Reduce(words, func(a, b string) string { return a + " " + b })
// sentence -> "Go is awesome"
```
### Sum and Product ( examples of reduce)
```go
func Sum[A Number](input []A) A {
	return Reduce(input, func(a1, a2 A) A { return a1 + a2 })
}
func Product[A Number](input []A) A {
	return Reduce(input, func(a1, a2 A) A { return a1 * a2 })
}
```
```go
nums := []int{1, 2, 3, 4}
total := Sum(nums)       // total -> 10
product := Product(nums) // product -> 24
```
### Set
```go
// Set function that works with any slice (basic types or structs)
// returns nil on empty input.
func Set[T comparable](slice []T) []T
```
```go
nums := []int{1, 2, 2, 3}
unique := Set(nums)
```
### Union
```go
// Union takes 2 sets of comparable types and returns their union.
// returns nil else

func Union[T comparable](set1, set2 []T) []T
```
```go
a := []int{1, 2}
b := []int{2, 3}
union := Union(a, b)
// union -> [1, 2, 3]
```
### Intersection
```go
// Intersection returns a slice containing the intersection of two slices.
// returns nil if there is no intersection between set1 and set2.

func Intersection[T comparable](set1, set2 []T) []T 
```
```go
a := []int{1, 2}
b := []int{2, 3}
inter := Intersection(a, b)
// inter -> [2]
```
### Difference
```go
// Difference returns a slice containing elements in set1 but not in set2.
// returns nil if sets are equal
func Difference[T comparable](set1, set2 []T) []T 
```
```go
a := []int{1, 2, 3}
b := []int{2, 3}
diff := Difference(a, b)
// diff -> [1]
```
### Exists
```go

// Exists takes in any array and returns true if predicate is true for any element.
//
//	-type Predicate func(A any) bool
func Exists[A any](input []A, pred Predicate[A]) bool
```
```go
nums := []int{1, 2, 3}
hasEven := Exists(nums, func(n int) bool { return n%2 == 0 })
```
### For all
```go
// ForAll returns true if all elements satisfy predicate.
func ForAll[A any](input []A, pred Predicate[A]) bool {
```
```go
nums := []int{2, 4, 6}
allEven := ForAll(nums, func(n int) bool { return n%2 == 0 })
// allEven -> true
```
## Logical operators
### NOT
```go
func Not[A any](predicate func(A) bool) func(A) bool {
```
```go
notEven := Not(func(n int) bool { return n%2 == 0 })
odd := notEven(3)
// odd -> true
```
### AND
```go
func And[A any](second func(A) bool) func(func(A) bool) func(A) bool 
```
```go
greaterThanTwo := func(n int) bool { return n > 2 }
even := func(n int) bool { return n%2 == 0 }
andFunc := And(even)(greaterThanTwo)

result := andFunc(4)
// result -> true
```
### OR
```go
func Or[A any](second func(A) bool) func(func(A) bool) func(A) bool
```
```go
greaterThanTwo := func(n int) bool { return n > 2 }
even := func(n int) bool { return n%2 == 0 }
orFunc := Or(even)(greaterThanTwo)

result := orFunc(3)
// result -> true
```




