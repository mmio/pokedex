package main

import (
	"fmt"
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input string
		expected []string
	} {
		{
			input: " hello world ",
				expected: []string{"hello", "world"},
			},
			{
			input: "",
				expected: []string{},
			},
			{
			input: " ",
				expected: []string{},
			},
			{
			input: fmt.Sprintf("\n"),
				expected: []string{},			
			},
			{
			input: "apwoeirjpaoweifjpawoiefjpoaiwefjpaiwefjpawifjepoiawjfepoijawepfoijawoefijapoweijpowaiefjpoawifj",
				expected: []string{"apwoeirjpaoweifjpawoiefjpoaiwefjpaiwefjpawifjepoiawjfepoijawepfoijawoefijapoweijpowaiefjpoawifj"},
			},
			{
			input: "a b c d",
				expected: []string{"a", "b", "c", "d"},
			},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("actual expected don't have same len: %v %v", len(actual), len(c.expected))
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("actualWord doesn't match expectedWord: %v %v at index %v", word, expectedWord, i)
			}
		}
	}
	
}
