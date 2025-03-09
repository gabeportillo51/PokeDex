package main
import "testing"

func TestCleanInput(t *testing.T) {
	// This is a tester function to ensure that the cleanInput function works as intended
	test_cases := []struct {
		input string
		expected []string
	}{
		{	// testing capital letters
			input: "HeY hOw'S iT gOiNg, mAn?",
			expected: []string{"hey", "how's", "it", "going,", "man?"},
		},
		{	// testing capital letters + trailing/leading white space
			input: "    WhAt tHe Fuck is GoinG on MaN?   ",
			expected: []string{"what", "the", "fuck", "is", "going", "on", "man?"},
		},
		{	// testing no spaces
			input: "heyhowsitgoingbrother?",
			expected: []string{"heyhowsitgoingbrother?"},
		},
	}

	for _, c := range test_cases {
		actual_result := cleanInput(c.input)
		if len(actual_result) != len(c.expected) {
			t.Errorf("Error: incompatible lengths")
		}
		for index := range actual_result {
			word := actual_result[index]
			expected_word := c.expected[index]
			if word != expected_word {
				t.Errorf("Error: mismatched words - %s and %s", word, expected_word)
			}
		}
	}
}
