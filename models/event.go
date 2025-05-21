package models

type Event struct {
	Name string `json:"name"`
	Data any    `json:"data"`
}
