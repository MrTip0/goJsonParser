package jsonParser

import "testing"

func TestEmptyArray(t *testing.T) {
	val, err := ParseJson([]rune(`[]`))
	if err != nil {
		t.Errorf("an error occurred: %v", err)
	}
	if len(val.([]any)) != 0 {
		t.Error("should be empty")
	}
}

func TestEmptyObject(t *testing.T) {
	val, err := ParseJson([]rune(`{}`))
	if err != nil {
		t.Errorf("an error occurred: %v", err)
	}
	if len(val.(map[string]any)) != 0 {
		t.Error("should be empty")
	}
}

func TestParsable1(t *testing.T) {
	val, err := ParseJson([]rune(`
	{
		"greetings": [
			"ciao",
			"hello",
			"こんにちは"
		],
		"number": 12.549,
		"numbers": [
			1,
			2,
			3,
			4,
			5
		]
	}`))
	if err != nil {
		t.Errorf("an error occurred: %v", err)
	}
	expected := map[string]any{
		"greetings": []any{"ciao", "hello", "こんにちは"},
		"number":    12.549,
		"numbers":   []any{int64(1), int64(2), int64(3), int64(4), int64(5)},
	}
	casted := val.(map[string]any)
	compareMaps(expected, casted, t)
}

func TestParsable2(t *testing.T) {
	val, err := ParseJson([]rune(`
	[
		"ciao",
		"hello",
		"こんにちは",
		"number",
		12.0549,
		{
			"propriety": "value",
			"boolean": true,
			"boolean2": false,
			"nullval": null
		}
	]`))
	if err != nil {
		t.Errorf("an error occurred: %v", err)
	}
	expected := []any{
		"ciao",
		"hello",
		"こんにちは",
		"number",
		12.0549,
		map[string]any{
			"propriety": "value",
			"boolean":   true,
			"boolean2":  false,
			"nullval":   nil,
		},
	}
	casted := val.([]any)
	compareArrays(expected, casted, t)
}

func TestError1(t *testing.T) {
	_, err := ParseJson([]rune(`[invalid_value]`))
	if err == nil {
		t.Error("should not parse")
	}
}

func TestError2(t *testing.T) {
	_, err := ParseJson([]rune(`{invalid_name}`))
	if err == nil {
		t.Error("should not parse")
	}
}

func TestError3(t *testing.T) {
	_, err := ParseJson([]rune(`["valid_name": "but this is not an object"]`))
	if err == nil {
		t.Error("should not parse")
	}
}

func TestError4(t *testing.T) {
	_, err := ParseJson([]rune(`[pietra]`))
	if err == nil {
		t.Error("should not parse")
	}
}

func compareMaps(expected map[string]any, value map[string]any, t *testing.T) {
	for k, expectedVal := range expected {
		val, ok := value[k]
		if !ok {
			t.Errorf("attribute not parsed: %s", k)
		}
		switch expectedVal.(type) {
		case map[string]any:
			compareMaps(expectedVal.(map[string]any), val.(map[string]any), t)
		case []any:
			compareArrays(expectedVal.([]any), val.([]any), t)
		default:
			if val != expectedVal {
				t.Errorf("different value, expected: %v, got: %v",
					expectedVal, val)
			}
		}
	}
}

func compareArrays(expected []any, value []any, t *testing.T) {
	if len(expected) != len(value) {
		t.Errorf("different arrays, expected: %v, got: %v", expected, value)
	}
	for i, expectedVal := range expected {
		val := value[i]
		switch expectedVal.(type) {
		case map[string]any:
			compareMaps(expectedVal.(map[string]any), val.(map[string]any), t)
		case []any:
			compareArrays(expectedVal.([]any), val.([]any), t)
		default:
			if val != expectedVal {
				t.Errorf("different value, expected: %v, got: %v",
					expectedVal, val)
			}
		}
	}
}
