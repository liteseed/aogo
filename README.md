# AOGO

[![protocol.land](https://arweave.net/eZp8gOeR8Yl_cyH9jJToaCrt2He1PHr0pR4o-mHbEcY)](https://protocol.land/#/repository/d8e7b91b-1025-47a5-9ea8-364451f496f9)

Interact with [AO](https://ao.arweave.dev) in Go.

## Getting Started

```go
package main

import (
 "log"
 "os"

 "github.com/liteseed/goar/signer"
 "github.com/liteseed/goar/types"
 "github.com/liteseed/aogo"
)

// AO MODULE ID - Get Latest from "https://github.com/permaweb/aos"
const MODULE = "1PdCJiXhNafpJbvC-sjxWTeNzbf9Q_RfUNs84GYoPm0"

func main() {

 // Make a Signer
 s, err := signer.FromPath("./wallet.json")
 if err != nil {
  log.Fatal(err)
 }

 // Initialize an AO Struct
 ao, err := aogo.New()
 if err != nil {
  log.Fatal(err)
 }

 data := []byte{1, 2, 3}

 // Spawn a process with some data
 // Note: Spawn Process has some delay before the process shows up on ao.link and aos
 pID, err := ao.SpawnProcess(MODULE, data, []types.Tag{}, itemSigner)
 if err != nil {
  log.Fatal(err)
 }
 log.Println(pID)

  // Send a message to your AO process
  mID, err := ao.SendMessage("jysQej65l7KHRZi93csg0rvdmciJNL9hteM1N_yakpE", "", []types.Tag{{Name: "Action", Value: "Balance"}}, "", s)
  if err != nil {
  log.Fatal(err)
 }
 log.Println(mID)

 // Read message result
 res, err := ao.LoadResult("jysQej65l7KHRZi93csg0rvdmciJNL9hteM1N_yakpE", messageId)
 if err != nil {
  log.Fatal(err)
 }
 log.Println(res)

 // DryRun a Message and get it's result 
 dres, err := ao.DryRun(aogo.Message{Target: "qAbHIghbi7lb0Y8KdZr8q9dvH8xwXpawCXgGSD8OpJk", Data: "sgg2DZhTFPDIpI4bpZldRpP2RCobgtUstIlcxyh5mCA", Owner: s.Owner(), From: s.Address, Tags: []types.Tag{{Name: "Action", Value: "Balance"}}})
 if err != nil {
  log.Fatal(err)
 }
 log.Println(dres.Messages[0]["Data"])

```
