package apispec

import (
	"fmt"
	"sort"
	"strings"
)

type SettingsOptioner func(*SettingsOptions)

type SettingsOptions struct {
	prefix      string
	separator   string
	dquoteValue bool
}

// WithPrefix is an option to set a prefix.
func WithPrefix(prefix string) SettingsOptioner {
	return func(opts *SettingsOptions) {
		opts.prefix = prefix
	}
}

// WithPrefix is an option to set a separator.
func WithSeparator(separator string) SettingsOptioner {
	return func(opts *SettingsOptions) {
		opts.separator = separator
	}
}

func WithValueDoubleQuotes() SettingsOptioner {
	return func(opts *SettingsOptions) {
		opts.dquoteValue = true
	}
}

func (s *Settings) TokensMap(options ...SettingsOptioner) map[string]string {
	// Apply options
	opts := &SettingsOptions{
		separator: "_",
	}
	for _, opt := range options {
		opt(opts)
	}

	result := make(map[string]string)
	suffix := "token"

	if s.Tokens != nil {
		for tokenName, tokenData := range *s.Tokens {
			if tokenData.Key != nil {
				prefix := opts.prefix
				if prefix != "" {
					prefix += opts.separator
				}
				// Create the key in uppercase format
				key := strings.ToUpper(prefix + tokenName + opts.separator + suffix)
				val := *tokenData.Key
				if opts.dquoteValue {
					val = fmt.Sprintf("%q", val)
				}
				result[key] = val
			}
		}
	}

	return result
}

func (s *Settings) TokensSlice(options ...SettingsOptioner) []string {
	tokensMap := s.TokensMap(append([]SettingsOptioner{WithValueDoubleQuotes()}, options...)...)

	res := []string{}
	for k, v := range tokensMap {
		res = append(res, fmt.Sprintf("%s=%s", k, v))
	}
	sort.Strings(res)

	return res
}
