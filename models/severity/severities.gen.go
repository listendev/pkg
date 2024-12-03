// Package severity provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version 2.4.1 DO NOT EDIT.
package severity

// Defines values for Severity.
const (
	Empty  Severity = ""
	High   Severity = "high"
	Low    Severity = "low"
	Medium Severity = "medium"
)

// Severity defines model for Severity.
type Severity string
