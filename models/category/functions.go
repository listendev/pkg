package category

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (c Category) String() string {
	switch c {
	case AdjacentNetwork:
		fallthrough
	case Cis:
		fallthrough
	case Container:
		fallthrough
	case Cve:
		fallthrough
	case Filesystem:
		fallthrough
	case Local:
		fallthrough
	case Network:
		fallthrough
	case Physical:
		fallthrough
	case Process:
		fallthrough
	case Users:
		return string(c)
	}

	return "unknown"
}

func New(input string) (Category, error) {
	c := strings.ToLower(input)

	switch c {
	case AdjacentNetwork.String():
		return AdjacentNetwork, nil
	case Cis.String():
		return Cis, nil
	case Container.String():
		return Container, nil
	case Cve.String():
		return Cve, nil
	case Filesystem.String():
		return Filesystem, nil
	case Local.String():
		return Local, nil
	case Network.String():
		return Network, nil
	case Physical.String():
		return Physical, nil
	case Process.String():
		return Process, nil
	case Users.String():
		return Users, nil
	}

	return "unknown", fmt.Errorf("the input %q is not a category", input)
}

func (c *Category) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("couldn't unmarshal category data %q", string(data))
	}

	o, e := New(str)
	if e != nil {
		return e
	}
	*c = o

	return nil
}

func (c Category) MarshalJSON() ([]byte, error) {
	str := c.String()
	if _, err := New(str); err != nil {
		return nil, err
	}
	lower := strings.ToLower(str)

	return []byte(fmt.Sprintf("%q", lower)), nil
}
