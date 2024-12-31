package elastic

type Doc struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	HitCount int    `json:"hit_count"`
}

type Request struct {
	Query Query `json:"query"`
}

type Query struct {
	FunctionScore *FunctionScore `json:"function_score,omitempty"`
	MultiMatch    *MultiMatch    `json:"multi_match,omitempty"`
}

type FunctionScore struct {
	Query       Query       `json:"query"`
	FieldFactor FieldFactor `json:"field_value_factor"`
	BoostMode   string      `json:"boost_mode"`
}

type FieldFactor struct {
	Field    string  `json:"field"`
	Factor   float64 `json:"factor"`
	Modifier string  `json:"modifier"`
	Missing  float64 `json:"missing"`
}

type MultiMatch struct {
	Query     string   `json:"query"`
	Fields    []string `json:"fields"`
	Fuzziness string   `json:"fuzziness"`
}

type Result struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Hits     Hits `json:"hits"`
}

type Hits struct {
	Total    Total   `json:"total"`
	MaxScore float64 `json:"max_score"`
	Hits     []Hit   `json:"hits"`
}

type Total struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

type Hit struct {
	Index  string  `json:"_index"`
	ID     string  `json:"_id"`
	Score  float64 `json:"_score"`
	Source Doc     `json:"_source"`
}
