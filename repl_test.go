package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    " Charmander Bulbasaur PIKACHU ",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input:    " ChaRmAnder BulBaSAur PikAchu ",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input:    "  charman  pikapoo  ",
			expected: []string{"charman", "pikapoo"},
		},
		// add more cases here
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		// Check the length of the actual slice against the expected slice
		// if they don't match, use t.Errorf to print an error message
		// and fail the test
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			// Check each word in the slice
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
			if word != expectedWord {
				t.Errorf("For input %q, at index %d expected %q, got %q", c.input, i, expectedWord, word)
			}
		}
	}
}
