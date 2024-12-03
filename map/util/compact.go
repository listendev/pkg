package maputil

import (
	stockmap "github.com/ghetzel/go-stockutil/maputil"
	typeutil "github.com/listendev/pkg/type/util"
)

func Compact(input map[string]interface{}) (map[string]interface{}, error) {
	output := make(map[string]interface{})

	if err := stockmap.Walk(input, func(value interface{}, path []string, isLeaf bool) error {
		if !typeutil.IsEmpty(value) {
			if typeutil.IsArray(value) {
				stockmap.DeepSet(output, path, value)

				return stockmap.SkipDescendants
			} else if isLeaf {
				stockmap.DeepSet(output, path, value)
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return output, nil
}
