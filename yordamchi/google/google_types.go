package google

type Request struct {
	Contents          []Content `json:"contents,omitempty"`
	SystemInstruction Content   `json:"systemInstruction,omitempty"`
}

type Content struct {
	Parts []Part `json:"parts,omitempty"`
	Role  string `json:"role,omitempty"`
}

type Part struct {
	Text string `json:"text,omitempty"`
}

type Response struct {
	Candidates    []Candidate   `json:"candidates"`
	UsageMetadata UsageMetadata `json:"usageMetadata"`
	ModelVersion  string        `json:"modelVersion"`
}

type Candidate struct {
	Content      Content `json:"content"`
	FinishReason string  `json:"finishReason"`
	Index        int     `json:"index"`
}

type UsageMetadata struct {
	PromptTokenCount     int            `json:"promptTokenCount"`
	CandidatesTokenCount int            `json:"candidatesTokenCount"`
	TotalTokenCount      int            `json:"totalTokenCount"`
	PromptTokensDetails  []TokensDetail `json:"promptTokensDetails"`
	ThoughtsTokenCount   int            `json:"thoughtsTokenCount"`
}

type TokensDetail struct {
	Modality   string `json:"modality"`
	TokenCount int    `json:"tokenCount"`
}
