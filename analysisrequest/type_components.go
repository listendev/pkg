package analysisrequest

import (
	"fmt"
	"strings"

	"github.com/leodido/go-urn"
	"github.com/listendev/pkg/ecosystem"
)

type Framework string

const (
	None      Framework = "nop"
	Scheduler Framework = "scheduler"
	Hoarding  Framework = "hoarding"
)

type Collector string

const (
	NoCollector                     Collector = "nop"
	DynamicInstrumentationCollector Collector = "dynamic"
	AdvisoryCollector               Collector = "advisory"
	AICollector                     Collector = "ai"
	TyposquatCollector              Collector = "typosquat"
	MetadataCollector               Collector = "metadata"
	StaticAnalysisCollector         Collector = "static"
)

type TypeComponents struct {
	Framework       Framework
	Collector       Collector
	CollectorAction string
	Ecosystem       ecosystem.Ecosystem
	EcosystemAction string
	Format          string
	Parent          *TypeComponents
}

// ResultFile returns the filename of the result file for the current Components.
//
// Note it tries to always use characters safe for S3 keys (see https://docs.aws.amazon.com/AmazonS3/latest/userguide/object-keys.html).
func (c TypeComponents) ResultFile() string {
	if c.Parent != nil {
		return c.Parent.ResultFile()
	}

	filename := string(c.Collector)

	suffix := ""

	if len(c.CollectorAction) > 0 {
		suffix += fmt.Sprintf("(%s)", c.CollectorAction)
	}

	if len(c.EcosystemAction) > 0 {
		suffix += fmt.Sprintf("!%s!", c.EcosystemAction)
	}

	if c.Format != "" {
		suffix += "." + strings.TrimPrefix(c.Format, ".")
	}

	filename += suffix

	return filename
}

func (c TypeComponents) HasEcosystem() bool {
	_, err := ecosystem.FromUint64(uint64(c.Ecosystem))

	return err == nil
}

func (c TypeComponents) HasEcosystemAction() bool {
	return len(c.EcosystemAction) > 0
}

func (c TypeComponents) HasCollectorAction() bool {
	return len(c.CollectorAction) > 0
}

func (c TypeComponents) HasFormat() bool {
	return len(c.Format) > 0
}

func (c TypeComponents) ToURN() *urn.URN {
	u := &urn.URN{
		ID: string(c.Framework),
		SS: string(c.Collector),
	}
	if c.HasCollectorAction() {
		u.SS += "," + c.CollectorAction
	}
	if c.Parent == nil {
		if c.HasEcosystem() {
			u.SS += "!" + c.Ecosystem.Case()
		}
		if c.HasEcosystemAction() {
			u.SS += "," + c.EcosystemAction
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
