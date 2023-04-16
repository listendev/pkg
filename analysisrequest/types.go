package analysisrequest

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/leodido/go-urn"
)

type Ecosystem string

const (
	NPMEcosystem Ecosystem = "npm"
)

type Type int

const (
	Nop Type = iota + 1
	NPMInstallWhileFalco
	NPMTestWhileFalco
	NPMDepsDev
	NPMGPT4InstallWhileFalco

	_maxType
)

func MaxType() Type {
	return _maxType - 1
}

func init() {
	numTypes := int(MaxType())
	if len(typeURNs) != numTypes {
		panic("some type is missing its URN definition")
	}
}

type TypeComponents struct {
	Framework        string
	Collector        string
	CollectorActions []string
	Ecosystem        Ecosystem
	EcosystemActions []string
	Format           string
	Parent           *TypeComponents
}

// typeURNs maps the enums to their string representations.
//
// We use strings that are URNs (see https://www.rfc-editor.org/rfc/rfc2141).
// The format is urn:<framework>:<collector[,<action>{0,}]>[!<ecosystem>[,<action>]{0,}].<format>
//
// What does each component mean?
//   - <framework> is the framework/platform meant to process the analysis request with the current type.
//   - <collector> is the collector such a platform is meant to execute.
//   - <ecosystem> represents the ecosystem (ie., language/package manager) the analysis request refers to.
//   - <action>{0,} are the specific actions of the ecosystem that the collector will execute.
//     or the actions of the collector itself.
//
// Notice only the framework part is case-insensitive.
var typeURNs = map[Type]string{
	Nop:                  "urn:NOP:nop",
	NPMInstallWhileFalco: "urn:scheduler:falco!npm,install.json",
	NPMTestWhileFalco:    "urn:scheduler:falco!npm,test.json",
	NPMDepsDev:           "urn:hoarding:depsdev!npm.json",
	// NPMGPT4InstallWhileFalco represents analysis requests to enrich the NPMInstallWhileFalco results
	NPMGPT4InstallWhileFalco: "urn:scheduler:falco!npm,install.json+urn:hoarding:gpt4,context",
}

// TODO: enforce enrichers (+urn:...) to do not specify ecosystem, ecosystem actions, and format

func ToType(s string) (Type, error) {
	uu, ok := urn.Parse([]byte(s))
	if !ok {
		return 0, fmt.Errorf("cannot convert non URN input to types")
	}
	uuu := uu.Normalize()

	for t := range typeURNs {
		u := t.ToURN()
		if u != nil && u.ID == uuu.ID && u.SS == uuu.SS {
			return t, nil
		}
	}

	return 0, fmt.Errorf("couldn't convert %s to any type", s)
}

func (t Type) ToURN() *urn.URN {
	return t.Components().ToURN()
}

func (t Type) HasEcosystem() bool {
	return t.Components().Ecosystem != ""
}

func (t Type) String() string {
	// Assuming we do not forget to correctly define any type...
	representation := typeURNs[t]

	return representation
}

func (t Type) Parent() (Type, error) {
	c := t.Components()
	if c.Parent != nil {
		return ToType(c.Parent.ToURN().String())
	}

	return 0, fmt.Errorf("type %q isn't an enricher, thus it doesn't have a parent type", t.String())
}

func (t Type) Components() TypeComponents {
	// Assuming we do not forget to correctly define any type...
	s := t.String()
	u, _ := urn.Parse([]byte(s))
	n := u.Normalize()

	// Is this an enricher?
	enrich := false
	cc := TypeComponents{}
	hierarchy := strings.Split(n.SS, "+")
	n.SS = hierarchy[0]
	if len(hierarchy) > 1 {
		cu, _ := urn.Parse([]byte(hierarchy[1]))
		cc = componentsFromString(cu.Normalize())
		enrich = true
	}

	c := componentsFromString(n)

	ret := TypeComponents{
		Framework:        c.Framework,
		Collector:        c.Collector,
		CollectorActions: c.CollectorActions,
		Ecosystem:        c.Ecosystem,
		EcosystemActions: c.EcosystemActions,
		Format:           c.Format,
	}
	if enrich {
		ret.Parent = &c
	}

	// Override (some) parent values with child ones
	// When these conditions met it means the current type is an enricher
	// Thus, it only makes sense to override the framework and the collector values where the augumenting will happen
	// while keeping the parent's ecosystem and format info
	if cc.Framework != "" && cc.Framework != c.Framework {
		ret.Framework = cc.Framework
	}
	if cc.Collector != "" && cc.Collector != c.Collector {
		ret.Collector = cc.Collector
		ret.CollectorActions = cc.CollectorActions
	}
	if cc.Format != "" && cc.Format != c.Format {
		ret.Format = cc.Format
	}

	return ret
}

// componentsFromString parses the SS part of the custom URN format we use.
//
// The format of the SS part is: <collector[,<action>{0,}]>[!<ecosystem>[,<action>]{0,}][.<format>]
func componentsFromString(n *urn.URN) TypeComponents {
	// .<format>
	firstSplit := strings.Split(n.SS, ".")
	format := ""
	if len(firstSplit) > 1 {
		format = firstSplit[1]
	}

	// <collector[,<action>{0,}]>
	remainings := strings.Split(firstSplit[0], "!")
	collector := remainings[0]
	cActions := []string{}
	firstRemainingsActions := strings.Split(remainings[0], ",")
	if len(firstRemainingsActions) > 1 {
		collector = firstRemainingsActions[0]
		cActions = firstRemainingsActions[1:]
	}

	// <ecosystem>[,<action>]{0,}
	ecosystem := ""
	eActions := []string{}
	if len(remainings) > 1 {
		rest := strings.Split(remainings[1], ",")
		ecosystem = rest[0]
		if len(rest) > 1 {
			eActions = rest[1:]
		}
	}

	tc := TypeComponents{
		Framework:        n.ID,
		Collector:        collector,
		CollectorActions: cActions,
		Ecosystem:        Ecosystem(ecosystem),
		EcosystemActions: eActions,
		Format:           format,
	}

	return tc
}

func (t Type) MarshalJSON() ([]byte, error) {
	u := t.ToURN()
	if u == nil {
		return nil, fmt.Errorf("couldn't marshal because type is not an URN")
	}

	return json.Marshal(u.String())
}

func (t *Type) UnmarshalJSON(data []byte) error {
	var typeStr string
	if err := json.Unmarshal(data, &typeStr); err != nil {
		return err
	}
	res, err := ToType(typeStr)
	if err != nil {
		return err
	}
	*t = res

	return nil
}

func (c TypeComponents) ResultFile() string {
	if c.Parent != nil {
		return c.Parent.ResultFile()
	}

	filename := c.Collector

	suffix := ""

	colActionsSuffix := strings.Join(c.CollectorActions, ",")
	if len(colActionsSuffix) > 0 {
		suffix += fmt.Sprintf("(%s)", colActionsSuffix)
	}

	ecoActionsSuffix := strings.Join(c.EcosystemActions, ",")
	if len(ecoActionsSuffix) > 0 {
		suffix += fmt.Sprintf("[%s]", ecoActionsSuffix)
	}

	filename += suffix + "." + strings.TrimPrefix(c.Format, ".")

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
		ID: c.Framework,
		SS: c.Collector,
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
