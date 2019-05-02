package opsgenie

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/giantswarm/microerror"
)

type Config struct {
	Key string
}

type Client struct {
	httpClient *http.Client
}

func New(config Config) (*Client, error) {
	if config.Key == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Key must not be empty", config)
	}

	httpClient := &http.Client{
		Transport: opsgenieTransport{
			transport: http.DefaultTransport,
			key:       config.Key,
		},
	}

	c := &Client{
		httpClient: httpClient,
	}

	return c, nil
}

func (c *Client) doCountRequest(query string) (int, error) {
	type CountResponseData struct {
		Count int `json:"count"`
	}

	type CountResponse struct {
		Data CountResponseData `json:"data"`
	}

	url := "https://api.opsgenie.com/v2/alerts/count"

	if query != "" {
		url = fmt.Sprintf("%s?query=%s", url, query)
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, microerror.Mask(err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return 0, microerror.Mask(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll((res.Body))
	if err != nil {
		return 0, microerror.Mask(err)
	}

	cr := CountResponse{}
	if err := json.Unmarshal(body, &cr); err != nil {
		return 0, microerror.Mask(err)
	}

	return cr.Data.Count, nil
}

func (c *Client) CountAlerts() (int, error) {
	return c.doCountRequest("")
}

func (c *Client) CountOpenAlerts() (int, error) {
	return c.doCountRequest("status%3Aopen")
}

func (c *Client) CountClosedAlerts() (int, error) {
	return c.doCountRequest("status%3Aclosed")
}
