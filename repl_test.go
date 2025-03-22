package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "   ",
			expected: []string{},
		},
		{
			input:    "one",
			expected: []string{"one"},
		},
		{
			input:    "  one   two   three  ",
			expected: []string{"one", "two", "three"},
		},
		{
			input:    "\t tab\tseparated\twords\t",
			expected: []string{"tab", "separated", "words"},
		},
		{
			input:    "line\nbreaks\nincluded",
			expected: []string{"line", "breaks", "included"},
		},
		{
			input:    "mix of \t spaces \n and lines",
			expected: []string{"mix", "of", "spaces", "and", "lines"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("\nExpected: %v\nActual: %v", expectedWord, word)
			}
		}
	}
}
