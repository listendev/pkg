// Package manifest provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package manifest

// Defines values for Manifest.
const (
	None        Manifest = ""
	PackageJSON Manifest = "package.json"
)

// Manifest defines model for Manifest.
type Manifest string
