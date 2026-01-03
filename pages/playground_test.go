package pages

import (
	"testing"

	"github.com/tidwall/gjson"
)

// Basic logic test to ensure GJSON is integrated correctly and logic holds
func TestGJSONLogic(t *testing.T) {
	jsonContent := `{
  "name": {"first": "Tom", "last": "Anderson"},
  "age": 37,
  "children": ["Sara","Alex","Jack"]
}`
	query := "name.last"

	res := gjson.Get(jsonContent, query)
	if res.String() != "Anderson" {
		t.Errorf("Expected Anderson, got %s", res.String())
	}
}

// Test Component Logic (mocked)
func TestUpdateResult(t *testing.T) {
	comp := &GJSONPlayground{}
	comp.JSONContent = `{"key": "value"}`
	comp.Query = "key"

	comp.calculateResult()

	if comp.Result != "value" {
		t.Errorf("Expected 'value', got '%s'", comp.Result)
	}
}

// Test Empty Query
func TestEmptyQuery(t *testing.T) {
	comp := &GJSONPlayground{}
	comp.JSONContent = `{"key": "value"}`
	comp.Query = ""

	comp.calculateResult()

	if comp.Result != "" {
		t.Errorf("Expected empty string, got '%s'", comp.Result)
	}
}

// Test Invalid Path (Result not found)
func TestInvalidPath(t *testing.T) {
	comp := &GJSONPlayground{}
	comp.JSONContent = `{"key": "value"}`
	comp.Query = "non.existent.path"

	comp.calculateResult()

	if comp.Result != "" {
		t.Errorf("Expected empty string for non-existent path, got '%s'", comp.Result)
	}
	if comp.PathFound {
		t.Errorf("Expected PathFound to be false")
	}
}
