package models

// assert model
type Assert struct {
	Status   []AssertItem `json:"status"`
	JSONPath []AssertItem `json:"jsonpath"`
	Cookie   []AssertItem `json:"cookie"`
	Header   []AssertItem `json:"header"`
	Body     []AssertItem `json:"body"`
}

// assert item model
type AssertItem struct {
	Type      string `json:"type"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	ValueType string `json:"value_type"`
}
