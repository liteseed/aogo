# Argo

Interact with [AO](https://ao.arweave.dev) in Go.

## Getting Started

```go
package main

import (
 "log"

 "github.com/liteseed/argo/ao"
 "github.com/liteseed/argo/signer"
 "github.com/liteseed/argo/transaction"
)

func main() {
 // Make a Signer
 s, err := signer.New("./data/wallet.json")
 if err != nil {
  log.Fatal(err)
 }

 // Initialize an AO Struct
 ao := ao.New()

 processId := "your-process-id"
 data := "your-data"
 tags := []transaction.Tag{{Name: "", Value: ""}}
 // Send a message to your AO process
 messageId, err := ao.SendMessage(processId, data, tags, s)
 if err != nil {
  log.Fatal(err)
 }

 // Read message result
 res, err := ao.ReadResult(processId, messageId)
 if err != nil {
  log.Fatal(err)
 }

 log.Println(res.Messages)
}

```
