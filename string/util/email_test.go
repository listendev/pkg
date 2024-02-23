package stringutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatchEmails(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			input: "Sample text with emails user1@example.com and user2@example.com",
			expected: []string{
				"user1@example.com",
				"user2@example.com",
			},
		},
		{
			input:    "No emails in this text",
			expected: []string(nil),
		},
		{
			input:    "Only one email: test@example.com",
			expected: []string{"test@example.com"},
		},
		{
			input: "Multiple emails separated by space: test1@example.com test2@example.com",
			expected: []string{
				"test1@example.com",
				"test2@example.com",
			},
		},
		{
			input:    "Email with numbers: test-123@example.co.uk",
			expected: []string{"test-123@example.co.uk"},
		},
		{
			input: "Emails with subdomains: user@mail.example.com, info@sub.domain.com",
			expected: []string{
				"user@mail.example.com",
				"info@sub.domain.com",
			},
		},
		{
			input: "Emails with numeric domain: user@123.com, info@1example123.co",
			expected: []string{
				"user@123.com",
				"info@1example123.co",
			},
		},
		{
			input: "Emails with international domain: user@тест.срб user@example.рф, idn@email.भारत info@xn--80akhbyknj4f.xn--90ais ñoñó1234@server.com",
			expected: []string{
				"user@тест.срб",
				"user@example.рф",
				"idn@email.भारत",
				"info@xn--80akhbyknj4f.xn--90ais",
				"ñoñó1234@server.com",
			},
		},
		{
			input: `#!$%&'*+-/=?^_` + `{}|~@example.org crazy: "()<>[]:,;@\\\"!#$%&'*+-/=?^_` + `{}| ~.a"@example.org`,
			expected: []string{
				"#!$%&'*+-/=?^_{}|~@example.org",
				`"()<>[]:,;@\\\"!#$%&'*+-/=?^_` + `{}| ~.a"@example.org`,
			},
		},
		{
			input: `prettyandsimple@example.com
			very.common@example.com
			disposable.style.email.with+symbol@example.com
			other.email-with-dash@example.com
			x@example.com (one-letter local part)
			"much.more unusual"@example.com
			"very.unusual.@.unusual.com"@example.com
			"very.(),:;<>[]\".VERY.\"very@\\ \"very\".unusual"@strange.example.com
			example-indeed@strange-example.com
			admin@mailserver1 (local domain name with no TLD)
			" "@example.org (space between the quotes)
			example@localhost (sent from localhost)
			example@s.solutions (see the List of Internet top-level domains)
			user@com
			user@localserver
			user@[IPv6:2001:db8::1]
			©other.email-with-dash@example.com
			?prettyandsimple@example.com

			Invalid email addresses[edit]
			Abc.example.com (no @ character)
			A@b@c@example.com (only one @ is allowed outside quotation marks)
			 a"b(c)d,e:f;g<h>i[j\k]l@example.com (none of the special characters in this local part are allowed outside quotation marks)
			just"not"right@example.com (quoted strings must be dot separated or the only element making up the local part)
			this is"not\allowed@example.com (spaces, quotes, and backslashes may only exist when within quoted strings and preceded by a backslash)
			this\ still\"not\\allowed@example.com (even if escaped (preceded by a backslash), spaces, quotes, and backslashes must still be contained by quotes)
			john..doe@example.com (double dot before @)
			with caveat: Gmail lets this through, Email address#Local-part the dots altogether
			john.doe@example..com (double dot after @)`,
			expected: []string{
				"prettyandsimple@example.com",
				"very.common@example.com",
				"disposable.style.email.with+symbol@example.com",
				"other.email-with-dash@example.com",
				"x@example.com",
				`"much.more unusual"@example.com`,
				`"very.unusual.@.unusual.com"@example.com`,
				`"very.(),:;<>[]\".VERY.\"very@\\ \"very\".unusual"@strange.example.com`,
				"example-indeed@strange-example.com",
				`" "@example.org`,
				`example@s.solutions`,
				"©other.email-with-dash@example.com",
				"?prettyandsimple@example.com",
				"c@example.com", // From this one onwards it is extracting emails that would be invalid if matching from the beginning of each line
				"l@example.com",
				"right@example.com",
				"not\\allowed@example.com",
				"not\\\\allowed@example.com",
				"doe@example.com",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := MatchEmails(test.input)
			require.Equal(t, test.expected, result)
		})
	}
}
