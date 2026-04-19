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

func TestSaveAndLoad(t *testing.T) {
	cfg := &Config{
		SaveDir:  t.TempDir(),
		Pokedex:  make(map[string]Pokemon),
		PlayerLV: 5,
		PlayerXP: 120,
	}

	// save
	if err := saveGameState(cfg, nil); err != nil {
		t.Errorf("saveGameState function failed: %v", err)
	}

	// create a fresh empty config with same SaveDir

	cfg2 := &Config{
		SaveDir: cfg.SaveDir,
		Pokedex: make(map[string]Pokemon),
	}

	// load
	if err := loadGameState(cfg2, nil); err != nil {
		t.Errorf("loadGameState function failed: %v", err)
	}

	if cfg2.PlayerLV != 5 {
		t.Errorf("expected PlayerLV 5, got %d", cfg.PlayerLV)
	}
	if cfg2.PlayerXP != 120 {
		t.Errorf("expected PlayerXP 120, got %d", cfg.PlayerXP)
	}
}
