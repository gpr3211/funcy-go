# FP Monad Package
This package provides implementations of common functional programming monads in Go, including `Maybe` and `Future`.

## List of Functions

### Maybe Monad
The `Maybe` monad represents an optional value. It can either be `Just` a value or `Nothing`.

#### `Just` and `Nothing`
```go
func Just[A any](a A) JustMaybe[A]
func Nothing[A any]() Maybe[A]
```
```go
maybeValue := Just(42)
noValue := Nothing[int]()
fmt.Println(maybeValue.GetOrElse(0)) // Output: 42
fmt.Println(noValue.GetOrElse(0))    // Output: 0
```

#### `Get` and `GetOrElse`
```go
func (j JustMaybe[A]) Get() A
func (j JustMaybe[A]) GetOrElse(def A) A
func (n NothingMaybe[A]) Get() A
func (n NothingMaybe[A]) GetOrElse(def A) A
```
```go
name := Just("Alice")
nilName := Nothing[string]()
fmt.Println(name.GetOrElse("Unknown"))  // Output: Alice
fmt.Println(nilName.GetOrElse("Unknown")) // Output: Unknown
```

#### `fmap`
```go
func fmap[A, B any](m Maybe[A], mapFunc func(A) B) Maybe[B]
```
```go
num := Just(10)
doubled := fmap(num, func(n int) int { return n * 2 })
fmt.Println(doubled.GetOrElse(-1)) // Output: 20
```

---

### Future Monad
The `Future` monad represents a computation that may complete asynchronously.

#### `NewFuture`
```go
func NewFuture[A any](compute func() (A, error)) *Future[A]
```
```go
future := NewFuture(func() (int, error) {
    time.Sleep(1 * time.Second)
    return 42, nil
})
result, err := future.Get()
fmt.Println(result, err) // Output: 42 <nil>
```

#### `Get` and `GetWithTimeout`
```go
func (f *Future[A]) Get() (A, error)
func (f *Future[A]) GetWithTimeout(timeout time.Duration) (A, error,bool)
// returns A,nil,false if Success
// returns A,err,false if error in computation
// returns A,err,true if timout
// 
```
```go
result, err,timeOut := future.GetWithTimeout(2 * time.Second)
fmt.Println(result, err,timeOut)
```

#### `Map`
```go
func Map[A, B any](f *Future[A], fn func(A) B) *Future[B]
```
```go
doubled := Map(future, func(x int) int { return x * 2 })
fmt.Println(doubled.Get()) // Output: 84 <nil>
```

#### `FlatMap`
```go
func FlatMap[A, B any](f *Future[A], fn func(A) *Future[B]) *Future[B]
```
```go
resultFuture := FlatMap(future, func(x int) *Future[string] {
    return NewFuture(func() (string, error) {
        return fmt.Sprintf("Result is %d", x), nil
    })
})
fmt.Println(resultFuture.Get()) // Output: "Result is 42" <nil>
```

#### `Successful` and `Failed`
```go
func Successful[A any](value A) *Future[A]
func Failed[A any](err error) *Future[A]
```
```go
successFuture := Successful(100)
errorFuture := Failed[int](fmt.Errorf("an error occurred"))
fmt.Println(successFuture.Get()) // Output: 100 <nil>
fmt.Println(errorFuture.Get())    // Output: 0 an error occurred
```

#### `Sequence`
```go
func Sequence[A any](futures ...*Future[A]) *Future[[]A]
```
```go
future1 := Successful(10)
future2 := Successful(20)
combined := Sequence(future1, future2)
fmt.Println(combined.Get()) // Output: [10 20] <nil>
```

## More Complete examples

```go
// Safe map access
func getFromMap[K comparable, V any](m map[K]V, key K) Maybe[V] {
	if value, ok := m[key]; ok {
		return Just(value)
	}
	return Nothing[V]()
}
```
 - 
```go
 	data := map[string]int{
		"foo": 42,
		"bar": 84,
	}

	lookup1 := getFromMap(data, "foo")
	lookup2 := getFromMap(data, "nonexistent")

	fmt.Printf("Lookup foo: %v\n", lookup1.GetOrElse(-1))
	fmt.Printf("Lookup nonexistent: %v\n", lookup2.GetOrElse(-1))
    
```

```go
// Example types and helper functions
type User struct {
	Name string
	Age  int
}
func getUserByID(id int) Maybe[User] {
	users := map[int]User{
		1: {"Alice", 30},
		2: {"Bob", 25},
	}
	if user, exists := users[id]; exists {
		return Just(user)
	}
	return Nothing[User]()
}
```
 -  

```go
   func parseNumber(s string) Maybe[int] {
	if num, err := strconv.Atoi(s); err == nil {
		return Just(num)
	}
	return Nothing[int]()
}     
```
```go
	parseAndDouble := func(s string) Maybe[int] {
		return fmap(parseNumber(s), func(n int) int {
			return n * 2
		})
	}
	num1 := parseAndDouble("10")
	num2 := parseAndDouble("invalid")

	fmt.Printf("Parsed and doubled 10: %v\n", num1.GetOrElse(-1))
	fmt.Printf("Parsed and doubled invalid: %v\n", num2.GetOrElse(-1))

```
 -  safe access for nullable values.

```go
func fromNullable[A any](ptr *A) Maybe[A] {
	if ptr == nil {
		return Nothing[A]()
	}
	return Just(*ptr)
}
```
```go
	name := "John"
	maybeNamePtr := fromNullable(&name)
	fmt.Printf("Name: %v\n", maybeNamePtr.GetOrElse("Unknown"))

	var nilName *string
	maybeNilPtr := fromNullable(nilName)
	fmt.Printf("Nil name: %v\n", maybeNilPtr.GetOrElse("Unknown"))
```

