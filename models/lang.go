package models

import (
	"github.com/PaesslerAG/gval"
	"github.com/PaesslerAG/jsonpath"
)

var lang = gval.Full(jsonpath.Language())
