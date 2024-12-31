package elastic

import (
	"context"
	"encoding/json"
	"net/http"
)

func id(ctx context.Context) string {
	return ctx.Value("id").(string)
}

func searchRequest(query string) Request {
	return Request{
		Query: Query{
			FunctionScore: &FunctionScore{
				Query: Query{
					MultiMatch: &MultiMatch{
						Query:     query,
						Fields:    []string{"username", "name"},
						Fuzziness: "AUTO",
					},
				},
				FieldFactor: FieldFactor{
					Field:    "hit_count",
					Factor:   1,
					Modifier: "log1p",
					Missing:  0,
				},
				BoostMode: "sum",
			},
		},
	}
}

func searchResponse(response *http.Response) ([]string, error) {
	var result Result
	err := json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, hit := range result.Hits.Hits {
		ids = append(ids, hit.ID)
	}

	return ids, nil
}
