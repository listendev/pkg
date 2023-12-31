package category

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
)

// all is the list of all the supported categories.
var all = []Category{}

func init() {
	for i := range _Category_index {
		if i == len(_Category_index)-1 {
			break
		}
		all = append(all, Category(i+1))
	}
}

func Categories() []Category {
	return all
}

func FromUint64(input uint64) (Category, error) {
	for _, c := range all {
		if uint64(c) == input {
			return c, nil
		}
	}

	return 0, fmt.Errorf("couldn't find a category matching input (%d)", input)
}

func FromString(input string) (Category, error) {
	for _, c := range all {
		if c.String() == input || string(c.Case()) == input || strings.EqualFold(c.String(), input) {
			return c, nil
		}
	}

	return 0, fmt.Errorf("couldn't find a category matching input %q", input)
}

// Scan implements the sql.Scanner interface.
func (c *Category) Scan(source any) error {
	if x, ok := source.(string); ok {
		cat, err := FromString(x)
		if err != nil {
			return err
		}
		*c = cat

		return nil
	}
	if x, ok := source.(uint64); ok {
		cat, err := FromUint64(x)
		if err != nil {
			return err
		}
		*c = cat
	}

	return fmt.Errorf("cannot scan %T into Category", source)
}

func (c *Category) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("couldn't unmarshal category data %q", string(data))
	}

	o, e := FromString(str)
	if e != nil {
		return e
	}
	*c = o

	return nil
}

func (c Category) MarshalJSON() ([]byte, error) {
	s := c.Case()

	return json.Marshal(s)
}

func (c Category) Case() Case {
	s := c.String()
	delim := strcase.ToDelimited(s, ' ')

	return Case(delim)
}

type Case string

func (x Case) Original() string {
	return strcase.ToCamel(string(x))
}

func init() {
	strcase.ConfigureAcronym("cve", "CVE")
	strcase.ConfigureAcronym("cis", "CIS")
}

//go:generate stringer -type=Category
