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
arJSON := `{"type": "npm", "snowflake_id": "1524854487523524608", "name": "chalk"}`
// you can use the observability package to create a context with tracing and logging here
ctx := observability.NewNopContext()
arbuilder := analysisrequest.NewAnalysisRequestBuilder()
regClient := npm.NewNPMRegistryClient(npm.NPMRegistryClientConfig{})
arbuilder.WithNPMRegistryClient(regClient)
ar, _ := arbuilder.FromJSON(ctx, arJSON)
spew.Dump(ar)
```
