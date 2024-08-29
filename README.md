# gotask

a simple task queue based on golang

## Usage

```go
package main

import "github.com/sumwai/gotask"


type Test struct {
	Name string
}

func main() {
	task := gotask.New()
	task.DebugMode(true)

	task.AddTask("hello", func(params ...any) any {
		var a, b, c string
		var d int
		var e Test
		if err := gotask.Params(params).Parse(&a, &b, &c, &d, &e); err != nil {
			return gotask.Exit{Code: -1, Message: err.Error()}
		}
		return gotask.Params{
			1,
			2,
			3,
		}
	})

	task.AddTask("world", func(params ...any) any {
		var a, b, c int
		if err := gotask.Params(params).Parse(&a, &b, &c); err != nil {
			return gotask.Exit{Code: -1, Message: err.Error()}
		}
		return nil
	})

    // return `Exit.Data` or `final result`
	task.Run("1", "2", "3", 4, Test{Name: "test"})
}

```

```text
2024/08/29 18:07:36 [Task] [Start] hello
2024/08/29 18:07:36 [Task] [Params] [1 2 3 4 {test}]
2024/08/29 18:07:36 [Task] [Done] hello
2024/08/29 18:07:36 [Task] [Result] [1 2 3]
2024/08/29 18:07:36 [Task] 
2024/08/29 18:07:36 [Task] [Start] world
2024/08/29 18:07:36 [Task] [Params] [1 2 3]
2024/08/29 18:07:36 [Task] [Done] world
2024/08/29 18:07:36 [Task] [Result] <nil>
2024/08/29 18:07:36 [Task] 
```

## Struct

### Params

Params is a slice of any type, It has a method `Parse` to parse the parameters into the specified types.

```go
type Params []any

func (p Params) Parse(args ...any) error
```

### Example

```go
package main

import "fmt"

func main() {
	var a, b, c int
	if err := (Params{1, 2, 3}).Parse(&a, &b, &c); err != nil {
		panic(err)
	}
	fmt.Println(a, b, c)
}
```

```
1 2 3
```

### Errors
#### the number of parameters is not equal to the number of specified types.
```go
func main() {
    var a, b int
    if err := (Params{1}).Parse(&a, &b); err != nil {
        fmt.Println(err)
    }
    fmt.Println(a, b)
}
```
> params length mismatch, want: 2, but: 1

#### the specified type is not a pointer.
```go
func main() {
    var a int
    if err := (Params{1}).Parse(a); err != nil {
        fmt.Println(err)
    }
    fmt.Println(a)
}
```
> params #0 is not a pointer

#### the specified type is not match the parameter type.

```go
func main() {
	var a string
	if err := (Params{1}).Parse(&a); err != nil {
		fmt.Println(err)
	}
	fmt.Println(a)
}
```
> params #0 type mismatch:string != int

#### cannot set

```go
// I don't know how to trig this error.
// that's code is
if !paramsv.Elem().CanSet() {
    return fmt.Errorf("params #" + strconv.Itoa(i) + " cannot set")
}
```

> params #0 cannot set

### Exit

```go
type Exit struct {
	Code    int
	Message string
    Data    any
}
``` 

You can use this to exit the task and return a result.

```go
return gotask.Exit{Code: -1, Message: err.Error()}
```


## License

MIT