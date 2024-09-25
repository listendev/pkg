package apispec

import (
	"testing"
)

func TestTokensMap(t *testing.T) {
	// Prepare test data
	key1 := "ABC"
	key2 := "ZZZ"
	tokens := map[string]struct {
		Key *string `json:"key,omitempty"`
	}{
		"openai": {Key: &key1},
		"claude": {Key: &key2},
	}

	settings := Settings{
		Tokens: &tokens,
	}

	// Expected result
	expected := map[string]string{
		"OPENAI_TOKEN": "ABC",
		"CLAUDE_TOKEN": "ZZZ",
	}

	// Call the method
	result := settings.TokensMap()

	// Check if the result matches the expected output
	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for key, expectedValue := range expected {
		if resultValue, exists := result[key]; !exists || resultValue != expectedValue {
			t.Errorf("For key %s, expected value %s, got %s", key, expectedValue, resultValue)
		}
	}
}

func TestTokensSlice(t *testing.T) {
	// Prepare test data
	key1 := "ABC"
	key2 := "ZZZ"
	tokens := map[string]struct {
		Key *string `json:"key,omitempty"`
	}{
		"openai": {Key: &key1},
		"claude": {Key: &key2},
	}

	settings := Settings{
		Tokens: &tokens,
	}

	// Expected result
	expected := []string{
		"CLAUDE_TOKEN=ZZZ",
		"OPENAI_TOKEN=ABC",
	}

	// Call the method
	result := settings.TokensSlice()

	// Check if the result matches the expected output
	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for key, expectedValue := range expected {
		if resultValue := result[key]; resultValue != expectedValue {
			t.Errorf("Expected value %s, got %s", expectedValue, resultValue)
		}
	}
}

func TestTokensSliceWithValueDoubleQuotes(t *testing.T) {
	// Prepare test data
	key1 := "ABC"
	key2 := "ZZZ"
	tokens := map[string]struct {
		Key *string `json:"key,omitempty"`
	}{
		"openai": {Key: &key1},
		"claude": {Key: &key2},
	}

	settings := Settings{
		Tokens: &tokens,
	}

	// Expected result
	expected := []string{
		"CLAUDE_TOKEN=\"ZZZ\"",
		"OPENAI_TOKEN=\"ABC\"",
	}

	// Call the method
	result := settings.TokensSlice(WithValueDoubleQuotes())

	// Check if the result matches the expected output
	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for key, expectedValue := range expected {
		if resultValue := result[key]; resultValue != expectedValue {
			t.Errorf("Expected value %s, got %s", expectedValue, resultValue)
		}
	}
}

func TestTokensMapWithQuote(t *testing.T) {
	// Prepare test data
	key1 := "ABC"
	key2 := "ZZZ"
	tokens := map[string]struct {
		Key *string `json:"key,omitempty"`
	}{
		"openai": {Key: &key1},
		"claude": {Key: &key2},
	}

	settings := Settings{
		Tokens: &tokens,
	}

	// Expected result
	expected := map[string]string{
		"OPENAI_TOKEN": "\"ABC\"",
		"CLAUDE_TOKEN": "\"ZZZ\"",
	}

	// Call the method
	result := settings.TokensMap(WithValueDoubleQuotes())

	// Check if the result matches the expected output
	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for key, expectedValue := range expected {
		if resultValue, exists := result[key]; !exists || resultValue != expectedValue {
			t.Errorf("For key %s, expected value %s, got %s", key, expectedValue, resultValue)
		}
	}
}

func TestTokensMapWithPrefix(t *testing.T) {
	// Prepare test data
	key1 := "ABC"
	key2 := "ZZZ"
	tokens := map[string]struct {
		Key *string `json:"key,omitempty"`
	}{
		"openai": {Key: &key1},
		"claude": {Key: &key2},
	}

	settings := Settings{
		Tokens: &tokens,
	}

	// Expected result
	expected := map[string]string{
		"LSTN_OPENAI_TOKEN": "ABC",
		"LSTN_CLAUDE_TOKEN": "ZZZ",
	}

	// Call the method
	result := settings.TokensMap(WithPrefix("LSTN"))

	// Check if the result matches the expected output
	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for key, expectedValue := range expected {
		if resultValue, exists := result[key]; !exists || resultValue != expectedValue {
			t.Errorf("For key %s, expected value %s, got %s", key, expectedValue, resultValue)
		}
	}
}
