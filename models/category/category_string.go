// Code generated by "stringer -type=Category"; DO NOT EDIT.

package category

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[AdjacentNetwork-8]
	_ = x[Advisory-7]
	_ = x[CIS-6]
	_ = x[Container-5]
	_ = x[Cybersquatting-11]
	_ = x[Filesystem-1]
	_ = x[Local-9]
	_ = x[Network-3]
	_ = x[Physical-10]
	_ = x[Process-2]
	_ = x[Users-4]
}

const _Category_name = "FilesystemProcessNetworkUsersContainerCISAdvisoryAdjacentNetworkLocalPhysicalCybersquatting"

var _Category_index = [...]uint8{0, 10, 17, 24, 29, 38, 41, 49, 64, 69, 77, 91}

func (i Category) String() string {
	i -= 1
	if i >= Category(len(_Category_index)-1) {
		return "Category(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _Category_name[_Category_index[i]:_Category_index[i+1]]
}
