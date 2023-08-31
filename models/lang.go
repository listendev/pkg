package models

import (
	"github.com/PaesslerAG/gval"
	"github.com/PaesslerAG/jsonpath"
)

// lang is used to evaluate JSONPath expressions.
//
// lang is the union of the following languages: base, arithmetic, bitmask, text, propositional logic, ljson languages, and jsonpath prefix extensions. Plus the following operators:
//
//	in: a in b is true iff value a is an element of array b.
//	??: a ?? b returns a if a is not false or nil, otherwise n (colaesce operator).
//	?: a ? b : c returns b if bool a is true, otherwise b.
//
// And the following functions:
//
//	date(a): it parses string `a``. Notice `a`` must match RFC3339, ISO8601, ruby date, or unix date.
//
// The base language contains equal (==) and not equal (!=), perentheses and general support for variables, constants and functions. It contains true and false, (floating point) number, string  ("" or “) and char (”) constants.
//
// The arithmetic language contains base, plus(+), minus(-), divide(/), power(**), negative(-) and numerical order (<=,<,>,>=).
//
// The bitmask language contains base, bitwise and(&), bitwise or(|) and bitwise not(^).
//
// The text language contains base, lexical order on strings (<=,<,>,>=), regex match (=~) and regex not match (!~).
//
// The propositional logic language contains not(!), and (&&), or (||).
//
// The ljson language contains json objects ({string:expression,...}) and json arrays ([expression, ...]).
//
// The jsonpath extension language adds support for the $ and @ prefix extensions as per JSONPath grammar.
var lang = gval.Full(jsonpath.Language())
