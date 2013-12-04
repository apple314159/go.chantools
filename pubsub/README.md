# PubSub package for Go channels

[![GoDoc](http://godoc.org/github.com/satori/go.chantools/pubsub?status.png)](http://godoc.org/github.com/satori/go.chantools/pubsub)

This package provides pure Go implementation of Publish/Subscribe pattern for Go channels.

## Installation

Use the `go` command:

	$ go get github.com/satori/go.chantools/pubsub

## Example

```go
package main

import (
	"fmt"
	"github.com/satori/go.chantools/pubsub"
)

func main() {
	queue1 := make(chan int, 1)
	pubsub.Sub(queue1, "i_wanna_integer")
	defer pubsub.Unsub(queue1)

	pubsub.Pub(42, "i_wanna_integer")

	i := <-queue1

	fmt.Printf("Got int: %d", i)

	queue2 := make(chan string, 1)
	pubsub.Sub(queue2, "i_wanna_string")
	defer pubsub.Unsub(queue2)

	pubsub.Pub("fourty-two", "i_wanna_string")

	s := <-queue2

	fmt.Printf("Got string: %s", s)
}
```

## Documentation

[Documentation](http://godoc.org/github.com/satori/go.chantools/pubsub) is hosted at GoDoc project.

## Copyright

Copyright (C) 2013 by Maxim Bublis <b@codemonkey.ru>

PubSub package released under MIT License.
See [LICENSE](https://github.com/satori/go.chantools/blob/master/LICENSE) for details.
