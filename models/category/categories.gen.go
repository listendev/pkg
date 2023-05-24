// Package category provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.5-0.20230513000919-14548c7e7bbe DO NOT EDIT.
package category

// Defines values for Category.
const (
	AdjacentNetwork Category = 8
	Advisory        Category = 7
	CIS             Category = 6
	Container       Category = 5
	Cybersquatting  Category = 11
	Filesystem      Category = 1
	Local           Category = 9
	Network         Category = 3
	Physical        Category = 10
	Process         Category = 2
	Users           Category = 4
)

// Category defines model for Category.
type Category uint64