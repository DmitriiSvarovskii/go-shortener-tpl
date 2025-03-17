package models

const (
	TypeSimpleUtterance = "SimpleUtterance"
)

type Request struct {
	Url string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}
