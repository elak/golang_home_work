package hw02_unpack_string //nolint:golint,stylecheck

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type test struct {
	input    string
	expected string
	err      error
}

func TestUnpack(t *testing.T) {
	for _, tst := range [...]test{
		{
			input:    "a4bc2d5e",
			expected: "aaaabccddddde",
		},
		{
			input:    "abccd",
			expected: "abccd",
		},
		{
			input:    "3abc",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "45",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "aaa10b",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "",
			expected: "",
		},
		{
			input:    "aaa0b",
			expected: "aab",
		},
		{
			input:    "🐈3",
			expected: "🐈🐈🐈",
		},
		{
			input:    "🐈2🦉",
			expected: "🐈🐈🦉",
		},
		{
			input:    "\u65e5本5\U00008a9e",
			expected: "\u65e5本本本本本\U00008a9e",
		},
		// Потребуется нормализация уникода, и это будет уже не распаковка, а совсем наоборот
		// {
		// 	input:    "\U00000438\U000003062",
		// 	expected: "йй",
		// },
	} {
		result, err := Unpack(tst.input)
		require.Equal(t, tst.err, err, fmt.Sprintf("Error unpacking '%s'", tst.input))
		require.Equal(t, tst.expected, result, fmt.Sprintf("Error unpacking '%s'", tst.input))
	}
}

func TestUnpackWithEscape(t *testing.T) {
	for _, tst := range [...]test{
		{
			input:    `qwe\4\5`,
			expected: `qwe45`,
		},
		{
			input:    `qwe\45`,
			expected: `qwe44444`,
		},
		{
			input:    `qwe\\5`,
			expected: `qwe\\\\\`,
		},
		{
			input:    `qwe\\\3`,
			expected: `qwe\3`,
		},
		{
			input:    `\qwe`,
			expected: `qwe`,
		},
		{
			input:    `\`,
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    `qwe\\\3\`,
			expected: "",
			err:      ErrInvalidString,
		},
	} {
		result, err := Unpack(tst.input)
		require.Equal(t, tst.err, err, fmt.Sprintf("Error unpacking '%s'", tst.input))
		require.Equal(t, tst.expected, result, fmt.Sprintf("Error unpacking '%s'", tst.input))
	}
}
