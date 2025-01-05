package kuroko2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Client struct {
	client http.Client
	url    string
}

func NewClient(url string, name string, apikey string) Client {
	return Client{
		client: http.Client{
			Transport: basicAuthTransport{name: name, apikey: apikey},
		},
		url: url,
	}
}

type basicAuthTransport struct {
	name   string
	apikey string
}

func (t basicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(t.name, t.apikey)
	return http.DefaultTransport.RoundTrip(req)
}

type JobDefinition struct {
	Id                 int64    `json:"id"`
	Name               string   `json:"name"`
	Description        string   `json:"description"`
	Admins             []int64  `json:"user_id"`
	Script             string   `json:"script"`
	Cron               []string `json:"cron"`
	Tags               []string `json:"tags"`
	NotifyCancellation bool     `json:"notify_cancellation"`
	Suspended          bool     `json:"suspended"`
	PreventMulti       int32    `json:"prevent_multi"`
	SlackChannel       string   `json:"slack_channel"`
}

func (c Client) GetJobDefinition(ctx context.Context, id int64) (JobDefinition, error) {
	var d JobDefinition

	url := fmt.Sprintf("%s/definitions/%d", c.url, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return d, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return d, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return d, fmt.Errorf("GET %s returned unexpected status code: %d", url, resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return d, err
	}

	normalizeJobDefinition(&d)

	return d, nil
}

type JobDefinitionModel struct {
	Name               string   `json:"name"`
	Description        string   `json:"description"`
	Script             string   `json:"script"`
	Admins             []int64  `json:"user_id"`
	Cron               []string `json:"cron"`
	Tags               []string `json:"tags"`
	NotifyCancellation bool     `json:"notify_cancellation"`
	Suspended          bool     `json:"suspended"`
	PreventMulti       int32    `json:"prevent_multi"`
	SlackChannel       string   `json:"slack_channel"`
}

func (c Client) CreateJobDefinition(ctx context.Context, model JobDefinitionModel) (JobDefinition, error) {
	var d JobDefinition

	body, err := json.Marshal(model)
	if err != nil {
		return d, err
	}
	url := fmt.Sprintf("%s/definitions", c.url)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return d, err
	}
	req.Header.Add("content-type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return d, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return d, fmt.Errorf("POST %s/definitions returned unexpected status code: %d", url, resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return d, err
	}

	normalizeJobDefinition(&d)

	return d, nil
}

func (c Client) UpdateJobDefinition(ctx context.Context, id int64, model JobDefinitionModel) error {
	body, err := json.Marshal(model)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/definitions/%d", c.url, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Add("content-type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("UPDATE %s returned unexpected status code: %d", url, resp.StatusCode)
	}

	return nil
}

func (c Client) DeleteJobDefinition(ctx context.Context, id int64) error {
	url := fmt.Sprintf("%s/definitions/%d", c.url, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("DELETE %s returned unexpected status code: %d", url, resp.StatusCode)
	}

	return nil
}

func normalizeJobDefinition(d *JobDefinition) {
	d.Description = strings.ReplaceAll(d.Description, "\r\n", "\n")
	d.Script = strings.ReplaceAll(d.Script, "\r\n", "\n")
}
