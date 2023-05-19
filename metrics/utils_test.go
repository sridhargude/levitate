package metrics

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestUnmarshalJSON(t *testing.T) {
	// Test case 1: Valid JSON data
	rawJsonData := `{"metric":{"__name__":"up","job":"node_exporter","instance":"localhost:9100"},"values":[0]}`
	expectedData := Data{
		Metric: MetricData{
			Name: "up",
			Data: map[string]interface{}{
				"__name__": "up",
				"job":      "node_exporter",
				"instance": "localhost:9100",
			},
		},
		Values: []int{0},
	}
	testUnmarshalJSON(t, rawJsonData, expectedData)

	// Test case 2: Missing "__name__" field
	rawJsonData = `{"metric":{"job":"node_exporter","instance":"localhost:9100"},"values":[0]}`
	expectedData = Data{}
	testUnmarshalJSON(t, rawJsonData, expectedData)

	// Test case 3: Invalid "__name__" field type
	rawJsonData = `{"metric":{"__name__":123,"job":"node_exporter","instance":"localhost:9100"},"values":[0]}`
	expectedData = Data{}
	testUnmarshalJSON(t, rawJsonData, expectedData)

	// Test case 4: Empty "values" field
	rawJsonData = `{"metric":{"__name__":"up","job":"node_exporter","instance":"localhost:9100"},"values":[]}`
	expectedData = Data{
		Metric: MetricData{
			Name: "up",
			Data: map[string]interface{}{
				"__name__": "up",
				"job":      "node_exporter",
				"instance": "localhost:9100",
			},
		},
		Values: []int{},
	}
	testUnmarshalJSON(t, rawJsonData, expectedData)
}

// Helper function to test json.Unmarshal and compare the result
func testUnmarshalJSON(t *testing.T, rawJsonData string, expectedData Data) {
	var actualData Data
	err := json.Unmarshal([]byte(rawJsonData), &actualData)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Compare the expected and actual data
	if !compareData(expectedData, actualData) {
		t.Errorf("Unmarshaled data does not match the expected data")
	}
}

// Helper function to compare two Data structs
func compareData(expected Data, actual Data) bool {
	// Compare MetricData
	if expected.Metric.Name != actual.Metric.Name {
		return false
	}

	// Compare Values
	if len(expected.Values) != len(actual.Values) {
		return false
	}
	for i := range expected.Values {
		if expected.Values[i] != actual.Values[i] {
			return false
		}
	}

	// Compare MetricData's Data field
	for key, value := range expected.Metric.Data {
		if actual.Metric.Data[key] != value {
			return false
		}
	}

	return true
}

func TestTextRequest(t *testing.T) {
	textData := `up{job="node_exporter",instance="localhost:9100",k=1} 0`
	// Process text data
	textDataObj, err := ProcessDataText(textData)
	if err != nil {
		fmt.Println("Error processing text data:", err)
	} else {
		fmt.Println("Text data:")
		fmt.Println("Metric Name:", textDataObj.Metric.Name)
		fmt.Println("Metric Data:", textDataObj.Metric.Data)
		fmt.Println("Values:", textDataObj.Values)
	}
}
