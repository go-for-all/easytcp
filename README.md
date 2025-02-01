Simplified TCP with automatic message framing.

<hr>

General guidelines for this repository:
1. **License**: this code is AGPLv3-licensed; by using it, you agree to comply with the terms.
2. **Issues**: only bug reports will be reviewed; do not submit support or feature requests.
3. **Pull requests**: accepted only if they directly address a reported bug linked to an issue.

<hr>

# The Problem
Go's standard TCP returns `net.Conn` which requires you to pass an explicit data length to read.

Thus, you have to manually handle length and backlogging to correctly parse exchanged messages.

# The Solution
easytcp's `Tunnel` with `ReadMessage()` and `WriteMessage()` to replace `conn.Read` and `conn.Write`.

`Tunnel.ReadMessage()` is guaranteed to return exactly what you passed in to `Tunnel.WriteMessage()`.

easytcp automatically handles message length and will block until the full message is received.

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
**Note: message must be non-empty and no more than than 1 MiB (MaxMessageSize = 1 << 20).**

It is the caller's responsibility to chunk data into 1 MiB messages if necessary.

#### Read a message:
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
