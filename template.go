package main

type Template struct {
	Id      string   `json:"id" bson:"_id"`
	Name    string   `json:"name"`
	Content string   `json:"content"`
	Format  string   `json:"format"`
	Tokens  []string `json:"tokens"`
}
