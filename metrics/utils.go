package metrics

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// UnmarshalJSON converts JSON data into Data struct in metrics
func (m *MetricData) UnmarshalJSON(data []byte) error {
	// Unmarshal metric data into a temporary map
	tmp := make(map[string]interface{})
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	// Extract the "__name__" field and store the rest in Data
	name, ok := tmp["__name__"].(string)
	if !ok {
		return fmt.Errorf("__name__ field is missing or not a string")
	}

	m.Name = name
	// Add the __name__ to  labels
	// delete(tmp, "__name__")
	m.Data = tmp

	return nil
}

// ProcessDataText converts OpenMetrics data from text format into Data struct in metrics
func ProcessDataText(textData string) (*Data, error) {
	dataObj := &Data{
		Values: make([]int, 1),
	}

	// Extract metric name and labels
	re := regexp.MustCompile(`^(.*?)\{(.*?)\} (.*)$`)
	match := re.FindStringSubmatch(textData)
	if match == nil {
		return nil, fmt.Errorf("invalid text data format")
	}

	dataObj.Metric.Name = match[1]
	dataObj.Metric.Data = make(map[string]interface{})

	// Add Metric Name
	dataObj.Metric.Data["__name__"] = dataObj.Metric.Name

	labelPairs := strings.Split(match[2], ",")

	// Add the rest of the labels
	for _, pair := range labelPairs {
		kv := strings.Split(pair, "=")
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		value := strings.Trim(kv[1], `"`)

		dataObj.Metric.Data[key] = value
	}

	// Extract value
	_, err := fmt.Sscanf(match[3], "%d", &dataObj.Values[0])
	if err != nil {
		return nil, fmt.Errorf("invalid value in text data: %v", err)
	}

	return dataObj, nil
}

func StructToMap(obj interface{}) (map[string]interface{}, error) {
	objValue := reflect.ValueOf(obj)
	if objValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input is not a struct")
	}

	data := make(map[string]interface{})
	objType := objValue.Type()
	for i := 0; i < objValue.NumField(); i++ {
		field := objType.Field(i)
		fieldValue := objValue.Field(i).Interface()
		data[field.Name] = fieldValue
	}

	return data, nil
}
