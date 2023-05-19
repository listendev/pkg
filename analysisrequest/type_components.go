package analysisrequest

import (
	"fmt"
	"strings"

	"github.com/leodido/go-urn"
)

type Framework string

const (
	None      Framework = "nop"
	Scheduler Framework = "scheduler"
	Hoarding  Framework = "hoarding"
)

type Collector string

const (
	NoCollector        Collector = "nop"
	FalcoCollector     Collector = "falco"
	DepsDevCollector   Collector = "depsdev"
	GPT4Collector      Collector = "gpt4"
	TyposquatCollector Collector = "typosquat"
)

type Ecosystem string

const (
	NPMEcosystem Ecosystem = "npm"
)

type TypeComponents struct {
	Framework        Framework
	Collector        Collector
	CollectorActions []string
	Ecosystem        Ecosystem
	EcosystemActions []string
	Format           string
	Parent           *TypeComponents
}

func (c TypeComponents) ResultFile() string {
	if c.Parent != nil {
		return c.Parent.ResultFile()
	}

	filename := string(c.Collector)

	suffix := ""

	colActionsSuffix := strings.Join(c.CollectorActions, ",")
	if len(colActionsSuffix) > 0 {
		suffix += fmt.Sprintf("(%s)", colActionsSuffix)
	}

	ecoActionsSuffix := strings.Join(c.EcosystemActions, ",")
	if len(ecoActionsSuffix) > 0 {
		suffix += fmt.Sprintf("[%s]", ecoActionsSuffix)
	}

	if c.Format != "" {
		suffix += "." + strings.TrimPrefix(c.Format, ".")
	}

	filename += suffix

	return filename
}

func (c TypeComponents) HasEcosystem() bool {
	return c.Ecosystem != ""
}

func (c TypeComponents) HasEcosystemActions() bool {
	return len(c.EcosystemActions) > 0
}

func (c TypeComponents) HasCollectorActions() bool {
	return len(c.CollectorActions) > 0
}

func (c TypeComponents) HasFormat() bool {
	return len(c.Format) > 0
}

func (c TypeComponents) ToURN() *urn.URN {
	u := &urn.URN{
		ID: string(c.Framework),
		SS: string(c.Collector),
	}
	if c.HasCollectorActions() {
		for _, action := range c.CollectorActions {
			u.SS += "," + action
		}
	}
	if c.Parent == nil {
		if c.HasEcosystem() {
			u.SS += "!" + string(c.Ecosystem)
		}
		if c.HasEcosystemActions() {
			for _, action := range c.EcosystemActions {
				u.SS += "," + action
			}
		}
		if c.HasFormat() {
			u.SS += "." + c.Format
		}
	} else {
		p := c.Parent.ToURN()
		if p == nil {
			return p
		}
		p.SS += "+" + u.String()
		p.Normalize()

		return p
	}

	return u
}
