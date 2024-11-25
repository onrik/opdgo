# opdgo

[![Go Report Card](https://goreportcard.com/badge/github.com/onrik/opdgo)](https://goreportcard.com/report/github.com/onrik/opdgo)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/onrik/opdgo)](https://pkg.go.dev/github.com/onrik/opdgo)

Golang client for [OpenPanel.dev](https://openpanel.dev) 


### Usage
```go
package main

import (
    "log"
	
    "github.com/onrik/opdgo"
)

func main() {
    client := opdgo.New("<client_id>", "<client_secret>", nil)

    client.Track("purchase_completed", {
        "product_id": "123",
        "price": 99.99,
        "currency": "USD"
    })
}

```