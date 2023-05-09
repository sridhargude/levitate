package metrics

type Data struct {
	Metric MetricData `json:"metric"`
	Values []int      `json:"values"`
}

type MetricData struct {
	Name string                 `json:"__name__"`
	Data map[string]interface{} `json:"-"`
}
