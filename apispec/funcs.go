package apispec

import "strings"

type SettingsOptioner func(*SettingsOptions)

type SettingsOptions struct {
	prefix    string
	separator string
}

// WithPrefix is an option to set a prefix
func WithPrefix(prefix string) SettingsOptioner {
	return func(opts *SettingsOptions) {
		opts.prefix = prefix
	}
}

// WithPrefix is an option to set a separator
func WithSeparator(separator string) SettingsOptioner {
	return func(opts *SettingsOptions) {
		opts.separator = separator
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

	if s.Tokens != nil {
		for tokenName, tokenData := range *s.Tokens {
			if tokenData.Key != nil {
				// Create the key in uppercase format
				prefix := opts.prefix
				if prefix != "" {
					prefix += opts.separator
				}
				key := strings.ToUpper(prefix+tokenName+opts.separator) + "TOKEN"
				result[key] = *tokenData.Key
			}
		}
	}

	return result
}
