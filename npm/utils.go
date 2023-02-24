package npm

import "strings"

// SplitName returns the organization and the name part of an NPM package.
//
// Notice it assumes the package name is a valid one, thus it doesn't perform any validation.
func SplitName(packagename string) (string, string) {
	parts := strings.Split(packagename, "/")

	if len(parts) == 1 {
		return "", parts[0]
	}

	return parts[0], strings.Join(parts[1:], "")
}
