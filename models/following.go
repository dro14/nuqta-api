package models

type Following struct {
	Uid     string  `json:"uid"`
	Posts   []*Post `json:"posts"`
	Reposts []*Post `json:"reposts"`
}
