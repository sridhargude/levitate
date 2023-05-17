package hll

import (
	"encoding/json"
	h "github.com/mtchavez/go-hll/hll"
	"github.com/sridhargude/levitate/metrics"
)

var HLLDone chan bool
var HLLChan chan metrics.Data

type MetricsHLL struct {
	hll          *h.Hll
	Count        int64
	TotalMetrics int64
}

func Init() *MetricsHLL {
	return &MetricsHLL{hll: h.New(), Count: 0, TotalMetrics: 0}
}

func (m *MetricsHLL) Add(dict map[string]interface{}) error {
	// Convert the dictionary to a JSON string
	dictBytes, err := json.Marshal(dict)
	if err != nil {
		return err
	}
	dictString := string(dictBytes)

	// Add the JSON string to HLL data struct
	m.hll.Add(dictString)
	// Update the Count
	m.Count = int64(m.hll.Count())
	m.TotalMetrics++
	return nil
}

func (m *MetricsHLL) CalculateDistinctMetrics(dictionaries []map[string]interface{}) (int64, error) {

	for _, dictionary := range dictionaries {
		m.Add(dictionary)
	}

	// Retrieve the estimated distinct count of the metrics dict
	m.Count = int64(m.hll.Count())
	return m.Count, nil
}
