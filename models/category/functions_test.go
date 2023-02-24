package category

import (
	"testing"
)

func TestCategories(t *testing.T) {
	cats := Categories()
	for _, c := range cats {
		if c == 0 {
			t.Errorf("categories() returned an empty category")
		}
	}
	if len(cats) != len(all) {
		t.Errorf("categories() returned %d categories, expected %d", len(cats), len(all))
	}
}
