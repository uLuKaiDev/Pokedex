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
			input:    "HELLO World",
			expected: []string{"hello", "world"},
		},
		{
			input:    "LET'S TEST some more",
			expected: []string{"let's", "test", "some", "more"},
		},
		{
			input:    "1234      5678",
			expected: []string{"1234", "5678"},
		},
	}

	for _, c := range cases {
		actual := CleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Actual length(%v) does not match expected length(%v)", len(actual), len(c.expected))
			return
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("Words do not match. Actual: %v, Expected: %v", actual[i], c.expected[i])
				return
			}
		}
	}
}
