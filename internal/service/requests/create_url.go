package requests

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CreateUrl struct {
	Url string `json:"url" validate:"required,url"`
}

func NewCreateUrl(r *http.Request) (*CreateUrl, error) {
	var req CreateUrl

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("error decoding create url request: %v", err)
	}

	err := validate.Struct(req)
	if err != nil {
		return nil, fmt.Errorf("error validating create url request: %v", err)
	}

	return &req, nil
}
