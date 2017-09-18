# WealthEngine Go

A fully featured WealthEngine client delivered in a Golang SDK, for quickly getting up and running with the [WealthEngine REST API](https://dev.wealthengine.com/api).

# Getting Started
```bash
go get github.com/zackproser/wealthengine-go
```

```go
package main

import (
  "encoding/json"
  "fmt"

  "github.com/zackproser/wealthengine"
)

func main() {
  // Instantiate a new WealthEngine API client,
  // Passing in your APIKey and the environment you want to hit ("Prod" or "Dev")
  w := wealthengine.New("278s6sfd-z2cy-569h-f35g-7dwa4090g2", "Prod")

  // Score a profile by address
  score, scoreErr := w.ScoreOneByAddress("proser", "zachary", "3089 nunya", "business", "berkeley", "CA", "94704", "my-model")

  // Score a profile by email
  //score, scoreErr := w.ScoreOneByEmail("zackproser@gmail.com", "proser", "zachary", "my-model")

  // Score a profile by phone
  //score, scoreErr := w.ScoreOneByPhone("15103267023", "zachary", "proser", "my-model")

  if scoreErr != nil {
    fmt.Printf("scoreErr: %v", scoreErr)
  }

  scoreJson, mErr := json.Marshal(score)

  if mErr != nil {
    fmt.Printf("Score err: %v", mErr)
  }

  fmt.Printf("Score: \n\n%v", string(scoreJson))

    // Match a profile by email and name
    //profile, err := w.MatchOneByEmail("zackproser@gmail.com", "Zachary", "Proser", "full")

    // Match a profile by address
    //profile, err := w.MatchOneByAddress("proser", "zachary", "somewhere", "else", "berkeley", "CA", "94704", "basic")

    // Match a profile by phone
    profile, err := w.MatchOneByPhone("15553339021", "zachary", "proser", "basic")

    if err != nil {
      fmt.Printf("err: %v\n", err)
    }

    profileJson, mErr := json.Marshal(profile)

    if mErr != nil {
      fmt.Printf("Error marshaling JSON: %v\n", mErr)
    }

    fmt.Printf("json: \n%v", string(profileJson))

    // Prepare a batch processing request with several BatchLookups
    b := wealthengine.Batch{
      []wealthengine.BatchLookup{{"Proser", "Zach", "1234 st", "apt 1", "somecity", "CA", "94711", "zackproser@gmail.com", "1235103333"}, {"Someone", "Else", "1234 st", "apt 1", "Anaheim", "CA", "947111", "zackproser@gmail.com", "1235103333"}, {"Quabbity", "Ashwitz", "Wakka st", "apt 1", "Columbia", "MD", "947111", "zackproser@gmail.com", "47472727889"}},
    }

    fmt.Printf("Batch: %v\n", b)

    // Batch jobs return a job ID for looking up the status and results later
    batchJobId, err := w.FindMany(&b, "full")

    fmt.Printf("BatchID: %v", batchJobId)

    if batchJobId != nil {
      //Use the batch ID to check on the status
      status, statusErr := w.GetBatchJobStatus(batchJobId)

      if statusErr != nil {
        fmt.Printf("StatusErr: %v\n", statusErr)
      }

      fmt.Printf("Status: %v\n", status)

      // Use the batch ID to get the final job results
      results, finalErr := w.GetBatchJobResults(batchJobId)

      if finalErr != nil {
        fmt.Printf("finalErr: %v", finalErr)
      }

      resultsJson, resultsJsonErr := json.Marshal(results)

      if resultsJsonErr != nil {
        fmt.Printf("resultsJsonErr: %v", resultsJsonErr)
      }

      fmt.Printf("Results: %v", string(resultsJson))
    }
}
```

# Documentation

[Godocs](https://godoc.org/github.com/zackproser/wealthengine-go)