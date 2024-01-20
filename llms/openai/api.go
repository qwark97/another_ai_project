package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vargspjut/wlog"

	"github.com/qwark97/assistant/llms/openai/model"
)

type api struct {
	log wlog.Logger
	url string
	key string
}

func (a api) askModel(ctx context.Context, r model.Request) (model.Response, error) {
	var container model.Response
	uri := fmt.Sprintf("%s%s", a.url, "chat/completions")

	data, err := json.Marshal(r)
	if err != nil {
		return model.Response{}, err
	}

	body := bytes.NewReader(data)

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, body)
	if err != nil {
		return model.Response{}, nil
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.key))

	c := &http.Client{}
	response, err := c.Do(request)
	if err != nil {
		return model.Response{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		a.log.Error(response.Status)
		return container, fmt.Errorf(response.Status)
	}

	err = json.NewDecoder(response.Body).Decode(&container)
	if err != nil {
		return container, err
	}

	return container, nil
}
