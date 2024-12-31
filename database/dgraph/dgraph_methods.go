package dgraph

import (
	"context"
	"encoding/json"
	"time"

	"github.com/dgraph-io/dgo/v240/protos/api"
)

type School struct {
	Name string `json:"name,omitempty"`
}

type Loc struct {
	Type   string    `json:"type,omitempty"`
	Coords []float64 `json:"coordinates,omitempty"`
}

type Person struct {
	Uid      string     `json:"uid,omitempty"`
	Name     string     `json:"name,omitempty"`
	Age      int        `json:"age,omitempty"`
	Dob      *time.Time `json:"dob,omitempty"`
	Married  bool       `json:"married"`
	Raw      []byte     `json:"raw_bytes,omitempty"`
	Friends  []Person   `json:"friend,omitempty"`
	Location Loc        `json:"loc,omitempty"`
	School   []School   `json:"school,omitempty"`
}

func (d *Dgraph) GetSchema() (string, error) {
	op := &api.Operation{}
	op.Schema = `
		name: string @index(exact) .
		age: int .
		married: bool .
		loc: geo .
		dob: datetime .
	`

	ctx := context.Background()
	err := d.client.Alter(ctx, op)
	if err != nil {
		return "", err
	}

	resp, err := d.client.NewTxn().Query(ctx, `schema(pred: [name, age]) {type}`)
	if err != nil {
		return "", err
	}

	return string(resp.Json), nil
}

func (d *Dgraph) SetObject() (*Person, error) {
	loc, _ := time.LoadLocation("Asia/Tashkent")
	dob := time.Date(2002, time.October, 14, 16, 30, 0, 0, loc)

	person := Person{
		Uid:     "_:doniyorbek",
		Name:    "Doniyorbek",
		Age:     22,
		Married: false,
		Location: Loc{
			Type:   "Point",
			Coords: []float64{1.1, 2},
		},
		Dob: &dob,
		Raw: []byte("raw_bytes"),
		Friends: []Person{
			{Name: "Bob", Age: 24},
			{Name: "Charlie", Age: 29},
		},
		School: []School{
			{Name: "Crown Public School"},
		},
	}

	ctx := context.Background()
	pb, err := json.Marshal(person)
	if err != nil {
		return nil, err
	}

	mu := &api.Mutation{CommitNow: true}
	mu.SetJson = pb
	assigned, err := d.client.NewTxn().Mutate(ctx, mu)
	if err != nil {
		return nil, err
	}

	variables := map[string]string{"$id1": assigned.Uids["doniyorbek"]}
	query := `query Me($id1: string){
		me(func: uid($id1)) {
			name
			dob
			age
			loc
			raw_bytes
			married
			friend @filter(eq(name, "Bob")){
				name
				age
			}
			school {
				name
			}
		}
	}`

	resp, err := d.client.NewTxn().QueryWithVars(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	type Root struct {
		Me []Person `json:"me"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return nil, err
	}

	return &r.Me[0], nil
}
