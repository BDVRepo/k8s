# Values

Go has various value types including strings, integers, floats, booleans, etc. Here are a few basic examples.

```golang
package main

import "fmt"

func main() {

    fmt.Println("go" + "lang")

    fmt.Println("1+1 =", 1+1)
    fmt.Println("7.0/3.0 =", 7.0/3.0)

    fmt.Println(true && false)
    fmt.Println(true || false)
    fmt.Println(!true)
}
```
Strings, which can be added together with +.
	
```golang
    fmt.Println("go" + "lang")
```    

Integers and floats.
	
```golang
    fmt.Println("1+1 =", 1+1)
    fmt.Println("7.0/3.0 =", 7.0/3.0)
```    

Booleans, with boolean operators as you’d expect.