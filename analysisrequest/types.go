package analysisrequest

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/garnet-org/pkg/ecosystem"
	"github.com/leodido/go-urn"
	"golang.org/x/exp/maps"
)

type Type int

// Those are the constants representing the analysis request types.
//
// When adding a new one append it after the existing ones (before the _maxType constant).
const (
	Nop Type = iota + 1
	NPMInstallWhileFalco
	NPMDepsDev
	NPMGPT4InstallWhileFalco
	// NPMTestWhileFalco
	NPMTyposquat
	NPMMetadataEmptyDescription
	NPMMetadataVersion
	NPMMetadataMaintainersEmailCheck

	NPMStaticAnalysisEnvExfiltration Type = iota + 10 // 18 // Do not forget to specify the type Type when using iota to reserve space for previous types
	NPMStaticAnalysisDetachedProcessExecution
	NPMStaticAnalysisShadyLinks
	NPMStaticAnalysisEvalBase64
	NPMStaticAnalysisInstallScript
	NPMStaticNonRegistryDependency

	_maxType
)

func LastType() Type {
	return _maxType - 1
}

func init() {
	maxID := LastType()

	typeIDs := Types()

	var maxType Type
	for _, t := range typeIDs {
		if t > maxType {
			maxType = t
		}
	}

	if maxType != maxID {
		panic("some type is missing its URN definition")
	}
}

func createType(f Framework, c Collector, cAction string, e ecosystem.Ecosystem, eAction string, format string) string {
	u := urn.URN{}
	if len(f) == 0 || len(c) == 0 {
		panic("Missing mandatory framework and collector names")
	}
	u.ID = string(f)
	u.SS = string(c)
	if len(cAction) > 0 {
		u.SS += fmt.Sprintf(",%s", cAction)
	}
	if e != 0 {
		u.SS += fmt.Sprintf("!%s", e.Case())
		if len(eAction) > 0 {
			u.SS += fmt.Sprintf(",%s", eAction)
		}
	}

	if len(format) > 0 {
		u.SS += fmt.Sprintf(".%s", format)
	}

	return u.String()
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
	Nop:                  createType(None, NoCollector, "", ecosystem.None, "", ""),
	NPMInstallWhileFalco: createType(Scheduler, FalcoCollector, "", ecosystem.Npm, "install", "json"),
	NPMDepsDev:           createType(Hoarding, DepsDevCollector, "", ecosystem.Npm, "", "json"),
	// NPMGPT4InstallWhileFalco represents analysis requests to enrich the NPMInstallWhileFalco results
	NPMGPT4InstallWhileFalco: "urn:scheduler:falco!npm,install.json+urn:hoarding:gpt4,context",
	// NPMTestWhileFalco:     "urn:scheduler:falco!npm,test.json",
	NPMTyposquat:                              createType(Hoarding, TyposquatCollector, "", ecosystem.Npm, "", "json"),
	NPMMetadataEmptyDescription:               createType(Hoarding, MetadataCollector, "empty_descr", ecosystem.Npm, "", "json"),
	NPMMetadataVersion:                        createType(Hoarding, MetadataCollector, "version", ecosystem.Npm, "", "json"),
	NPMMetadataMaintainersEmailCheck:          createType(Hoarding, MetadataCollector, "email_check", ecosystem.Npm, "", "json"),
	NPMStaticAnalysisEnvExfiltration:          createType(Hoarding, StaticAnalysisCollector, "exfiltrate_env", ecosystem.Npm, "", "json"),
	NPMStaticAnalysisDetachedProcessExecution: createType(Hoarding, StaticAnalysisCollector, "detached_process_exec", ecosystem.Npm, "", "json"),
	NPMStaticAnalysisShadyLinks:               createType(Hoarding, StaticAnalysisCollector, "shady_links", ecosystem.Npm, "", "json"),
	NPMStaticAnalysisEvalBase64:               createType(Hoarding, StaticAnalysisCollector, "base64_eval", ecosystem.Npm, "", "json"),
	NPMStaticAnalysisInstallScript:            createType(Hoarding, StaticAnalysisCollector, "install_script", ecosystem.Npm, "", "json"),
	NPMStaticNonRegistryDependency:            createType(Hoarding, StaticAnalysisCollector, "non_registry_dependency", ecosystem.Npm, "", "json"),
}

func Types() []Type {
	return maps.Keys(typeURNs)
}

// TODO: enforce types to have max 1 collector action and max 1 ecosystem action
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
	return t.Components().HasEcosystem()
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
		Framework:       c.Framework,
		Collector:       c.Collector,
		CollectorAction: c.CollectorAction,
		Ecosystem:       c.Ecosystem,
		EcosystemAction: c.EcosystemAction,
		Format:          c.Format,
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
		ret.CollectorAction = cc.CollectorAction
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
	cAction := ""
	firstRemainingsActions := strings.Split(remainings[0], ",")
	if len(firstRemainingsActions) > 1 {
		collector = firstRemainingsActions[0]
		cAction = firstRemainingsActions[1]
	}

	// <ecosystem>[,<action>]{0,}
	eco := ""
	eAction := ""
	if len(remainings) > 1 {
		rest := strings.Split(remainings[1], ",")
		eco = rest[0]
		if len(rest) > 1 {
			eAction = rest[1]
		}
	}

	// FIXME: assuming the eco string is always a valid ecosystem
	eee, _ := ecosystem.FromString(eco)

	tc := TypeComponents{
		Framework:       Framework(n.ID),
		Collector:       Collector(collector),
		CollectorAction: cAction,
		Ecosystem:       eee,
		EcosystemAction: eAction,
		Format:          format,
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
