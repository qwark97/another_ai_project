package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vargspjut/wlog"

	"github.com/qwark97/another_ai_project/llms/model"
)

type api struct {
	log          wlog.Logger
	url          string
	embeddingURL string
	key          string
}

func (a api) askModel(ctx context.Context, r model.Request) (model.Response, error) {
	var container model.Response
	uri := fmt.Sprintf("%s%s", a.url, "chat/completions")

	data, err := json.Marshal(r)
	if err != nil {
		return container, err
	}
	a.log.Debug("%s", string(data))

	body := bytes.NewReader(data)

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, body)
	if err != nil {
		return container, nil
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.key))

	c := &http.Client{}
	response, err := c.Do(request)
	if err != nil {
		return container, err
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&container)
	if err != nil {
		return container, err
	}
	a.log.Debugf("full response: %+v", container)

	if response.StatusCode != http.StatusOK {
		a.log.Error(response.Status)
		return container, fmt.Errorf(response.Status)
	}

	return container, nil
}

func (a api) getEmbeddings(ctx context.Context, r model.EmbeddingRequest) (model.EmbeddingResponse, error) {
	a.log.Debugf("full request: %+v", r)
	var container model.EmbeddingResponse

	data, err := json.Marshal(r)
	if err != nil {
		return container, err
	}

	body := bytes.NewReader(data)

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, a.embeddingURL, body)
	if err != nil {
		return container, nil
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.key))

	c := &http.Client{}
	response, err := c.Do(request)
	if err != nil {
		return model.EmbeddingResponse{}, err
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
	a.log.Debugf("full response: %+v", container)

	return container, nil
}
