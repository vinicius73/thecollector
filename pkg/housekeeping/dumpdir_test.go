package housekeeping_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vinicius73/thecollector/pkg/housekeeping"
)

func TestIsDumpDir(t *testing.T) {
	testCases := []struct {
		dir      string
		expected bool
	}{
		{
			dir:      "2021/01/01/database",
			expected: true,
		},
		{
			dir:      "2021/01/01/database/",
			expected: false,
		},
		{
			dir:      "2021/01/database/other",
			expected: false,
		},
		{
			dir:      "FOOO/01/01/database/other/",
			expected: false,
		},
		{
			dir:      "2021/OO/01/other",
			expected: false,
		},
		{
			dir:      "2021/01/gol",
			expected: false,
		},
		{
			dir:      "2021/01/01",
			expected: false,
		},
		{
			dir:      "2021/01/01/",
			expected: false,
		},
		{
			dir:      "2021/01/AA/database",
			expected: false,
		},
	}

	for _, tc := range testCases {
		actual := tc
		t.Run(actual.dir, func(t *testing.T) {
			got := housekeeping.IsDumpDir(tc.dir)
			assert.Equalf(t, actual.expected, got, "IsDumpDir(%s) = %v; want %v", actual.dir, got, actual.expected)
		})
	}
}
