package metrics

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestJSON(t *testing.T) {
	jsonData := `{"metric":{"__name__":"up","job":"node_exporter","instance":"localhost:9100"},"values":[0]}`

	// Process JSON data
	var jsonDataObj Data
	err := json.Unmarshal([]byte(jsonData), &jsonDataObj)
	if err != nil {
		fmt.Println("Error processing JSON data:", err)
	} else {
		fmt.Println("JSON data:")
		fmt.Println("Metric Name:", jsonDataObj.Metric.Name)
		fmt.Println("Metric Data:", jsonDataObj.Metric.Data)
		fmt.Println("Values:", jsonDataObj.Values)
	}

}

func TestText(t *testing.T) {
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
