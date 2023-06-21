package analysisrequest

import (
	"fmt"
	"strings"
)

type Ecosystem string

const (
	NPMEcosystem Ecosystem = "npm"
)

var (
	allEcosystems = []Ecosystem{
		NPMEcosystem,
	}
)

func GetEcosystemFrom(input string) (Ecosystem, error) {
	x := strings.ToLower(input)
	switch x {
	case "npm":
		return NPMEcosystem, nil
	default:
		return "", fmt.Errorf("couldn't find an ecosystem matching the input string %q", input)
	}
}
