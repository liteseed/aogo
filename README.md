# AOGO

[![protocol.land](https://arweave.net/eZp8gOeR8Yl_cyH9jJToaCrt2He1PHr0pR4o-mHbEcY)](https://protocol.land/#/repository/d8e7b91b-1025-47a5-9ea8-364451f496f9)

Interact with [AO](https://ao.arweave.dev) in Go.

## Getting Started

```go
package main

import (
 "log"
 "os"

 "github.com/everFinance/goar"
 "github.com/everFinance/goar/types"
 "github.com/liteseed/aogo"
)

func main() {

 // Make a Signer
 s, err := goar.NewSignerFromPath("./keys/wallet.json")
 if err != nil {
  log.Fatal(err)
 }
 itemSigner, err := goar.NewItemSigner(s)
 // Initialize an AO Struct
 if err != nil {
  log.Fatal(err)
 }

 // Initialize an AO Struct
 ao, err := aogo.New()
 if err != nil {
  log.Fatal(err)
 }

 data := []byte{1, 2, 3}

 //Spawn a process with some data
 processId, err := ao.SpawnProcess(data, []types.Tag{}, itemSigner)
 if err != nil {
  log.Fatal(err)
 }
 log.Println(processId)

 // Send a message to your AO process
 messageId, err := ao.SendMessage(processId, data, []types.Tag{{Name: "Action", Value: "Eval"}}, "", itemSigner)
 if err != nil {
  log.Fatal(err)
 }
 log.Println(messageId)
 // Read message result
 res, err := ao.ReadResult(processId, messageId)
 if err != nil {
  log.Fatal(err)
 }
 log.Println(res)
}

```
