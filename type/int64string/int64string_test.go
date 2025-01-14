package int64string

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestInt64String_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    Int64String
		expected string
	}{
		{"Positive Value", 123, "\"123\""},
		{"Negative Value", -456, "\"-456\""},
		{"Zero Value", 0, "\"0\""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := json.Marshal(test.input)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, string(result))
		})
	}
}

func TestInt64String_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Int64String
	}{
		{"Positive Value", "\"123\"", 123},
		{"Negative Value", "\"-456\"", -456},
		{"Zero Value", "\"0\"", 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result Int64String
			err := json.Unmarshal([]byte(test.input), &result)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestInt64String_MarshalBSONValue(t *testing.T) {
	tests := []struct {
		name     string
		input    Int64String
		expected bson.RawValue
	}{
		{"Positive Value", 123, bson.RawValue{Type: bson.TypeString, Value: bson.Raw("123")}},
		{"Negative Value", -456, bson.RawValue{Type: bson.TypeString, Value: bson.Raw("-456")}},
		{"Zero Value", 0, bson.RawValue{Type: bson.TypeString, Value: bson.Raw("0")}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.input.MarshalBSONValue()
			assert.NoError(t, err)
			assert.Equal(t, test.expected, result)
		})
	}
}
