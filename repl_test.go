package main

import ("testing")

func TestCleanInput(t *testing.T) {

cases := []struct {
		input string
		expected []string
	}{
		{
			input:   " hello world ",
			expected: []string{"hello", "world"},
		},
		{
			input: "charmander    charizard Mew",
			expected: []string{"charmander", "charizard", "mew"},
		},
		{
			input: "aSh Brock mistY",
			expected: []string{"ash", "brock", "misty"},
		},
	}


	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("test failed, not expected number of items")
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("test failed, not correct word")
			}
		}
	}
}