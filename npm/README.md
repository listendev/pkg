# npm

This package allows to get some information about a npm package from the npm registry.


## Installation

```
go env -w GOPRIVATE=github.com/garnet-org/*
go get github.com/garnet-org/pkg/npm
```


## Usage


### npm.GetPackageList

```go
pkg, err := npm.GetPackageList(context.Background(), "react")
if err != nil {
    panic(err)
}

fmt.Println(pkg.Name)
fmt.Println(pkg.Versions["0.0.1"].Dist.Shasum)
```

### npm.GetPackageVersion

```go
pkg, err := npm.GetPackageVersion(context.Background(), "react", "0.0.1")
if err != nil {
    panic(err)
}

fmt.Println(pkg.Name)
fmt.Println(pkg.Dist.Shasum)
```
