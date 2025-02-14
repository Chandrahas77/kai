package daos

import (
	"testing"
)

func TestParseRiskFactors_ValidJSON(t *testing.T) {
	jsonString := `["Remote Code Execution", "High CVSS Score"]`
	expected := []string{"Remote Code Execution", "High CVSS Score"}

	result := parseRiskFactors(jsonString)
	if len(result) != len(expected) || result[0] != expected[0] || result[1] != expected[1] {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestParseRiskFactors_InvalidJSON(t *testing.T) {
	jsonString := `invalid-json`

	result := parseRiskFactors(jsonString)
	if len(result) != 0 {
		t.Errorf("expected empty array on invalid JSON, got %v", result)
	}
}
