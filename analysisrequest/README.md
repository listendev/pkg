# Analysis Request

This package provides a way to create the request to analyze a package.

## Installation

```
go env -w GOPRIVATE=github.com/garnet-org/*
go get github.com/garnet-org/pkg/analysisrequest
```

## Usage

### Unmarshal a request from JSON

```go
package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/garnet-org/pkg/analysisrequest"
	"github.com/garnet-org/pkg/npm"
	"github.com/garnet-org/pkg/observability"
)

func main() {
	arJSON := `{"type": "urn:scheduler:falco!npm,install.json", "snowflake_id": "1524854487523524608", "name": "chalk"}`
	// you can use the observability package to create a context with tracing and logging here
	ctx := observability.NewNopContext()
	arbuilder, _ := analysisrequest.NewBuilder(ctx)
	regClient, _ := npm.NewNPMRegistryClient(npm.NPMRegistryClientConfig{})
	arbuilder.WithNPMRegistryClient(regClient)
	ar, _ := arbuilder.FromJSON([]byte(arJSON))
	spew.Dump(ar.(analysisrequest.NPM))
}
```
