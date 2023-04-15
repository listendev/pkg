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

	if len(typeFormats) != numTypes {
		panic("some type is missing its format definition")
	}
}

type TypeComponents struct {
	Framework string
	Collector string
	Ecosystem Ecosystem
	Actions   []string
	Format    string
}

// typeURNs maps the enums to their string representations.
//
// We use strings that are URNs (see https://www.rfc-editor.org/rfc/rfc2141).
// The format is urn:<framework>:<collector>[!<ecosystem>[.<action>]{0,}]
//
// Notice only the framework part is case-insensitive.
var typeURNs = map[Type]string{
	Nop:                  "urn:NOP:nop",
	NPMInstallWhileFalco: "urn:scheduler:falco!npm.install",
	NPMTestWhileFalco:    "urn:scheduler:falco!npm.test",
	NPMDepsDev:           "urn:hoarding:depsdev!npm",
	// FIXME: we need a way to represent enriching collectors
	// EnrichFalcoAlertsWithGPT: {"urn:hoarding:enrichfalcoalertswithgpt"},
}

var typeFormats = map[Type]string{
	Nop:                  "json",
	NPMInstallWhileFalco: "json",
	NPMTestWhileFalco:    "json",
	NPMDepsDev:           "json",
}

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

func (t Type) Components() TypeComponents {
	// Assuming we do not forget to correctly define any type...
	u, _ := urn.Parse([]byte(t.String()))
	n := u.Normalize()

	others := strings.Split(n.SS, "!")
	ecosystem := ""
	actions := []string{}
	if len(others) > 1 {
		rest := strings.Split(others[1], ".")
		ecosystem = rest[0]
		if len(rest) > 1 {
			actions = rest[1:]
		}
	}

	tc := TypeComponents{
		Framework: n.ID,
		Collector: others[0],
		Ecosystem: Ecosystem(ecosystem),
		Actions:   actions,
		Format:    typeFormats[t],
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
	filename := c.Collector
	suffix := strings.Join(c.Actions, ",")
	if len(suffix) > 0 {
		filename = fmt.Sprintf("%s(%s)", filename, suffix)
	}
	filename += "." + strings.TrimPrefix(c.Format, ".")

	return filename
}

func (c TypeComponents) HasEcosystem() bool {
	return c.Ecosystem != ""
}

func (c TypeComponents) HasActions() bool {
	return len(c.Actions) > 0
}

func (c TypeComponents) ToURN() *urn.URN {
	u := &urn.URN{
		ID: c.Framework,
		SS: c.Collector,
	}
	if c.HasEcosystem() {
		u.SS += "!" + string(c.Ecosystem)
	}
	if c.HasActions() {
		for _, action := range c.Actions {
			u.SS += "." + action
		}
	}

	return u
}
