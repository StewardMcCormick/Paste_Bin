package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateAPIKey(t *testing.T) {
	cases := []struct {
		name     string
		value    string
		expected bool
	}{
		{
			"Valid key",
			"pb_live_abcd_abcdefghijkl",
			true,
		},
		{
			"Key with not enough fields",
			"pb_live_",
			false,
		},
		{
			"Key with too short first_random field",
			"pb_live_abc_abcdefghijkl",
			false,
		},
		{
			"Key with too short second_random field",
			"pb_live_abcd_abcdef",
			false,
		},
		{
			"Key with too short all fields",
			"pb_live_abc_abcdefg",
			false,
		},
		{
			"Key with invalid prefix",
			"pb_invalid_abcd_abcdefghijkl",
			false,
		},
		{
			"Empty key",
			"",
			false,
		},
		{
			"Invalid key",
			"key",
			false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := ValidateAPIKey(tc.value)

			assert.Equal(t, tc.expected, res)
		})
	}
}
