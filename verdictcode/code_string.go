// Code generated by "stringer -type=Code"; DO NOT EDIT.

package verdictcode

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[DDN01-1001]
	_ = x[FNI001-1]
	_ = x[FNI002-2]
	_ = x[FNI003-3]
	_ = x[MDN01-1021]
	_ = x[MDN02-1022]
	_ = x[MDN03-1023]
	_ = x[MDN04-1024]
	_ = x[STN001-1101]
	_ = x[STN002-1102]
	_ = x[STN003-1103]
	_ = x[STN004-1104]
	_ = x[STN005-1105]
	_ = x[STN006-1106]
	_ = x[STN007-1107]
	_ = x[STN008-1108]
	_ = x[TSN01-1011]
	_ = x[UNK-0]
}

const (
	_Code_name_0 = "UNKFNI001FNI002FNI003"
	_Code_name_1 = "DDN01"
	_Code_name_2 = "TSN01"
	_Code_name_3 = "MDN01MDN02MDN03MDN04"
	_Code_name_4 = "STN001STN002STN003STN004STN005STN006STN007STN008"
)

var (
	_Code_index_0 = [...]uint8{0, 3, 9, 15, 21}
	_Code_index_3 = [...]uint8{0, 5, 10, 15, 20}
	_Code_index_4 = [...]uint8{0, 6, 12, 18, 24, 30, 36, 42, 48}
)

func (i Code) String() string {
	switch {
	case i <= 3:
		return _Code_name_0[_Code_index_0[i]:_Code_index_0[i+1]]
	case i == 1001:
		return _Code_name_1
	case i == 1011:
		return _Code_name_2
	case 1021 <= i && i <= 1024:
		i -= 1021
		return _Code_name_3[_Code_index_3[i]:_Code_index_3[i+1]]
	case 1101 <= i && i <= 1108:
		i -= 1101
		return _Code_name_4[_Code_index_4[i]:_Code_index_4[i+1]]
	default:
		return "Code(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
