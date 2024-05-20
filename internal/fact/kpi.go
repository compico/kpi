package fact

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	HttpClient *http.Client
	BaseUrl    string
	Token      string
}

func NewClient(token string) *Client {
	return &Client{
		HttpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		Token:   token,
		BaseUrl: "https://development.kpi-drive.ru",
	}
}

type SaveFactResponse struct {
	Messages struct {
		Error   []string `json:"error"`
		Warning []string `json:"warning"`
		Info    []string `json:"info"`
	} `json:"MESSAGES"`
	Data struct {
		IndicatorToMoFactId int `json:"indicator_to_mo_fact_id"`
	} `json:"DATA"`
	Status string `json:"STATUS"`
}

// SendFact отправка фактов в api сервер
func (c *Client) SendFact(ctx context.Context, model *Fact) (*SaveFactResponse, error) {
	path := "/_api/facts/save_fact"
	payload, contentType, err := model.Payload()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseUrl+path, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Add("Authorization", "Bearer "+c.Token)
	resp, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("error closing response body: %v\n", err)
		}
	}()

	saveFactResponse := &SaveFactResponse{}

	return saveFactResponse, json.NewDecoder(resp.Body).Decode(saveFactResponse)
}
