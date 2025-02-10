package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world   ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input:    "    ChIkoRitA BaYLEEF MEGanium   ",
			expected: []string{"chikorita", "bayleef", "meganium"},
		},
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "       ",
			expected: []string{},
		},
		{
			input:    "     cHaRIZARD",
			expected: []string{"charizard"},
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("handles input: %s", c.input), func(t *testing.T) {
			actual := cleanInput(c.input)

			assert.Equal(t, c.expected, actual)
		})
	}
}
