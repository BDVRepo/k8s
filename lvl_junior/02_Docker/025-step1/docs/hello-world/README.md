# Hello world!

Our first program will print the classic “hello world” message. Here’s the full source code.

```golang
package main

import "fmt"

func main() {
    fmt.Println("hello world")
}
```

To run the program, put the code in hello-world.go and use go run.
	
```bash 
go run hello-world.go

#hello world
```

Sometimes we’ll want to build our programs into binaries. We can do this using go build.
	
```bash 
go build hello-world.go
ls

#hello-world    hello-world.go
```

We can then execute the built binary directly.
	
```bash 
./hello-world

#hello world
```

Now that we can run and build basic Go programs, let’s learn more about the language.