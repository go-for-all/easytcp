Simplified TCP with automatic message framing.

# The Problem
Go's standard `net.Dial` and `net.Listener.Accept` return a `net.Conn` in which you must pass an explicit size to read and write from. If you are using TCP for exchanging "messages", you have to manually handle length and backlogging to correctly parse your messages.

# The Solution
`easytcp` simplifies this by providing `Tunnel.ReadMessage()` and `Tunnel.WriteMessage()`, analogous to `net.Conn.Read` and `net.Conn.Write` respectively but `Tunnel.ReadMessage()` is guaranteed to return exactly what you passed in to `Tunnel.WriteMessage()`. It automatically handles message length and blocks until the full message is received.

# Usage
#### Write a message:
```go
package main

import (
	"time"

	"github.com/go-for-all/easytcp"
)

func main() {
	message := []byte("hello world")

	tunnel, err := easytcp.DialTimeout(":8000", 1*time.Second)
	if err != nil {
		panic(err)
	}
	if err := tunnel.WriteMessage(message); err != nil {
		panic(err)
	}
}
```
Note: message length must be greater than zero and less than 1 MiB (MaxMessageSize = 1 << 20).

#### Reading a message:
```go
package main

import (
	"fmt"

	"github.com/go-for-all/easytcp"
)

func main() {
	listener, err := easytcp.Listen(":8000")
	if err != nil {
		panic(err)
	}
	for {
		tunnel, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		message, err := tunnel.ReadMessage()
		if err != nil {
			panic(err)
		}
		fmt.Println(string(message)) // "hello world"

		if err := tunnel.Close(); err != nil {
			panic(err)
		}
	}
}
```

# License
All code in this repository is licensed under the GNU AGPLv3; you may not use this code unless you agree to the license.
