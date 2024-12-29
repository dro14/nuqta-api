package elasticsearch

import (
	"fmt"
	"strings"
)

func (e *Elasticsearch) CreateIndex(index string) error {
	res, err := e.client.Indices.Create(
		index,
		e.client.Indices.Create.WithBody(strings.NewReader(`{
			"mappings": {
				"properties": {
					"username": {
						"type": "text",
						"analyzer": "standard",
						"fields": {
							"keyword": {
								"type": "keyword",
								"ignore_above": 256
							}
						}
					}
				}
			}
		}`)),
	)
	if err != nil {
		return err
	}
	if res.IsError() {
		return fmt.Errorf("error creating index: %s", res.String())
	}
	return nil
}
