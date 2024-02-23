package stringutil

import (
	"regexp"
)

const (
	// emailPattern matches email addresses (also IDN) into text.
	// Using it ^ at the beginning (asserting start of line position) would make it only extract valid emails
	// Anyways we can't be sure the input
	emailPattern = `(?i)(([^<>()\[\]\.,;:\s@\"]+(\.[^<>()\[\]\.,;:\s@\"]+)*)|(\".+\"))@(([^<>()[\]\.,;:\s@\"]+\.)+[^<>()[\]\.,;:\s@\"]{2,})`
)

var (
	emailRegex = regexp.MustCompile(emailPattern)
)

func MatchEmails(s string) []string {
	res := emailRegex.FindAllString(s, -1)

	return res
}
